# Auto-generated Terraform test for bitbucket_pipelines
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_pipelines" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    pipeline_uuid = "pipeline-uuid"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_pipelines.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_pipelines"
  }
}

run "create_pipelines" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    pipeline_uuid = "pipeline-uuid"
  }

  # Resource create should produce a plan without errors
  assert {
    condition     = bitbucket_pipelines.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_pipelines"
  }
}
