resource "bitbucket_branch_restrictions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  branch_match_kind = "glob"
  kind = "require_approvals_to_merge"
  pattern = "main"
}
