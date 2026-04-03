resource "bitbucket_commit_statuses" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  key = "build-key"
}
