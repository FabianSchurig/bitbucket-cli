data "bitbucket_repo_deploy_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  key_id = "123"
}

output "repo_deploy_keys_response" {
  value = data.bitbucket_repo_deploy_keys.example.api_response
}
