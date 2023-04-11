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

func makeGithubApiRequest(method string, url string, data string, accessToken string) (int, string) {
  req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
  checkError(err)

  req.Header.Set("Accept", "application/vnd.github+json")
  req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
  req.Header.Set("Authorization", "Bearer " + accessToken)

  client := &http.Client{}
  resp, err := client.Do(req)
  checkError(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  checkError(err)
  return resp.StatusCode, string(body)
}

func getHookIdIfWebhookAlreadyExists(repoHooksUrl string, targetUrl string, access_token string) (string) {
  respStatusCode, body := makeGithubApiRequest("GET", repoHooksUrl, "", access_token)
  if respStatusCode != http.StatusOK {
    panic(fmt.Sprintf("Failed to fetch webhooks! Received %d response with message %s", respStatusCode, body))
  } else {
    var hooks []map[string]interface{}
    if err := json.Unmarshal([]byte(body), &hooks); err != nil { panic(err) }
    for _, hook := range hooks {
      if hook["config"].(map[string]interface{})["url"].(string) == targetUrl {
        return fmt.Sprintf("%f", hook["id"])
      }
    }
  }
  return ""
}

func configureGithubWebhook(svcDef SvcDefinition) {
  fmt.Println("Configuring github webhook for app:", svcDef.Name)
  githubRepoRe := regexp.MustCompile(`.*github\.com\/(.*)\.git$`)
  githubOwnerAndRepo := githubRepoRe.FindStringSubmatch(svcDef.GitRepo)
  data := "{\"name\":\"web\",\"active\":true,\"events\":[\"push\",\"pull_request\"],\"config\":{\"url\":\"" + svcDef.GithubWebhookUrl + "\",\"content_type\":\"json\",\"secret\":\"" + svcDef.GithubWebhookSecretToken + "\"}}"
  repoHooksUrl := "https://api.github.com/repos/" + githubOwnerAndRepo[1] + "/hooks"
  hookId := getHookIdIfWebhookAlreadyExists(repoHooksUrl, svcDef.GithubWebhookUrl, svcDef.GithubWebhookPAT)
  if hookId == "" {
    respStatusCode, body := makeGithubApiRequest("POST", repoHooksUrl, data, svcDef.GithubWebhookPAT)
    if respStatusCode != http.StatusCreated {
      panic(fmt.Sprintf("Failed to create webhook! Received %d response with message %s", respStatusCode, body))
    }
  } else {
    fmt.Println("Webhook exists already...Updating in place!")
    respStatusCode, body := makeGithubApiRequest("PATCH", repoHooksUrl + "/" + hookId, data, svcDef.GithubWebhookPAT)
    if respStatusCode != http.StatusOK {
      panic(fmt.Sprintf("Failed to create webhook! Received %d response with message %s", respStatusCode, body))
    }
  }
  fmt.Println("Webhook configured successfully!")
}

func main() {
  svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
  checkError(err)
  var svcDefs []SvcDefinition
  if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil { panic(err) }
  for _, svcDef := range svcDefs { configureGithubWebhook(svcDef) }
}
