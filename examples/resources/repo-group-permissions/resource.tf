resource "bitbucket_repo_group_permissions" "example" {
  group_slug = "developers"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
