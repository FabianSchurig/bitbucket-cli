# Auto-generated Terraform test for bitbucket_default_reviewers
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_default_reviewers" {
  command = apply

  variables {
    repo_slug = "my-repo"
    target_username = "jdoe"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_default_reviewers.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_default_reviewers"
  }
}

run "create_default_reviewers" {
  command = apply

  variables {
    repo_slug = "my-repo"
    target_username = "jdoe"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_default_reviewers.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_default_reviewers"
  }
}
