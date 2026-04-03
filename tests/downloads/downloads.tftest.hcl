# Auto-generated Terraform test for bitbucket_downloads
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_downloads" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    filename = "artifact.zip"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_downloads.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_downloads"
  }
}

run "create_downloads" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    filename = "artifact.zip"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_downloads.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_downloads"
  }
}
