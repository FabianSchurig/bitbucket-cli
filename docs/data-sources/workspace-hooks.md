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
  workspace = "my-workspace"
  uid = "webhook-uuid"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `uid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
