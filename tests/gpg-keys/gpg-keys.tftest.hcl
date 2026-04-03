# Auto-generated Terraform test for bitbucket_gpg_keys
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_gpg_keys" {
  command = apply

  variables {
    fingerprint = "AA:BB:CC:DD"
    selected_user = "jdoe"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_gpg_keys.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_gpg_keys"
  }
}

run "create_gpg_keys" {
  command = apply

  variables {
    fingerprint = "AA:BB:CC:DD"
    selected_user = "jdoe"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_gpg_keys.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_gpg_keys"
  }
}
