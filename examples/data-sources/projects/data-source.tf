data "bitbucket_projects" "example" {
  workspace = "my-workspace"
}

output "projects_response" {
  value = data.bitbucket_projects.example.api_response
}
