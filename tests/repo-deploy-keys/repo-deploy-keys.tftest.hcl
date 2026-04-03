# Auto-generated Terraform test for bitbucket_repo_deploy_keys
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_repo_deploy_keys" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    key_id = "123"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_repo_deploy_keys.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_repo_deploy_keys"
  }
}

run "create_repo_deploy_keys" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    key_id = "123"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_repo_deploy_keys.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_repo_deploy_keys"
  }
}
