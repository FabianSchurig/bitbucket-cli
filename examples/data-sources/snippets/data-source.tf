data "bitbucket_snippets" "example" {
  workspace = "my-workspace"
  encoded_id = "snippet-id"
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
