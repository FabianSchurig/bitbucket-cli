# Auto-generated Terraform test for bitbucket_annotations
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

mock_provider "bitbucket" {}

run "read_annotations" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    commit = "abc123def"
    report_id = "report-uuid"
    annotation_id = "{annotation-id}"
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.bitbucket_annotations.test.id != ""
    error_message = "Expected non-empty id for data source bitbucket_annotations"
  }
}

run "create_annotations" {
  command = apply

  variables {
    workspace = "my-workspace"
    repo_slug = "my-repo"
    commit = "abc123def"
    report_id = "report-uuid"
    annotation_id = "{annotation-id}"
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = bitbucket_annotations.test.id != ""
    error_message = "Expected non-empty id for resource bitbucket_annotations"
  }
}
