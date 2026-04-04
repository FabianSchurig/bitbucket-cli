data "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
  key = "build-key"
}

output "commit_statuses_response" {
  value = data.bitbucket_commit_statuses.example.api_response
}
