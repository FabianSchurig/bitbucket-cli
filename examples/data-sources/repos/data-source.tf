data "bitbucket_repos" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
}
