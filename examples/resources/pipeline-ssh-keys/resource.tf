resource "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
