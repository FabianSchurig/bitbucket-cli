data "bitbucket_repos" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
}
