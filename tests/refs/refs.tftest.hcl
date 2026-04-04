# Auto-generated Terraform test for bitbucket_refs
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_refs" {
  command = apply

  variables {
    repo_slug = "my-repo"
    workspace = "my-workspace"
    name = "main"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_refs.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_refs"
  }
}

run "create_refs" {
  command = apply

  variables {
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_refs.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_refs"
  }
}
