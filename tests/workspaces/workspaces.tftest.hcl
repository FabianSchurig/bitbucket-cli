# Auto-generated Terraform test for bitbucket_workspaces
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_workspaces" {
  command = plan

  variables {
    workspace = "my-workspace"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_workspaces.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_workspaces"
  }
}
