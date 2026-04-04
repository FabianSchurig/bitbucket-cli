# Auto-generated Terraform test for bitbucket_deployment_variables
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_deployment_variables" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    environment_uuid = "env-uuid"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_deployment_variables.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_deployment_variables"
  }
}

run "create_deployment_variables" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    environment_uuid = "env-uuid"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_deployment_variables.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_deployment_variables"
  }
}
