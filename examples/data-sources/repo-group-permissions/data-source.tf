data "bitbucket_repo_group_permissions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  group_slug = "developers"
}

output "repo_group_permissions_response" {
  value = data.bitbucket_repo_group_permissions.example.api_response
}
