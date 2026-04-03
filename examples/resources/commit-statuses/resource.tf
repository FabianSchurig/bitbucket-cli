resource "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  key = "build-key"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
