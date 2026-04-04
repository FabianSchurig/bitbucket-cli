data "bitbucket_pipeline_caches" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_caches_response" {
  value = data.bitbucket_pipeline_caches.example.api_response
}
