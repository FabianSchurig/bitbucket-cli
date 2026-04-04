data "bitbucket_snippets" "example" {
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
