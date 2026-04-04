data "bitbucket_deployment_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
  variable_uuid = "{variable-uuid}"
}

output "deployment_variables_response" {
  value = data.bitbucket_deployment_variables.example.api_response
}
