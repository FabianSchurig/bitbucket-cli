data "bitbucket_repo_user_permissions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_user_permissions_response" {
  value = data.bitbucket_repo_user_permissions.example.api_response
}
