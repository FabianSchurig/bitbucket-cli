data "bitbucket_project_group_permissions" "example" {
  group_slug = "developers"
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_group_permissions_response" {
  value = data.bitbucket_project_group_permissions.example.api_response
}
