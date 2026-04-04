# Auto-generated Terraform test for bitbucket_project_deploy_keys
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_project_deploy_keys" {
  command = apply

  variables {
    key_id = "123"
    project_key = "PROJ"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_project_deploy_keys.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_project_deploy_keys"
  }
}

run "create_project_deploy_keys" {
  command = apply

  variables {
    key_id = "123"
    project_key = "PROJ"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_project_deploy_keys.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_project_deploy_keys"
  }
}
