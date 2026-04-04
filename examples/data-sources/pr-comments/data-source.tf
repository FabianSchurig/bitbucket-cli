data "bitbucket_pr_comments" "example" {
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
  comment_id = "1"
}

output "pr_comments_response" {
  value = data.bitbucket_pr_comments.example.api_response
}
