data "bitbucket_workspace_hooks" "example" {
  workspace = "my-workspace"
  uid = "webhook-uuid"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
