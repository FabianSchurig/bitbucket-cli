# Auto-generated Terraform test for bitbucket_user_emails
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_user_emails" {
  command = apply

  variables {
    email = "user@example.com"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_user_emails.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_user_emails"
  }
}
