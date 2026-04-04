data "bitbucket_pr" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
