# Auto-generated Terraform test for bitbucket_pipeline_ssh_keys
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_pipeline_ssh_keys" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_pipeline_ssh_keys.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_pipeline_ssh_keys"
  }
}
