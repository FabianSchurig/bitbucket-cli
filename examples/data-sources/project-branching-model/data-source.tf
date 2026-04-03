data "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_branching_model_response" {
  value = data.bitbucket_project_branching_model.example.api_response
}
