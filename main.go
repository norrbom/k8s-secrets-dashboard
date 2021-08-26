package main

import (
	"context"
	"os"
	"strings"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigs := make(map[string]string)

	type reportData struct {
		Report          report
		CountNamespace  int
		CountCompleted  int
		CountInProgress int
		CountInScope    int
		CountNotStarted int
	}
	reports := make(map[string]reportData)

	// load kube configs
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "KUBECONFIG_") {
			kubeconfigs[strings.TrimLeft(pair[0], "KUBECONFIG_")] = pair[1]
		}
	}

	for env, file := range kubeconfigs {
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", file)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		// access the API to get all secrets
		secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		// Iterate over all secrets to build a report
		report := make(report, 0)

		for _, s := range secrets.Items {
			if s.Type == "Opaque" {
				if len(s.Data) != 0 {
					for secret := range s.Data {
						if rr, found := report[s.Namespace]; found {
							rr.addSecret(secret)
						} else {
							rr := NewReportRow(env, s.Namespace, secret)
							report[s.Namespace] = rr
						}
					}
				}
			}
		}

		var CountNamespace int
		var CountCompleted int
		var CountInProgress int
		var CountInScope int
		var CountNotStarted int
		for _, row := range report {
			CountNamespace += 1
			if row.inScope() {
				CountInScope++
				row.setProgress()
				switch {
				case row.isCompleted():
					CountCompleted++
				case row.isInProgress():
					CountInProgress++
				}
			}
		}
		CountNotStarted = CountInScope - CountCompleted - CountInProgress
		reports[env] = reportData{
			Report: report, CountNamespace: CountNamespace, CountCompleted: CountCompleted,
			CountInProgress: CountInProgress, CountInScope: CountInScope, CountNotStarted: CountNotStarted}
	}
	// render report
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		panic(err.Error())
	}

	f, err := os.Create("./data/index.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	err = tmpl.Execute(f, reports)
	if err != nil {
		panic(err.Error())
	}
}
