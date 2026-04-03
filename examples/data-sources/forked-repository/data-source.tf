data "bitbucket_forked_repository" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "forked_repository_response" {
  value = data.bitbucket_forked_repository.example.api_response
}
