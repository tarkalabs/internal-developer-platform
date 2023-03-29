package main
import (
  "fmt"
  "os"
  "strings"
  "encoding/json"
)

func prefillMissingData(svcDef *SvcDefinition) {
  if (svcDef.Description == "") {
    svcDef.Description = "Generated via Tarkalabs IDP"
  }
  if (svcDef.SlugName == "") {
    svcDef.SlugName = svcDef.Name
  }
  if (svcDef.GitBranch == "") {
    svcDef.GitBranch = "main"
  }
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
    defs := strings.Split(predefinedTemplate, "|")
    frontend_def := strings.Split(defs[0], ":")
    backend_def  := strings.Split(defs[1], ":")
    svcDefs[0].Language     = frontend_def[0]
    svcDefs[0].MajorVersion = frontend_def[1]
    svcDefs[1].Language     = backend_def[0]
    svcDefs[1].MajorVersion = backend_def[1]
    svcDefs                 = svcDefs[0:2]
  }

  fmt.Println("App definitions:", len(svcDefs))
  for _, svcDef := range svcDefs {
    prefillMissingData(&svcDef)
    svcJson,_ := json.Marshal(svcDef)
    fmt.Println("App definition:", string(svcJson))
  }
}
