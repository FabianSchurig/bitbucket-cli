resource "bitbucket_project_user_permissions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
}
