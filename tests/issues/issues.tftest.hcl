# Auto-generated Terraform test for bitbucket_issues
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_issues" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    issue_id = "1"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_issues.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_issues"
  }
}

run "create_issues" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    issue_id = "1"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_issues.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_issues"
  }
}
