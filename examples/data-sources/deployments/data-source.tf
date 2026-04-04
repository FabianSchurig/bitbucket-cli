data "bitbucket_deployments" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "deployments_response" {
  value = data.bitbucket_deployments.example.api_response
}
