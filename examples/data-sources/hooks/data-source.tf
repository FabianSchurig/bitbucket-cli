data "bitbucket_hooks" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  uid = "webhook-uuid"
}

output "hooks_response" {
  value = data.bitbucket_hooks.example.api_response
}
