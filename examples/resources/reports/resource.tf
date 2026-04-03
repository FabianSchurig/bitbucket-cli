resource "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  reportId = "report-uuid"
}
