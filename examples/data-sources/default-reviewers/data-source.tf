data "bitbucket_default_reviewers" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  target_username = "jdoe"
}

output "default_reviewers_response" {
  value = data.bitbucket_default_reviewers.example.api_response
}
