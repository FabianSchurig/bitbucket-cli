data "bitbucket_commit_file" "example" {
  commit = "abc123def"
  path = "README.md"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commit_file_response" {
  value = data.bitbucket_commit_file.example.api_response
}
