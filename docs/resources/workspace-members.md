---
page_title: "bitbucket_workspace_members Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket workspace-members via the Bitbucket Cloud API.
---

# bitbucket_workspace_members (Resource)

Manages Bitbucket workspace-members via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_workspace_members" "example" {
  member = "{member-uuid}"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `member` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
