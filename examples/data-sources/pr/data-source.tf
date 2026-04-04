data "bitbucket_pr" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  pull_request_id = "1"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
