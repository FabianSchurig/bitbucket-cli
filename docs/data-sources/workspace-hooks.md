---
page_title: "bitbucket_workspace_hooks Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.
---

# bitbucket_workspace_hooks (Data Source)

Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspace_hooks" "example" {
  uid = "webhook-uuid"
  workspace = "my-workspace"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
```

## Schema

### Required
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
