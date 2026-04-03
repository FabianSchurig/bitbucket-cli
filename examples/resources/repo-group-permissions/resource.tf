resource "bitbucket_repo_group_permissions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  group_slug = "developers"
}
