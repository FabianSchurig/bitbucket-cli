# Auto-generated Terraform test for bitbucket_issue_comments
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_issue_comments" {
  command = apply

  variables {
    comment_id = "1"
    issue_id = "1"
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_issue_comments.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_issue_comments"
  }
}

run "create_issue_comments" {
  command = apply

  variables {
    comment_id = "1"
    issue_id = "1"
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_issue_comments.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_issue_comments"
  }
}
