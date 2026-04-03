data "bitbucket_branch_restrictions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  id = "1"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
