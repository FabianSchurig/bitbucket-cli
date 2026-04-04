data "bitbucket_gpg_keys" "example" {
  selected_user = "jdoe"
}

output "gpg_keys_response" {
  value = data.bitbucket_gpg_keys.example.api_response
}
