data "bitbucket_workspace_pipeline_variables" "example" {
  workspace = "my-workspace"
}

output "workspace_pipeline_variables_response" {
  value = data.bitbucket_workspace_pipeline_variables.example.api_response
}
