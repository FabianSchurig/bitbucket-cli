data "bitbucket_current_user" "example" {
}

output "current_user_response" {
  value = data.bitbucket_current_user.example.api_response
}
