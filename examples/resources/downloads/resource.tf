resource "bitbucket_downloads" "example" {
  filename = "artifact.zip"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
