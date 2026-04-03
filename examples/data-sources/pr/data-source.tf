data "bitbucket_pr" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pull_request_id = "1"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
