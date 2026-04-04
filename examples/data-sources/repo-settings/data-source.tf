data "bitbucket_repo_settings" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_settings_response" {
  value = data.bitbucket_repo_settings.example.api_response
}
