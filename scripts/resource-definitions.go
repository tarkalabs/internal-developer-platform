package main

var SupportedTemplates = map[string][]string{
  "nodejs": []string{"18"},
  "react": []string{"18"},
}

var SupportedPredefinedTemplates = []string{
  "nodejs:18|react:18",
}

type SvcDefinition struct {
  Type string `json:"type"`
  Language string `json:"language"`
  MajorVersion string `json:"major_version"`
  ProductName string `json:"product_name"`
  AdminName string `json:"admin_name"`
  AdminEmail string `json:"admin_email"`
  Name string `json:"name"`
  SlugName string `json:"slug_name"`
  Environment string `json:"environment"`
  EnvPrefix string `json:"env_prefix"`
  Description string `json:"description"`
  GitRepo string `json:"git_repo"`
  GitBranch string `json:"git_branch"`
  GithubDeployKey string `json:"github_deploy_key"`
  GithubWebhookPAT string `json:"github_webhook_pat"`
  DBConnString string `json:"db_conn_string"`
  // below attributes will get overridden while generating sources
  Namespace string `json:"namespace"`
  Domain string `json:"domain"`
  HttpPath string `json:"http_path"`
  GithubWebhookUrl string `json:"github_webhook_url"`
  GithubWebhookDomain string `json:"github_webhook_domain"`
  GithubWebhookPathPrefix string `json:"github_webhook_path_prefix"`
  GithubWebhookSecretToken string `json:"github_webhook_secret_token"`
  GeneratedFilesPath string `json:"generated_files_path"`
}
