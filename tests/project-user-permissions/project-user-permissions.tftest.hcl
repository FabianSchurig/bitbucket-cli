# Auto-generated Terraform test for bitbucket_project_user_permissions
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_project_user_permissions" {
  command = apply

  variables {
    project_key = "PROJ"
    selected_user_id = "{user-uuid}"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_project_user_permissions.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_project_user_permissions"
  }
}
