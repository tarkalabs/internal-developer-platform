package main
import (
  "fmt"
  "os"
  "bytes"
  "regexp"
  "io/ioutil"
  "net/http"
  "encoding/json"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func createGithubToken(url string, data string, access_token string) {
  req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
  checkError(err)

  req.Header.Set("Accept", "application/vnd.github+json")
  req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
  req.Header.Set("Authorization", "Bearer " + access_token)

  client := &http.Client{}
  resp, err := client.Do(req)
  checkError(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  checkError(err)
  if resp.StatusCode != http.StatusCreated {
    panic(fmt.Sprintf("Failed to create webhook! Received %d response with message %s", resp.StatusCode, string(body)))
  } else {
    fmt.Println(string(body))
  }
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }

  githubRepoRe := regexp.MustCompile(`.*github\.com\/(.*)\.git$`)

  for _, svcDef := range svcDefs {
    fmt.Println("Configuring github webhook for app:", svcDef.Name)
    githubOwnerAndRepo := githubRepoRe.FindStringSubmatch(svcDef.GitRepo)
    data := "{\"name\":\"web\",\"active\":true,\"events\":[\"push\",\"pull_request\"],\"config\":{\"url\":\"" + svcDef.GithubWebhookUrl + "\",\"content_type\":\"json\",\"secret\":\"" + svcDef.GithubWebhookSecretToken + "\"}}"
    createGithubToken("https://api.github.com/repos/" + githubOwnerAndRepo[1] + "/hooks", data, svcDef.GithubWebhookPAT)
    fmt.Println("Webhook configured successfully")
  }
}
