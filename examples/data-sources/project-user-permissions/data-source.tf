data "bitbucket_project_user_permissions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
}

output "project_user_permissions_response" {
  value = data.bitbucket_project_user_permissions.example.api_response
}
