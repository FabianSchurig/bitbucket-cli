data "bitbucket_refs" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  name = "main"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
