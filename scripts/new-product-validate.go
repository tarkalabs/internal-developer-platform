package main
import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

type SvcDefinition struct {
  Type string `json:"type"`
	Version string `json:"version"`
  Name string `json:"name"`
	Description string `json:"description"`
  GitRepo string `json:"git_repo"`
	GitBranch string `json:"git_branch"`
	GithubDeployKey string `json:"github_deploy_key"`
  DBConnString string `json:"db_conn_string"`
}

func main() {
	if len(strings.TrimSpace(os.Getenv("PREDEFINED_TEMPLATE"))) > 0 {
		fmt.Println("Using predefined template: ", os.Getenv("PREDEFINED_TEMPLATE"))
	} else {
		fmt.Println("Not a predefined template. Validating each app definition.")
		svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
		if err != nil { panic(err) }
		var svcDefs []SvcDefinition
		if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil {
			panic(err)
		}
		fmt.Println("App definitions found total: ", len(svcDefs))
    for _, svcDef := range svcDefs {
			svcJson,_ := json.Marshal(svcDef)
			fmt.Println("Validating app definition: ", string(svcJson))
			if strings.TrimSpace(svcDef.Name) == "" {
				panic("One of the app definitions doesn't have `name` defined!")
			} else if strings.TrimSpace(svcDef.Type) == "" {
					panic(svcDef.Name + " app definition doesn't have `type` defined!")
			} else if strings.TrimSpace(svcDef.Version) == "" {
				panic(svcDef.Name + " app definition doesn't have `version` defined!")
			} else if strings.TrimSpace(svcDef.GitRepo) == "" {
				panic(svcDef.Name + " app definition doesn't have `git_repo` defined!")
			} else if strings.TrimSpace(svcDef.GithubDeployKey) == "" {
				panic(svcDef.Name + " app definition doesn't have `github_deploy_key` defined!")
			} else {
				fmt.Printf("Microservice %s validation successful!\n", svcDef.Name)
			}
		}
	}
}

// go run template1.tmpl template2.tmpl | tee >(sed -e 's/.tmpl/.out/g' | xargs -n 1 touch)
