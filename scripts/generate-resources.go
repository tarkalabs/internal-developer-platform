package main
import (
  "fmt"
  "os"
  "strings"
  "os/exec"
  "path/filepath"
  "encoding/json"
  "text/template"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func prefillRequiredData(svcDef *SvcDefinition) {
  svcDef.Name = strings.ReplaceAll(strings.ToLower(svcDef.Name), "_", "-")
  if (svcDef.Description == "") { svcDef.Description = "Generated via Tarkalabs IDP" }
  if (svcDef.SlugName == "") { svcDef.SlugName = svcDef.Name }
  svcDef.SlugName = strings.ReplaceAll(strings.ToLower(svcDef.SlugName), "_", "-")
  if (svcDef.GitBranch == "") { svcDef.GitBranch = "main" }
  if (svcDef.Environment == "") { svcDef.Environment = "production" }
  if (svcDef.EnvPrefix == "") { svcDef.EnvPrefix = svcDef.Environment }
  svcDef.Namespace = svcDef.EnvPrefix + "-" + svcDef.SlugName
}

func generateKubernetesManifests(svcDef SvcDefinition) {
  tmplFiles, err := filepath.Glob(filepath.Join(os.Getenv("K8S_MANIFESTS_PATH"), svcDef.Language, "*.yml.tmpl"))
  checkError(err)
  for _, tmplFile := range tmplFiles {
    tmpl, err := template.ParseFiles(tmplFile)
    checkError(err)
    out, err := os.Create(filepath.Join(os.Getenv("OUTPUT_PATH"), "kubernetes", filepath.Base(tmplFile[:len(tmplFile)-5])))
    checkError(err)
    err = tmpl.Execute(out, svcDef)
    checkError(err)
    out.Close()
  }
}

func generateGithubWorkflow(svcDef SvcDefinition) {
  tmplFilePath := filepath.Join(os.Getenv("GITHUB_WORKFLOWS_PATH"), svcDef.Language, "deploy.yml.tmpl")
  tmpl, err := template.ParseFiles(tmplFilePath)
  checkError(err)
  out, err := os.Create(filepath.Join(os.Getenv("OUTPUT_PATH"), ".github", "workflows", filepath.Base(tmplFilePath[:len(tmplFilePath)-5])))
  checkError(err)
  err = tmpl.Execute(out, svcDef)
  checkError(err)
  out.Close()
}

func copyRequiredFiles(svcDef SvcDefinition) {
  filePaths, _ := filepath.Glob(filepath.Join(os.Getenv("APP_TEMPLATES_PATH"), svcDef.Language) + string(filepath.Separator) + "*")
  for _, filePath := range filePaths {
    fi, err := os.Lstat(filePath)
    checkError(err)
    if fi.Mode().IsRegular() {
      fmt.Println("Copying file", filePath)
      data, err := os.ReadFile(filePath)
      checkError(err)
      err = os.WriteFile(filepath.Join(os.Getenv("OUTPUT_PATH"), filepath.Base(filePath)), data, 0644)
      checkError(err)
    }
  }

  fmt.Println("Copying", svcDef.Language, svcDef.MajorVersion, "version specific files...")
  appTemplatePath := filepath.Join(os.Getenv("APP_TEMPLATES_PATH"), svcDef.Language, svcDef.MajorVersion)
  cmd := exec.Command("cp", "-rf", appTemplatePath + string(filepath.Separator), os.Getenv("OUTPUT_PATH"))
  fmt.Println(cmd)
  err := cmd.Run()
  checkError(err)
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }

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

  for _, svcDef := range svcDefs {
    fmt.Println("Generating resources for app:", svcDef.Name)
    prefillRequiredData(&svcDef)
    generateGithubWorkflow(svcDef)
    generateKubernetesManifests(svcDef)
    copyRequiredFiles(svcDef)
    fmt.Println("Generation of resources completed for app:", svcDef.Name)
  }
}
