data "bitbucket_pr" "example" {
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
