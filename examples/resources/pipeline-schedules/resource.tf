resource "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  cron_pattern = "example-value"
  target = "example-value"
}
