package main
import (
  "fmt"
  "os"
  "bytes"
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
  fmt.Printf("\nstdout:\n" + stdout.String() + "\nstderr:\n" + stderr.String() + "\n")
  checkError(err)
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }

  for _, svcDef := range svcDefs {
    fmt.Println("Creating tekton trigger for app:", svcDef.Name)
    helm_args := []string {
      "install",
      "--wait",
      // "--dry-run",
      "--namespace",
      os.Getenv("TEKTON_NAMESPACE"),
      fmt.Sprintf("%s-%s-tekton", svcDef.EnvPrefix, svcDef.SlugName),
      "--set",
      fmt.Sprintf("create.app_resources=false"),
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
      fmt.Sprintf("webhooks.github.token=%s", svcDef.GithubSecretToken),
      "--set",
      fmt.Sprintf("tekton.domain=hooks.%s", svcDef.Domain),
      "--set",
      fmt.Sprintf("tekton.triggerTemplate=tekton-%s-pipeline", svcDef.Language),
      "--set",
      fmt.Sprintf("tekton.namespace=%s", os.Getenv("TEKTON_NAMESPACE")),
      os.Getenv("HELM_CHART_PATH"),
    }
    runSystemCommand("helm", helm_args...)
    fmt.Println("Tekton resources created successfully")
  }
}
