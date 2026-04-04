resource "bitbucket_issue_comments" "example" {
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
