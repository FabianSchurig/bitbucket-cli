data "bitbucket_search" "example" {
  workspace = "my-workspace"
}

output "search_response" {
  value = data.bitbucket_search.example.api_response
}
