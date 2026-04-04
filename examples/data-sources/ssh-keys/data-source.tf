data "bitbucket_ssh_keys" "example" {
  key_id = "123"
  selected_user = "jdoe"
}

output "ssh_keys_response" {
  value = data.bitbucket_ssh_keys.example.api_response
}
