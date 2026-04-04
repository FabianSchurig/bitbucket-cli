data "bitbucket_hook_types" "example" {
  subject_type = "repository"
}

output "hook_types_response" {
  value = data.bitbucket_hook_types.example.api_response
}
