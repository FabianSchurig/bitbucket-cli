data "bitbucket_pipeline_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  variable_uuid = "{variable-uuid}"
}

output "pipeline_variables_response" {
  value = data.bitbucket_pipeline_variables.example.api_response
}
