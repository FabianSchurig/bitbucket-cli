resource "bitbucket_project_default_reviewers" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  selected_user = "jdoe"
}
