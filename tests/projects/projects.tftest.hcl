# Auto-generated Terraform test for bitbucket_projects
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_projects" {
  command = apply

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_projects.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_projects"
  }
}

run "create_projects" {
  command = apply

  variables {
    workspace = "my-workspace"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_projects.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_projects"
  }
}
