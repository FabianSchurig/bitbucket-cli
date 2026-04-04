---
page_title: "bitbucket_workspaces Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspaces via the Bitbucket Cloud API.
---

# bitbucket_workspaces (Data Source)

Reads Bitbucket workspaces via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspaces" "example" {
  workspace = "my-workspace"
}

output "workspaces_response" {
  value = data.bitbucket_workspaces.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `forking_mode` (String) Controls the rules for forking repositories within this workspace.
- `is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
- `is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
- `name` (String) The name of the workspace.
- `slug` (String) The short label that identifies this workspace.
- `updated_on` (String) updated_on
- `uuid` (String) The workspace's immutable id.
