data "bitbucket_project_branch_restrictions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "project_branch_restrictions_response" {
  value = data.bitbucket_project_branch_restrictions.example.api_response
}
