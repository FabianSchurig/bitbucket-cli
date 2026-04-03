resource "bitbucket_project_deploy_keys" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  key_id = "123"
}
