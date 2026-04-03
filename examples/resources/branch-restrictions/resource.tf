resource "bitbucket_branch_restrictions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  param_id = "1"
}
