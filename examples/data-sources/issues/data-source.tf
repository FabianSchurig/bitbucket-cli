data "bitbucket_issues" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  issue_id = "1"
}

output "issues_response" {
  value = data.bitbucket_issues.example.api_response
}
