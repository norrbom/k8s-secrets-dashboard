package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type GitOps struct {
	GitOpsRepositoryUrl      string
	CodeLines                int
	Project_key              string
	Pipeline_version         string
	Version                  string
	Type                     string
	Chart_name               string
	Helm_version             string
	Use_helm3                string
	Clusters                 []string
	Platform_instance_groups []map[string]interface{}
}

/*
factory method for GitOps, usage: gitOps := NewGitOps("stage1.us1", "service")
*/
func NewGitOps(environment string, namespace string) *GitOps {
	g := new(GitOps)
	g.GitOpsRepositoryUrl = config.GitOpsScheme + config.GitOpsHost + config.GitOpsProjectPath + "/" + config.getGitOpsRepo(environment)
	g.getConfigData(environment, namespace)
	g.getReleaseData(environment, namespace)
	return g
}

func (g *GitOps) getReleaseData(environment string, namespace string) {
	releaseFilePath := config.LocalGitOpsProjectPath + "/" + config.getGitOpsRepo(environment) + "/" + namespace + "/release-data.yaml"
	yamlFile, err := ioutil.ReadFile(releaseFilePath)
	if err == nil {
		if err := yaml.Unmarshal(yamlFile, g); err != nil {
			fmt.Printf("[WARN] Could not unmarshall file: %s, error: %s", yamlFile, err.Error())
		}
	}
}

func (g *GitOps) getConfigData(environment string, namespace string) {
	lineCount := 0
	configDir := config.LocalGitOpsProjectPath + "/" + config.getGitOpsRepo(environment) + "/" + namespace + "/config"
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		var files []string
		err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".yaml" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			file, _ := os.Open(file)
			fileScanner := bufio.NewScanner(file)
			for fileScanner.Scan() {
				lineCount++
			}
		}
	}
	g.CodeLines = lineCount
}
