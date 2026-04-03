resource "bitbucket_pr" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pull_request_id = "1"
}
