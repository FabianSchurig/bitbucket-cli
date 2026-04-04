# Auto-generated Terraform test for bitbucket_pipeline_known_hosts
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_pipeline_known_hosts" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    known_host_uuid = "{known-host-uuid}"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_pipeline_known_hosts.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_pipeline_known_hosts"
  }
}

run "create_pipeline_known_hosts" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    known_host_uuid = "{known-host-uuid}"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_pipeline_known_hosts.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_pipeline_known_hosts"
  }
}
