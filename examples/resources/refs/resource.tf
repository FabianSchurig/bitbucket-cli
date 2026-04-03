resource "bitbucket_refs" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  name = "main"
}
