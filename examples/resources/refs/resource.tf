resource "bitbucket_refs" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
