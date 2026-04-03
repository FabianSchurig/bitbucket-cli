resource "bitbucket_repo_deploy_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  key_id = "123"
}
