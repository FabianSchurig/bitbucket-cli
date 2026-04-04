data "bitbucket_gpg_keys" "example" {
  selected_user = "jdoe"
  fingerprint = "AA:BB:CC:DD"
}

output "gpg_keys_response" {
  value = data.bitbucket_gpg_keys.example.api_response
}
