# Auto-generated Terraform test for bitbucket_repos
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_repos" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_repos.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_repos"
  }
}

run "create_repos" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_repos.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_repos"
  }
}
