# Auto-generated Terraform test for bitbucket_pr_comments
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_pr_comments" {
  command = apply

  variables {
    pull_request_id = "1"
    repo_slug = "my-repo"
    workspace = "my-workspace"
    comment_id = "1"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_pr_comments.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_pr_comments"
  }
}

run "create_pr_comments" {
  command = apply

  variables {
    pull_request_id = "1"
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_pr_comments.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_pr_comments"
  }
}
