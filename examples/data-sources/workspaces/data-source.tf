data "bitbucket_workspaces" "example" {
}

output "workspaces_response" {
  value = data.bitbucket_workspaces.example.api_response
}
