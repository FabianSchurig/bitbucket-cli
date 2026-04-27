data "bitbucket_project_branch_restrictions_by_pattern" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "project_branch_restrictions_by_pattern_response" {
  value = data.bitbucket_project_branch_restrictions_by_pattern.example.api_response
}
