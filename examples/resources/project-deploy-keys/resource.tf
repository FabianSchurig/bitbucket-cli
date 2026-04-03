resource "bitbucket_project_deploy_keys" "example" {
  key_id = "123"
  project_key = "PROJ"
  workspace = "my-workspace"
}
