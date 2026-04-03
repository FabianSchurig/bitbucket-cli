data "bitbucket_project_user_permissions" "example" {
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}

output "project_user_permissions_response" {
  value = data.bitbucket_project_user_permissions.example.api_response
}
