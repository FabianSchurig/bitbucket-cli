# Auto-generated Terraform test for bitbucket_properties
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_properties" {
  command = plan

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    app_key = "my-app"
    property_name = "my-property"
  }

  # Data source read should produce a plan without errors
  assert {
    condition     = data.bitbucket_properties.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_properties"
  }
}
