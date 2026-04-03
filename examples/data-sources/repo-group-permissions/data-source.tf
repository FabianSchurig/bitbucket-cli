data "bitbucket_repo_group_permissions" "example" {
  group_slug = "developers"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_group_permissions_response" {
  value = data.bitbucket_repo_group_permissions.example.api_response
}
