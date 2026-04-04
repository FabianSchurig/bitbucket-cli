data "bitbucket_tags" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  name = "main"
}

output "tags_response" {
  value = data.bitbucket_tags.example.api_response
}
