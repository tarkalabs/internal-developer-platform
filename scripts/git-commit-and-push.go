package main
import (
  "fmt"
  "os"
  "bytes"
  "strings"
  "os/exec"
  "text/template"
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

func createAppShellScript(svcDef SvcDefinition) (string) {
  tmplFile := os.Getenv("GIT_SHELL_SCRIPT_TEMPLATE_PATH")
  tmpl, err := template.ParseFiles(tmplFile)
  checkError(err)
  out, err := os.CreateTemp("", svcDef.Name)
  checkError(err)
  defer out.Close()
  err = tmpl.Execute(out, svcDef)
  checkError(err)
  return out.Name()
}

func generateGitDeployKeyFile(svcDef SvcDefinition) (string) {
  file, err := os.CreateTemp("", svcDef.Name)
  checkError(err)
  defer file.Close()
  _, errWrite := file.WriteString(fmt.Sprintln(svcDef.GithubDeployKey))
  checkError(errWrite)
  return file.Name()
}

func commitAndPush(svcDef SvcDefinition) {
  // Deploy keys won't work with https urls
  svcDef.GitRepo = strings.Replace(svcDef.GitRepo, "https://github.com/", "git@github.com:", -1)
  shellFilePath := createAppShellScript(svcDef)
  gitDeployKeyFilePath := generateGitDeployKeyFile(svcDef)
  defer os.Remove(shellFilePath)
  defer os.Remove(gitDeployKeyFilePath)
  runSystemCommand("bash", shellFilePath, gitDeployKeyFilePath, "DEBUG_MODE=true")
  fmt.Println("Git initial push completed for application: " + svcDef.Name)
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }
  for _, svcDef := range svcDefs { commitAndPush(svcDef) }
}
