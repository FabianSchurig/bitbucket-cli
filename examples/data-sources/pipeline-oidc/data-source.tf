data "bitbucket_pipeline_oidc" "example" {
  workspace = "my-workspace"
}

output "pipeline_oidc_response" {
  value = data.bitbucket_pipeline_oidc.example.api_response
}
