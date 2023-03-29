package main

var SupportedTemplates = map[string][]string{
  "nodejs": []string{"18"},
  "react": []string{"18"},
}

var SupportedPredefinedTemplates = []string{
  "nodejs:18|react:18",
}

type SvcDefinition struct {
  Language string `json:"language"`
  MajorVersion string `json:"major_version"`
  Name string `json:"name"`
  SlugName string `json:"slug_name"`
  Environment string `json:"environment"`
  EnvPrefix string `json:"env_prefix"`
  Description string `json:"description"`
  GitRepo string `json:"git_repo"`
  GitBranch string `json:"git_branch"`
  GithubDeployKey string `json:"github_deploy_key"`
  DBConnString string `json:"db_conn_string"`
  // below attributes will get overridden while generating sources
  Namespace string `json:"image"`
}
