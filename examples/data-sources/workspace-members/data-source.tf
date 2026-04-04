data "bitbucket_workspace_members" "example" {
  workspace = "my-workspace"
}

output "workspace_members_response" {
  value = data.bitbucket_workspace_members.example.api_response
}
