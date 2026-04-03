resource "bitbucket_commits" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
}
