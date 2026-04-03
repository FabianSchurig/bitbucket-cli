resource "bitbucket_project_user_permissions" "example" {
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}
