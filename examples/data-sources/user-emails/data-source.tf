data "bitbucket_user_emails" "example" {
}

output "user_emails_response" {
  value = data.bitbucket_user_emails.example.api_response
}
