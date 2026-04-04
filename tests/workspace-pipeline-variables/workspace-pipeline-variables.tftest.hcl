# Auto-generated Terraform test for bitbucket_workspace_pipeline_variables
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_workspace_pipeline_variables" {
  command = apply

  variables {
    workspace = "my-workspace"
    variable_uuid = "{variable-uuid}"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_workspace_pipeline_variables.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_workspace_pipeline_variables"
  }
}

run "create_workspace_pipeline_variables" {
  command = apply

  variables {
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_workspace_pipeline_variables.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_workspace_pipeline_variables"
  }
}
