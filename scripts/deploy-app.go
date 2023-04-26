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
  var svcDef SvcDefinition
  err := json.Unmarshal([]byte(os.Getenv("SVC_JSON")), &svcDef)
  checkError(err)
  fmt.Println("Deploying app:", svcDef.Name)
  ownerAndRepoFormatted := strings.ToLower(strings.Replace(strings.ReplaceAll(svcDef.OwnerAndRepo, "_", "-"), "/", "-", 1))
  helm_args := []string {
    "upgrade", "--install", "--wait", // "--dry-run",
    "--namespace", os.Getenv("PLATFORM_NAMESPACE"),
    fmt.Sprintf("%s-%s-app", svcDef.EnvPrefix, svcDef.SlugName),
    "--set", fmt.Sprintf("platform.createPlatformResources=false"),
    "--set", fmt.Sprintf("productName=%s", svcDef.ProductName),
    "--set", fmt.Sprintf("appName=%s", svcDef.Name),
    "--set", fmt.Sprintf("appSlugName=%s", svcDef.SlugName),
    "--set", fmt.Sprintf("environment=%s", svcDef.Environment),
    "--set", fmt.Sprintf("envPrefix=%s", svcDef.EnvPrefix),
    "--set", fmt.Sprintf("namespace=%s", svcDef.Namespace),
    "--set", fmt.Sprintf("ingress.domain=%s", svcDef.AppDomain),
    "--set", fmt.Sprintf("ingress.httpPathPrefix=%s", svcDef.PathPrefix),
    "--set", fmt.Sprintf("platform.appIAMRoleARN=%s", os.Getenv("APP_IAM_ROLE_ARN")),
    "--set", fmt.Sprintf("deployment.container.image=%s", os.Getenv("IMAGE")),
    "--set", fmt.Sprintf("deployment.container.port=%s", os.Getenv("APP_PORT")),
    "--set", fmt.Sprintf("deployment.container.livenessProbe.tcpSocket.port=%s", os.Getenv("APP_PORT")),
    "--set", fmt.Sprintf("deployment.container.readinessProbe.httpGet.port=%s", os.Getenv("APP_PORT")),
    os.Getenv("HELM_CHART_PATH"),
  }
  if os.Getenv("DEBUG_MODE") == "true" {
    helm_args = append(helm_args, "--debug")
  }
  runSystemCommand("helm", helm_args...)
  fmt.Println("Application deployed successfully")
}
