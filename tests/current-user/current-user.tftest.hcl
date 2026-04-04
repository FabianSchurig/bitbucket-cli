# Auto-generated Terraform test for bitbucket_current_user
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_current_user" {
  command = apply

  variables {
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_current_user.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_current_user"
  }
}
