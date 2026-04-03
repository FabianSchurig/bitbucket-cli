# Auto-generated Terraform test for bitbucket_hooks
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_hooks" {
  command = apply

  variables {
    repo_slug = "my-repo"
    uid = "webhook-uuid"
    workspace = "my-workspace"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_hooks.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_hooks"
  }
}

run "create_hooks" {
  command = apply

  variables {
    repo_slug = "my-repo"
    uid = "webhook-uuid"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_hooks.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_hooks"
  }
}
