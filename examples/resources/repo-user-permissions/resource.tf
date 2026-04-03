resource "bitbucket_repo_user_permissions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  selected_user_id = "{user-uuid}"
}
