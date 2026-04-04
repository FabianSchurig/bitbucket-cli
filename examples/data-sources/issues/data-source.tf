data "bitbucket_issues" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  issue_id = "1"
}

output "issues_response" {
  value = data.bitbucket_issues.example.api_response
}
