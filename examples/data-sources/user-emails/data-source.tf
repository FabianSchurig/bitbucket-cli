data "bitbucket_user_emails" "example" {
  email = "user@example.com"
}

output "user_emails_response" {
  value = data.bitbucket_user_emails.example.api_response
}
