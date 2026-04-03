data "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_ssh_keys_response" {
  value = data.bitbucket_pipeline_ssh_keys.example.api_response
}
