data "bitbucket_workspace_hooks" "example" {
  uid = "webhook-uuid"
  workspace = "my-workspace"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
