# Auto-generated Terraform test for bitbucket_commit_statuses
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_commit_statuses" {
  command = apply

  variables {
    commit = "abc123def"
    key = "build-key"
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_commit_statuses.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_commit_statuses"
  }
}

run "create_commit_statuses" {
  command = apply

  variables {
    commit = "abc123def"
    key = "build-key"
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_commit_statuses.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_commit_statuses"
  }
}
