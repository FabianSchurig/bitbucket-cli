data "bitbucket_downloads" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  filename = "artifact.zip"
}

output "downloads_response" {
  value = data.bitbucket_downloads.example.api_response
}
