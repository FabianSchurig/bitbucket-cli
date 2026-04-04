---
page_title: "bitbucket_pipeline_known_hosts Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.
---

# bitbucket_pipeline_known_hosts (Resource)

Manages Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_pipeline_known_hosts" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  known_host_uuid = "{known-host-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `known_host_uuid` (String) Path parameter.

### Optional
- `hostname` (String) The hostname of the known host. (also computed from API response)
- `public_key_key` (String) The base64 encoded public key. (also computed from API response)
- `public_key_key_type` (String) The type of the public key. (also computed from API response)
- `public_key_md5_fingerprint` (String) The MD5 fingerprint of the public key. (also computed from API response)
- `public_key_sha256_fingerprint` (String) The SHA-256 fingerprint of the public key. (also computed from API response)
- `uuid` (String) The UUID identifying the known host. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
