data "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_schedules_response" {
  value = data.bitbucket_pipeline_schedules.example.api_response
}
