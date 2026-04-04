resource "bitbucket_default_reviewers" "example" {
  repo_slug = "my-repo"
  target_username = "jdoe"
  workspace = "my-workspace"
}
