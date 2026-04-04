data "bitbucket_refs" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  name = "main"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
