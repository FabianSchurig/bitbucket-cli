---
page_title: "bitbucket_workspace_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-permissions via the Bitbucket Cloud API.
---

# bitbucket_workspace_permissions (Data Source)

Reads Bitbucket workspace-permissions via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspace_permissions" "example" {
  workspace = "my-workspace"
}

output "workspace_permissions_response" {
  value = data.bitbucket_workspace_permissions.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
