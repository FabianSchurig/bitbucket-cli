data "bitbucket_branching_model" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branching_model_response" {
  value = data.bitbucket_branching_model.example.api_response
}
