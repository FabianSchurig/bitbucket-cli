# Auto-generated Terraform test for bitbucket_pipeline_schedules
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_pipeline_schedules" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    schedule_uuid = "{schedule-uuid}"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_pipeline_schedules.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_pipeline_schedules"
  }
}

run "create_pipeline_schedules" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_pipeline_schedules.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_pipeline_schedules"
  }
}
