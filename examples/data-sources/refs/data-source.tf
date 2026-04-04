data "bitbucket_refs" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
