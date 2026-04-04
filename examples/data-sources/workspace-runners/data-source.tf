data "bitbucket_workspace_runners" "example" {
  workspace = "my-workspace"
}

output "workspace_runners_response" {
  value = data.bitbucket_workspace_runners.example.api_response
}
