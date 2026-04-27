data "bitbucket_project_branch_restrictions_by_branch_type" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "project_branch_restrictions_by_branch_type_response" {
  value = data.bitbucket_project_branch_restrictions_by_branch_type.example.api_response
}
