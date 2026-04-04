---
page_title: "bitbucket_project_deploy_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_project_deploy_keys (Resource)

Manages Bitbucket project-deploy-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_project_deploy_keys" "example" {
  key_id = "123"
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
