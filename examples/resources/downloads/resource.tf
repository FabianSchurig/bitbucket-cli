resource "bitbucket_downloads" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  filename = "artifact.zip"
}
