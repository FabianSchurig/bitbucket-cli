data "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipelines_response" {
  value = data.bitbucket_pipelines.example.api_response
}
