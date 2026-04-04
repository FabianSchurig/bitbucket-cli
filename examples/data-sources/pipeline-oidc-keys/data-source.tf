data "bitbucket_pipeline_oidc_keys" "example" {
  workspace = "my-workspace"
}

output "pipeline_oidc_keys_response" {
  value = data.bitbucket_pipeline_oidc_keys.example.api_response
}
