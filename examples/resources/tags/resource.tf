resource "bitbucket_tags" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
