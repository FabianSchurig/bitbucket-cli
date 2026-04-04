data "bitbucket_pipeline_known_hosts" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  known_host_uuid = "{known-host-uuid}"
}

output "pipeline_known_hosts_response" {
  value = data.bitbucket_pipeline_known_hosts.example.api_response
}
