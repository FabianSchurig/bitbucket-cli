# Auto-generated Terraform test for bitbucket_reports
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_reports" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    commit = "abc123def"
    reportId = "report-uuid"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_reports.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_reports"
  }
}

run "create_reports" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    commit = "abc123def"
    reportId = "report-uuid"
  }

  # Resource create should produce a plan without errors
  assert {
    condition     = bitbucket_reports.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_reports"
  }
}
