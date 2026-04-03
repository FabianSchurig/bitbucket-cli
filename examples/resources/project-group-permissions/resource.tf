resource "bitbucket_project_group_permissions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  group_slug = "developers"
}
