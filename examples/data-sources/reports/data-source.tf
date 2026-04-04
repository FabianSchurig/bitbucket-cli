data "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
}

output "reports_response" {
  value = data.bitbucket_reports.example.api_response
}
