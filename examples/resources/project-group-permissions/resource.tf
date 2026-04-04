resource "bitbucket_project_group_permissions" "example" {
  group_slug = "developers"
  project_key = "PROJ"
  workspace = "my-workspace"
}
