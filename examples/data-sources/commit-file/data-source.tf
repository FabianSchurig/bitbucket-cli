data "bitbucket_commit_file" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  commit = "abc123def"
  path = "README.md"
}

output "commit_file_response" {
  value = data.bitbucket_commit_file.example.api_response
}
