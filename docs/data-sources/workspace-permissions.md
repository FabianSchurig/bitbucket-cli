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
- `workspace_created_on` (String) workspace.created_on
- `workspace_forking_mode` (String) Controls the rules for forking repositories within this workspace.
- `workspace_is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
- `workspace_is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
- `workspace_name` (String) The name of the workspace.
- `workspace_slug` (String) The short label that identifies this workspace.
- `workspace_updated_on` (String) workspace.updated_on
- `workspace_uuid` (String) The workspace's immutable id.
