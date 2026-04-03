data "bitbucket_default_reviewers" "example" {
  repo_slug = "my-repo"
  target_username = "jdoe"
  workspace = "my-workspace"
}

output "default_reviewers_response" {
  value = data.bitbucket_default_reviewers.example.api_response
}
