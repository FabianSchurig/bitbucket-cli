# Auto-generated Terraform test for bitbucket_workspace_members
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_workspace_members" {
  command = apply

  variables {
    member = "{member-uuid}"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_workspace_members.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_workspace_members"
  }
}
