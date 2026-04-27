resource "bitbucket_project_branch_restrictions_by_pattern" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  pattern = "example-value"
}
