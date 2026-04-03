resource "bitbucket_branch_restrictions" "example" {
  param_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
