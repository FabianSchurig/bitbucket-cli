resource "bitbucket_repo_deploy_keys" "example" {
  key_id = "123"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
