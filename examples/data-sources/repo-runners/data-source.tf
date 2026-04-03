data "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  runner_uuid = "{runner-uuid}"
}

output "repo_runners_response" {
  value = data.bitbucket_repo_runners.example.api_response
}
