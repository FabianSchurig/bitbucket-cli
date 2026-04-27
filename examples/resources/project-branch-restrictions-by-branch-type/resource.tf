resource "bitbucket_project_branch_restrictions_by_branch_type" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  branch_type = "example-value"
}
