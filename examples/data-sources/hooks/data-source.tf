data "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "hooks_response" {
  value = data.bitbucket_hooks.example.api_response
}
