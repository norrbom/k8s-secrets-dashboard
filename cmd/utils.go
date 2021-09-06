package main

import "strings"

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type byImportance []string

func (s byImportance) Len() int {
	return len(s)
}
func (s byImportance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byImportance) Less(i, j int) bool {
	env1 := strings.ToLower(s[i])
	env2 := strings.ToLower(s[j])
	if strings.Contains(env1, "data") && strings.Contains(env2, "data") {
		return environmentLess(env1, env2)
	}
	if strings.Contains(env1, "data") && !strings.Contains(env2, "data") {
		return false
	}
	if !strings.Contains(env1, "data") && strings.Contains(env2, "data") {
		return true
	} else {
		return environmentLess(env1, env2)
	}
}

func environmentLess(env1 string, env2 string) bool {
	if strings.Contains(env1, "prod") {
		return true
	}
	if strings.Contains(env1, "stage") && !strings.Contains(env2, "prod") {
		return true
	}
	if strings.Contains(env1, "pt") && (!strings.Contains(env2, "stage") && !strings.Contains(env2, "prod")) {
		return true
	}
	if strings.Contains(env1, "qa") && (!strings.Contains(env2, "stage") && !strings.Contains(env2, "prod") && !strings.Contains(env2, "pt")) {
		return true
	}
	if strings.Contains(env1, "si") && (!strings.Contains(env2, "stage") && !strings.Contains(env2, "prod") && !strings.Contains(env2, "pt") && !strings.Contains(env2, "qa")) {
		return true
	}
	return false
}
