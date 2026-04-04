data "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  uid = "webhook-uuid"
}

output "hooks_response" {
  value = data.bitbucket_hooks.example.api_response
}
