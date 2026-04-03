data "bitbucket_workspace_runners" "example" {
  workspace = "my-workspace"
  runner_uuid = "{runner-uuid}"
}

output "workspace_runners_response" {
  value = data.bitbucket_workspace_runners.example.api_response
}
