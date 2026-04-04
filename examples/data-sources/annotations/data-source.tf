data "bitbucket_annotations" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
}

output "annotations_response" {
  value = data.bitbucket_annotations.example.api_response
}
