data "bitbucket_branching_model" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "branching_model_response" {
  value = data.bitbucket_branching_model.example.api_response
}
