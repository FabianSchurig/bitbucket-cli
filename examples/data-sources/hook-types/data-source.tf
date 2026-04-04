data "bitbucket_hook_types" "example" {
}

output "hook_types_response" {
  value = data.bitbucket_hook_types.example.api_response
}
