data "bitbucket_project_deploy_keys" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_deploy_keys_response" {
  value = data.bitbucket_project_deploy_keys.example.api_response
}
