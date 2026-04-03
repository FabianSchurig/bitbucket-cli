# Auto-generated Terraform test for bitbucket_hooks
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_hooks" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    uid = "webhook-uuid"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_hooks.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_hooks"
  }
}

run "create_hooks" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    uid = "webhook-uuid"
  }

  # Resource create should produce a plan without errors
  assert {
    condition     = bitbucket_hooks.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_hooks"
  }
}
