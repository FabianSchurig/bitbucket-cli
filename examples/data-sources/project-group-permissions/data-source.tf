data "bitbucket_project_group_permissions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  group_slug = "developers"
}

output "project_group_permissions_response" {
  value = data.bitbucket_project_group_permissions.example.api_response
}
