data "bitbucket_workspaces" "example" {
  workspace = "my-workspace"
}

output "workspaces_response" {
  value = data.bitbucket_workspaces.example.api_response
}
