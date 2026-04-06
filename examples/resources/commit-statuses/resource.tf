resource "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
  key = "build-key"
  state = "SUCCESSFUL"
}
