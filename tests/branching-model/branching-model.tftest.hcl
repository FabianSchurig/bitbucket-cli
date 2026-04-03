# Auto-generated Terraform test for bitbucket_branching_model
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_branching_model" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_branching_model.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_branching_model"
  }
}
