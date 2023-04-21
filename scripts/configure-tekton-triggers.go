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

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }

  for _, svcDef := range svcDefs {
    fmt.Println("Creating tekton trigger for app:", svcDef.Name)
    ownerAndRepoFormatted := strings.ToLower(strings.Replace(strings.ReplaceAll(svcDef.OwnerAndRepo, "_", "-"), "/", "_"))
    helm_args := []string {
      "upgrade",
      "--install",
      "--wait",
      // "--dry-run",
      "--namespace",
      os.Getenv("PLATFORM_NAMESPACE"),
      fmt.Sprintf("%s-%s-tekton", svcDef.EnvPrefix, svcDef.SlugName),
      "--set",
      fmt.Sprintf("productName=%s", svcDef.ProductName),
      "--set",
      fmt.Sprintf("appName=%s", svcDef.Name),
      "--set",
      fmt.Sprintf("appSlugName=%s", svcDef.SlugName),
      "--set",
      fmt.Sprintf("environment=%s", svcDef.Environment),
      "--set",
      fmt.Sprintf("envPrefix=%s", svcDef.EnvPrefix),
      "--set",
      fmt.Sprintf("namespace=%s", svcDef.Namespace),
      "--set",
      fmt.Sprintf("platform.secretName=%s", svcDef.GitBranch + "_" + ownerAndRepoFormatted),
      "--set",
      fmt.Sprintf("platform.createAppResources=false"),
      "--set",
      fmt.Sprintf("platform.namespace=%s", os.Getenv("PLATFORM_NAMESPACE")),
      "--set",
      fmt.Sprintf("secretData.githubToken=%s", svcDef.GithubWebhookSecretToken),
      "--set",
      fmt.Sprintf("secretData.sshDeployKey=%s", svcDef.GithubDeployKey),
      "--set",
      fmt.Sprintf("secretData.ecrRepoUrl=%s", os.Getenv("BASE_ECR_URL") + "/" + svcDef.ProductName),
      "--set",
      fmt.Sprintf("tekton.pipelineBaseImage=%s", os.Getenv("PIPELINE_BASE_IMAGE")),
      "--set",
      fmt.Sprintf("tekton.ingress.domain=%s", svcDef.GithubWebhookDomain),
      "--set",
      fmt.Sprintf("tekton.ingress.pathPrefix=%s", svcDef.GithubWebhookPathPrefix),
      os.Getenv("HELM_CHART_PATH"),
    }
    if os.Getenv("DEBUG_MODE") == "true" {
      helm_args = append(helm_args, "--debug")
    }
    runSystemCommand("helm", helm_args...)
    fmt.Println("Tekton resources created successfully")
  }
}
