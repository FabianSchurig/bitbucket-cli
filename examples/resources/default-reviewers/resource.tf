resource "bitbucket_default_reviewers" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  target_username = "jdoe"
}
