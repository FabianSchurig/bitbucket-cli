data "bitbucket_ssh_keys" "example" {
  selected_user = "jdoe"
  key_id = "123"
}

output "ssh_keys_response" {
  value = data.bitbucket_ssh_keys.example.api_response
}
