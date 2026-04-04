data "bitbucket_users" "example" {
  selected_user = "jdoe"
}

output "users_response" {
  value = data.bitbucket_users.example.api_response
}
