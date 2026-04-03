data "bitbucket_tags" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "tags_response" {
  value = data.bitbucket_tags.example.api_response
}
