# Auto-generated Terraform test for bitbucket_project_branch_restrictions_by_branch_type
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_project_branch_restrictions_by_branch_type" {
  command = apply

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
    branch_type = "example-value"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_project_branch_restrictions_by_branch_type.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_project_branch_restrictions_by_branch_type"
  }
}

run "create_project_branch_restrictions_by_branch_type" {
  command = apply

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
    branch_type = "example-value"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_project_branch_restrictions_by_branch_type.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_project_branch_restrictions_by_branch_type"
  }
}
