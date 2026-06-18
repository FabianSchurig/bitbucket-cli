data "bitbucket_addon" "example" {
  addon_key = "example-value"
}

output "addon_response" {
  value = data.bitbucket_addon.example.api_response
}
