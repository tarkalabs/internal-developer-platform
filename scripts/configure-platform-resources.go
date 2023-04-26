package main
import (
  "fmt"
  "os"
  "bytes"
  "strings"
  "os/exec"
  "encoding/json"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func runSystemCommand(name string, args ...string) {
  cmd := exec.Command(name, args...)
  fmt.Println("cmd:", cmd)
  var stdout bytes.Buffer
  var stderr bytes.Buffer
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  err := cmd.Run()
  if os.Getenv("DEBUG_MODE") == "true" {
    fmt.Printf("\nstdout:\n" + stdout.String() + "\nstderr:\n" + stderr.String() + "\n")
  }
  checkError(err)
}

func getSecurityTrimmedSvcJson(svcDef SvcDefinition) (string) {
  // Clearing all security related keys
  svcDef.GithubDeployKey = ""
  svcDef.GithubWebhookPAT = ""
  svcDef.DBConnString = ""
  svcDef.GithubWebhookSecretToken = ""

  svcJson, err := json.Marshal(svcDef)
  checkError(err)
  return string(svcJson)
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }

  for _, svcDef := range svcDefs {
    fmt.Println("Creating tekton trigger for app:", svcDef.Name)
    ownerAndRepoFormatted := strings.ToLower(strings.Replace(strings.ReplaceAll(svcDef.OwnerAndRepo, "_", "-"), "/", "-", 1))
    helm_args := []string {
      "upgrade", "--install", "--wait", // "--dry-run",
      "--namespace", os.Getenv("PLATFORM_NAMESPACE"),
      fmt.Sprintf("%s-%s-tekton", svcDef.EnvPrefix, svcDef.SlugName),
      "--set", fmt.Sprintf("productName=%s", svcDef.ProductName),
      "--set", fmt.Sprintf("appName=%s", svcDef.Name),
      "--set", fmt.Sprintf("appSlugName=%s", svcDef.SlugName),
      "--set", fmt.Sprintf("environment=%s", svcDef.Environment),
      "--set", fmt.Sprintf("envPrefix=%s", svcDef.EnvPrefix),
      "--set", fmt.Sprintf("namespace=%s", svcDef.Namespace),
      "--set", fmt.Sprintf("platform.secretName=%s", svcDef.GitBranch + "-" + ownerAndRepoFormatted),
      "--set", fmt.Sprintf("platform.createAppResources=false"),
      "--set", fmt.Sprintf("platform.helmChartGithubUrl=%s", os.Getenv("HELM_CHART_GITHUB_URL")),
      "--set", fmt.Sprintf("platform.appIAMRoleARN=%s", os.Getenv("APP_IAM_ROLE_ARN")),
      "--set", fmt.Sprintf("platform.appPipelineIAMRoleARN=%s", os.Getenv("PIPELINE_IAM_ROLE_ARN")),
      "--set", fmt.Sprintf("platform.pipelineStorageClass=%s", os.Getenv("PIPELINE_DEFAULT_STORAGE_CLASS")),
      "--set", fmt.Sprintf("platform.namespace=%s", os.Getenv("PLATFORM_NAMESPACE")),
      "--set", fmt.Sprintf("secretData.svcJson='%s'", getSecurityTrimmedSvcJson(svcDef)),
      "--set", fmt.Sprintf("secretData.githubToken=%s", svcDef.GithubWebhookSecretToken),
      "--set", fmt.Sprintf("secretData.sshDeployKey=%s", svcDef.GithubDeployKey),
      "--set", fmt.Sprintf("secretData.ecrRepoUrl=%s", os.Getenv("BASE_ECR_URL") + "/" + svcDef.ProductName),
      "--set", fmt.Sprintf("tekton.pipelineBaseImage=%s", os.Getenv("PIPELINE_BASE_IMAGE")),
      "--set", fmt.Sprintf("tekton.ingress.domain=%s", svcDef.GithubWebhookDomain),
      "--set", fmt.Sprintf("tekton.ingress.pathPrefix=%s", svcDef.GithubWebhookPathPrefix),
      os.Getenv("HELM_CHART_PATH"),
    }
    if os.Getenv("DEBUG_MODE") == "true" {
      helm_args = append(helm_args, "--debug")
    }
    runSystemCommand("helm", helm_args...)
    fmt.Println("Tekton resources created successfully")
  }
}
