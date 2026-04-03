data "bitbucket_repo_deploy_keys" "example" {
  key_id = "123"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_deploy_keys_response" {
  value = data.bitbucket_repo_deploy_keys.example.api_response
}
