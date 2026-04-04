data "bitbucket_issue_comments" "example" {
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issue_comments_response" {
  value = data.bitbucket_issue_comments.example.api_response
}
