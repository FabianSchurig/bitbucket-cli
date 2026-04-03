data "bitbucket_project_default_reviewers" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  selected_user = "jdoe"
}

output "project_default_reviewers_response" {
  value = data.bitbucket_project_default_reviewers.example.api_response
}
