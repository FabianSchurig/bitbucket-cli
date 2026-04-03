data "bitbucket_commit_statuses" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  key = "build-key"
}

output "commit_statuses_response" {
  value = data.bitbucket_commit_statuses.example.api_response
}
