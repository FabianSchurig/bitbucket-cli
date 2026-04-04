---
page_title: "bitbucket_pipeline_known_hosts Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.
---

# bitbucket_pipeline_known_hosts (Data Source)

Reads Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_known_hosts" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  known_host_uuid = "{known-host-uuid}"
}

output "pipeline_known_hosts_response" {
  value = data.bitbucket_pipeline_known_hosts.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `known_host_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `hostname` (String) The hostname of the known host.
- `public_key_key` (String) The base64 encoded public key.
- `public_key_key_type` (String) The type of the public key.
- `public_key_md5_fingerprint` (String) The MD5 fingerprint of the public key.
- `public_key_sha256_fingerprint` (String) The SHA-256 fingerprint of the public key.
- `uuid` (String) The UUID identifying the known host.
