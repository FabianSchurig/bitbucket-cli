data "bitbucket_workspace_pipeline_variables" "example" {
  workspace = "my-workspace"
  variable_uuid = "{variable-uuid}"
}

output "workspace_pipeline_variables_response" {
  value = data.bitbucket_workspace_pipeline_variables.example.api_response
}
