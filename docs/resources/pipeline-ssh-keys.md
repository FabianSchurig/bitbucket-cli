---
page_title: "bitbucket_pipeline_ssh_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_ssh_keys (Resource)

Manages Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported

## Example Usage

```hcl
resource "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `private_key` (String) The SSH private key. This value will be empty when retrieving the SSH key pair. (also computed from API response)
- `public_key` (String) The SSH public key. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
