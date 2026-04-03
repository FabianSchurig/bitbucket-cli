# Auto-generated Terraform test for bitbucket_snippets
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_snippets" {
  command = apply

  variables {
    workspace = "my-workspace"
    encoded_id = "snippet-id"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_snippets.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_snippets"
  }
}

run "create_snippets" {
  command = apply

  variables {
    workspace = "my-workspace"
    encoded_id = "snippet-id"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_snippets.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_snippets"
  }
}
