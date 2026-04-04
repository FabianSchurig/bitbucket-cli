data "bitbucket_downloads" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  filename = "artifact.zip"
}

output "downloads_response" {
  value = data.bitbucket_downloads.example.api_response
}
