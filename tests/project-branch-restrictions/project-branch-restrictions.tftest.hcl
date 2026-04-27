# Auto-generated Terraform test for bitbucket_project_branch_restrictions
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_project_branch_restrictions" {
  command = apply

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
    pattern = "example-value"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_project_branch_restrictions.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_project_branch_restrictions"
  }
}
