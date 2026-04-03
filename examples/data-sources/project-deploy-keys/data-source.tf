data "bitbucket_project_deploy_keys" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  key_id = "123"
}

output "project_deploy_keys_response" {
  value = data.bitbucket_project_deploy_keys.example.api_response
}
