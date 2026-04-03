data "bitbucket_commits" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
}

output "commits_response" {
  value = data.bitbucket_commits.example.api_response
}
