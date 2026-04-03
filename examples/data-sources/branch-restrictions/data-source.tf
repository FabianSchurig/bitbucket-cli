data "bitbucket_branch_restrictions" "example" {
  param_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
