data "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pipeline_uuid = "pipeline-uuid"
}

output "pipelines_response" {
  value = data.bitbucket_pipelines.example.api_response
}
