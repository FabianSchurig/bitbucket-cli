data "bitbucket_addon" "example" {
}

output "addon_response" {
  value = data.bitbucket_addon.example.api_response
}
