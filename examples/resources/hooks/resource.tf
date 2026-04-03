resource "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  uid = "webhook-uuid"
  workspace = "my-workspace"
}
