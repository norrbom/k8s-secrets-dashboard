package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var config *Config

func init() {
	config = NewConfig()
}

type Config struct {
	GitUserPass            string
	GitOpsScheme           string
	GitOpsHost             string
	GitOpsProjectPath      string
	LocalGitOpsProjectPath string
	GitProjectPath         string
	GitOpsConfigPath       string
	Env2GitOpsRepos        map[string]string
	Kubeconfigs            map[string]string
	Environments           []string
}

func (c *Config) getGitOpsUrl(environment string) string {
	return c.GitOpsScheme + c.GitOpsHost + c.GitProjectPath + "/" + c.getGitOpsRepo(environment) + ".git"
}

/*
factory method for Config
*/
func NewConfig() *Config {
	c := new(Config)
	c.GitOpsScheme = "https://"
	c.GitOpsHost = "bitbucket.kindredgroup.com"
	c.GitOpsProjectPath = "/bitbucket/projects/DEPLOY/repos"
	c.LocalGitOpsProjectPath = "/tmp" + c.GitOpsProjectPath
	c.GitProjectPath = "/bitbucket/scm/deploy"
	c.GitOpsConfigPath = "/config"
	// TODO: read from GitOps config.yaml
	c.Env2GitOpsRepos = map[string]string{
		"stage1.us1": "prod1",
	}
	if os.Getenv("GIT_USER") != "" && os.Getenv("GIT_PASS") != "" {
		c.GitUserPass = os.Getenv("GIT_USER") + ":" + os.Getenv("GIT_PASS")
	} else {
		fmt.Println("[WARN] GIT_USER/GIT_PASS env variables not set!")
	}
	c.loadKubeconfigs()
	c.cloneGitOpsRepos()
	return c
}

func (c *Config) cloneGitOpsRepos() {
	for _, env := range c.Environments {
		fileSystemPath := c.LocalGitOpsProjectPath + "/" + c.getGitOpsRepo(env)
		if _, err := os.Stat(fileSystemPath); os.IsNotExist(err) {
			_, err := git.PlainClone(fileSystemPath, false, &git.CloneOptions{
				URL: c.getGitOpsUrl(env),
				Auth: &http.BasicAuth{
					Username: os.Getenv("GIT_USER"),
					Password: os.Getenv("GIT_PASS"),
				},
				ReferenceName: plumbing.ReferenceName("refs/heads/master"),
				SingleBranch:  true,
				Depth:         1,
			})
			if err != nil {
				panic("[ERROR] an error occured when cloning " + c.getGitOpsUrl(env) + ", error: " + err.Error())
			}
		} else {
			fmt.Println("[INFO] path already exists: + " + fileSystemPath + ", skipping!")
		}
	}
}

func (c *Config) loadKubeconfigs() {
	c.Kubeconfigs = map[string]string{}
	// load kube configs from env variables
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "KUBECONFIG_") {
			// dots needs to be passed as _ and converted back since they are not accepted in env variables
			c.Kubeconfigs[strings.ReplaceAll(strings.TrimLeft(pair[0], "KUBECONFIG_"), "_", ".")] = pair[1]
		}
	}
	// sort environments by importance
	c.Environments = make([]string, 0, len(c.Kubeconfigs))
	for k := range c.Kubeconfigs {
		c.Environments = append(c.Environments, k)
	}
	sort.Sort(byImportance(c.Environments))
}

func (c *Config) getGitOpsRepo(environment string) string {
	env := strings.ToLower(environment)
	if repoName, found := c.Env2GitOpsRepos[env]; found {
		return repoName
	} else {
		return strings.ToLower(env)
	}
}
