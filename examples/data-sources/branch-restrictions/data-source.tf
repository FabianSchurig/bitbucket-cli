data "bitbucket_branch_restrictions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  param_id = "1"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
