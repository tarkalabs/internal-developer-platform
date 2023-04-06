package main
import (
  "fmt"
  "os"
  "os/exec"
  "encoding/json"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func runSystemCommand(name string, args ...string) {
  cmd := exec.Command(name, args...)
  fmt.Println(cmd)
  err := cmd.Run()
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
      fmt.Sprintf("install --wait --name %s-%s-tekton", svcDef.EnvPrefix, svcDef.SlugName),
      fmt.Sprintf("--set create.app_resources=false"),
      fmt.Sprintf("--set productName=%s", svcDef.ProductName),
      fmt.Sprintf("--set appName=%s", svcDef.Name),
      fmt.Sprintf("--set appSlugName=%s", svcDef.SlugName),
      fmt.Sprintf("--set environment=%s", svcDef.Environment),
      fmt.Sprintf("--set envPrefix=%s", svcDef.EnvPrefix),
      fmt.Sprintf("--set webhooks.github.token=%s", svcDef.GithubSecretToken),
      fmt.Sprintf("--set tekton.domain=hooks.%s", svcDef.Domain),
      fmt.Sprintf("--set tekton.triggerTemplate=tekton-%s-pipeline", svcDef.Language),
      fmt.Sprintf("--set tekton.namespace=%s", os.Getenv("TEKTON_NAMESPACE")),
    }
    runSystemCommand("helm", helm_args...)
    fmt.Println("Tekton resources created successfully")
  }
}
