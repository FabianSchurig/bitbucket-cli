data "bitbucket_properties" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  app_key = "my-app"
  property_name = "my-property"
}

output "properties_response" {
  value = data.bitbucket_properties.example.api_response
}
