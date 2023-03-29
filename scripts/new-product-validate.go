package main
import (
  "fmt"
  "os"
  "strings"
  "encoding/json"
)

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

func validateMicroservicesDefinitions(svcDefs []SvcDefinition) bool {
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
  return true
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  if err != nil { panic(err) }
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil {
    panic(err)
  }

  predefinedTemplate := strings.TrimSpace(os.Getenv("PREDEFINED_TEMPLATE"))
  if len(predefinedTemplate) > 0 {
    if containsStringArray(SupportedPredefinedTemplates, predefinedTemplate) {
      fmt.Println("Predefined template:", predefinedTemplate, "defined!")
      defs := strings.Split(predefinedTemplate, "|")
      frontend_def := strings.Split(defs[0], ":")
      backend_def  := strings.Split(defs[1], ":")
      svcDefs[0].Language     = frontend_def[0]
      svcDefs[0].MajorVersion = frontend_def[1]
      svcDefs[1].Language     = backend_def[0]
      svcDefs[1].MajorVersion = backend_def[1]
      svcDefs = svcDefs[0:2]
      } else {
      panic("Predefined template named " + predefinedTemplate + " isn't found!")
    }
  }

  if validateMicroservicesDefinitions(svcDefs) {
    fmt.Println("All microservices validation successfully completed!")
  }
}
