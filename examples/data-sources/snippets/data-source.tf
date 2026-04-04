data "bitbucket_snippets" "example" {
  encoded_id = "snippet-id"
  workspace = "my-workspace"
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
