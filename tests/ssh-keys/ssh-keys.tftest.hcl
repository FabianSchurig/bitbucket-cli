# Auto-generated Terraform test for bitbucket_ssh_keys
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_ssh_keys" {
  command = apply

  variables {
    key_id = "123"
    selected_user = "jdoe"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_ssh_keys.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_ssh_keys"
  }
}

run "create_ssh_keys" {
  command = apply

  variables {
    key_id = "123"
    selected_user = "jdoe"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_ssh_keys.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_ssh_keys"
  }
}
