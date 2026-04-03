resource "bitbucket_hooks" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  uid = "webhook-uuid"
}
