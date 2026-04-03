# Auto-generated Terraform test for bitbucket_forked_repository
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "create_forked_repository" {
  command = apply

  variables {
    repo_slug = "my-repo"
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_forked_repository.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_forked_repository"
  }
}
