data "bitbucket_project_default_reviewers" "example" {
  project_key = "PROJ"
  selected_user = "jdoe"
  workspace = "my-workspace"
}

output "project_default_reviewers_response" {
  value = data.bitbucket_project_default_reviewers.example.api_response
}
