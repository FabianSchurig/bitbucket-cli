data "bitbucket_pipeline_config" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_config_response" {
  value = data.bitbucket_pipeline_config.example.api_response
}
