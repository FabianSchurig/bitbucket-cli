data "bitbucket_downloads" "example" {
  filename = "artifact.zip"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "downloads_response" {
  value = data.bitbucket_downloads.example.api_response
}
