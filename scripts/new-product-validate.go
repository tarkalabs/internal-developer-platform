package main
import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

var SupportedTemplates = map[string][]string{
	"nodejs": []string{"18"},
}

var SupportedPredefinedTemplates = []string{
	"nodejs:18|react:18",
}

type SvcDefinition struct {
  Language string `json:"language"`
	MajorVersion string `json:"major_version"`
  Name string `json:"name"`
	Description string `json:"description"`
  GitRepo string `json:"git_repo"`
	GitBranch string `json:"git_branch"`
	GithubDeployKey string `json:"github_deploy_key"`
  DBConnString string `json:"db_conn_string"`
}

func containsStringArray(arr []string, target string) bool {
	found := false
	for _, entry := range arr {
		if target == entry {
			found = true
			break
		}
	}
	return found
}

func main() {
	predefinedTemplate := strings.TrimSpace(os.Getenv("PREDEFINED_TEMPLATE"))
	if len(predefinedTemplate) > 0 {
		if containsStringArray(SupportedPredefinedTemplates, predefinedTemplate) {
			fmt.Println("Predefined template:", predefinedTemplate, "exists. Validation successful!")
		} else {
			panic("Predefined template named " + predefinedTemplate + " isn't found!")
		}
	} else {
		fmt.Println("Not a predefined template. Validating each app definition.")
		svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
		if err != nil { panic(err) }
		var svcDefs []SvcDefinition
		if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil {
			panic(err)
		}
		fmt.Println("App definitions found total:", len(svcDefs))
    for _, svcDef := range svcDefs {
			svcJson,_ := json.Marshal(svcDef)
			fmt.Println("Validating app definition:", string(svcJson))
			if strings.TrimSpace(svcDef.Name) == "" {
				panic("One of the app definitions doesn't have `name` defined!")
			} else if strings.TrimSpace(svcDef.Language) == "" {
					panic(svcDef.Name + " app definition doesn't have `language` defined!")
			} else if strings.TrimSpace(svcDef.MajorVersion) == "" {
				panic(svcDef.Name + " app definition doesn't have `major_version` defined!")
			} else if strings.TrimSpace(svcDef.GitRepo) == "" {
				panic(svcDef.Name + " app definition doesn't have `git_repo` defined!")
			} else if strings.TrimSpace(svcDef.GithubDeployKey) == "" {
				panic(svcDef.Name + " app definition doesn't have `github_deploy_key` defined!")
			} else {
				if _, ok := SupportedTemplates[svcDef.Language]; ok {
					if containsStringArray(SupportedTemplates[svcDef.Language], svcDef.MajorVersion) {
						fmt.Printf("Microservice %s validation successful!\n", svcDef.Name)
					} else {
						panic(svcDef.Name + " app definition " + svcDef.Language + " version " + svcDef.MajorVersion + " is not supported!")
					}
				} else {
					panic(svcDef.Name + " app definition has unsupported language type: " + svcDef.Language)
				}
			}
		}
	}
}

// go run template1.tmpl template2.tmpl | tee >(sed -e 's/.tmpl/.out/g' | xargs -n 1 touch)