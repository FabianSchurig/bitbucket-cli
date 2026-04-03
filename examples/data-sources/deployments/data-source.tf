data "bitbucket_deployments" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}

output "deployments_response" {
  value = data.bitbucket_deployments.example.api_response
}
