data "bitbucket_pr_comments" "example" {
  comment_id = "1"
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_comments_response" {
  value = data.bitbucket_pr_comments.example.api_response
}
