package main
import (
  "fmt"
  "os"
  "strings"
  "os/exec"
  "math/rand"
  "path/filepath"
  "encoding/json"
  "text/template"
  "regexp"
	"time"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func randomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func createFolders(path string) {
  cmd := exec.Command("mkdir", "-p", path)
  fmt.Println("Creating directory:", path)
  err := cmd.Run()
  checkError(err)
}

func prefillRequiredData(svcDef *SvcDefinition) {
  svcDef.Name = strings.ReplaceAll(strings.ToLower(svcDef.Name), "_", "-")
  svcDef.ProductName = strings.ReplaceAll(strings.ToLower(svcDef.ProductName), "_", "-")
  if (svcDef.Description == "") { svcDef.Description = "Generated via Tarkalabs IDP" }
  if (svcDef.ProductSlugName == "") { svcDef.ProductSlugName = svcDef.ProductName }
  if (svcDef.SlugName == "") { svcDef.SlugName = svcDef.Name }
  svcDef.SlugName = strings.ReplaceAll(strings.ToLower(svcDef.SlugName), "_", "-")
  if (svcDef.GitBranch == "") { svcDef.GitBranch = "main" }
  if (svcDef.Environment == "") { svcDef.Environment = "production" }
  if (svcDef.EnvPrefix == "") { svcDef.EnvPrefix = svcDef.Environment }

  svcDef.ProductName = strings.ReplaceAll(strings.ToLower(svcDef.ProductName), "_", "-")
  svcDef.Namespace = svcDef.EnvPrefix + "-" + svcDef.ProductName
  svcDef.Domain = strings.ToLower(svcDef.Namespace) + "." + os.Getenv("BASE_DOMAIN")
  if strings.ToLower(svcDef.Type) == "frontend" {
    svcDef.AppDomain = svcDef.EnvPrefix + "-" + svcDef.ProductSlugName + "." + os.Getenv("BASE_DOMAIN")
    svcDef.PathPrefix = "/"
  } else {
    svcDef.AppDomain = svcDef.EnvPrefix + "-" + svcDef.ProductSlugName + "-api." + os.Getenv("BASE_DOMAIN")
    svcDef.PathPrefix = "/api/"
  }

  svcDef.GithubWebhookSecretToken = randomPassword(14)
  if strings.TrimSpace(svcDef.GithubWebhookPAT) == "" {
    svcDef.GithubWebhookPAT = os.Getenv("GITHUB_WEBHOOK_ACCESS_TOKEN")
  }

  svcDef.GithubWebhookDomain = "hooks." + os.Getenv("BASE_DOMAIN")
  svcDef.GithubWebhookPathPrefix = "/" + svcDef.Namespace + "/" + svcDef.EnvPrefix + "-" + svcDef.SlugName
  svcDef.GithubWebhookUrl = "https://" + svcDef.GithubWebhookDomain + svcDef.GithubWebhookPathPrefix

  svcDef.AdminName = os.Getenv("ADMIN_NAME")
  svcDef.AdminEmail = os.Getenv("ADMIN_EMAIL")
  svcDef.GeneratedFilesPath = filepath.Join(os.Getenv("OUTPUT_PATH"), svcDef.Name)

  githubRepoRe := regexp.MustCompile(`.*github\.com\/(.*)\.git$`)
  svcDef.OwnerAndRepo = githubRepoRe.FindStringSubmatch(svcDef.GitRepo)[1]
}

// Not being used right now
func generateKubernetesManifests(svcDef SvcDefinition) {
  createFolders(filepath.Join(svcDef.GeneratedFilesPath, "kubernetes"))
  tmplFiles, err := filepath.Glob(filepath.Join(os.Getenv("K8S_MANIFESTS_PATH"), svcDef.Language, "*.yml.tmpl"))
  checkError(err)
  for _, tmplFile := range tmplFiles {
    tmpl, err := template.ParseFiles(tmplFile)
    checkError(err)
    out, err := os.Create(filepath.Join(svcDef.GeneratedFilesPath, "kubernetes", filepath.Base(tmplFile[:len(tmplFile)-5])))
    checkError(err)
    err = tmpl.Execute(out, svcDef)
    checkError(err)
    out.Close()
  }
}

// Not being used right now
func generateGithubWorkflow(svcDef SvcDefinition) {
  createFolders(filepath.Join(os.Getenv("OUTPUT_PATH"), svcDef.Name, ".github", "workflows"))
  tmplFilePath := filepath.Join(os.Getenv("GITHUB_WORKFLOWS_PATH"), svcDef.Language, "deploy.yml.tmpl")
  tmpl, err := template.ParseFiles(tmplFilePath)
  checkError(err)
  out, err := os.Create(filepath.Join(os.Getenv("OUTPUT_PATH"), svcDef.Name, ".github", "workflows", filepath.Base(tmplFilePath[:len(tmplFilePath)-5])))
  checkError(err)
  err = tmpl.Execute(out, svcDef)
  checkError(err)
  out.Close()
}

func runSystemCommand(name string, arg ...string) {
  cmd := exec.Command(name, arg...)
  fmt.Println(cmd)
  err := cmd.Run()
  checkError(err)
}

func copyRequiredFiles(svcDef SvcDefinition) {
  createFolders(svcDef.GeneratedFilesPath)
  filePaths, _ := filepath.Glob(filepath.Join(os.Getenv("APP_TEMPLATES_PATH"), svcDef.Language) + string(filepath.Separator) + "*")
  for _, filePath := range filePaths {
    fi, err := os.Lstat(filePath)
    checkError(err)
    if fi.Mode().IsRegular() {
      fmt.Println("Copying file", filePath)
      data, err := os.ReadFile(filePath)
      checkError(err)
      err = os.WriteFile(filepath.Join(svcDef.GeneratedFilesPath, filepath.Base(filePath)), data, 0644)
      checkError(err)
    }
  }

  fmt.Println("Copying", svcDef.Language, svcDef.MajorVersion, "version specific files...")
  appTemplatePath := filepath.Join(os.Getenv("APP_TEMPLATES_PATH"), svcDef.Language, svcDef.MajorVersion)
  runSystemCommand("bash", "-c", "cp -rf " + appTemplatePath + string(filepath.Separator) + "* " + svcDef.GeneratedFilesPath + string(filepath.Separator))
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
    svcDefs[0].Type         = "frontend"
    svcDefs[0].Language     = frontend_def[0]
    svcDefs[0].MajorVersion = frontend_def[1]
    svcDefs[1].Type         = "backend"
    svcDefs[1].Language     = backend_def[0]
    svcDefs[1].MajorVersion = backend_def[1]
    svcDefs                 = svcDefs[0:2]
  }
  fmt.Println("Clearing contents of output folder", os.Getenv("OUTPUT_PATH"))
  outputPath := os.Getenv("OUTPUT_PATH")
  if ! strings.HasSuffix(outputPath, "/") { outputPath += "/" }
  runSystemCommand("rm", "-rf", outputPath)
  for i, svcDef := range svcDefs {
    fmt.Println("Generating resources for app:", svcDef.Name)
    prefillRequiredData(&svcDef)
    svcDefs[i] = svcDef
    copyRequiredFiles(svcDef)
    fmt.Println("Generation of resources completed for app:", svcDef.Name)
  }

  servicesJson, err := json.Marshal(svcDefs)
  checkError(err)
  err = os.WriteFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"), servicesJson, 0644)
  checkError(err)
}
