resource "bitbucket_deployments" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}
