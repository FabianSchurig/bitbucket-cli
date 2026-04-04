# Auto-generated Terraform test for bitbucket_repo_runners
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_repo_runners" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    runner_uuid = "{runner-uuid}"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_repo_runners.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_repo_runners"
  }
}

run "create_repo_runners" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_repo_runners.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_repo_runners"
  }
}
