data "bitbucket_projects" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "projects_response" {
  value = data.bitbucket_projects.example.api_response
}
