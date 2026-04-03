# Auto-generated Terraform test for bitbucket_projects
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_projects" {
  command = plan

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_projects.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_projects"
  }
}

run "create_projects" {
  command = plan

  variables {
    workspace = "my-workspace"
    project_key = "PROJ"
  }

  # Resource create should produce a plan without errors
  assert {
    condition     = bitbucket_projects.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_projects"
  }
}
