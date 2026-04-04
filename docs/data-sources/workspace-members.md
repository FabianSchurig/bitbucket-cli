---
page_title: "bitbucket_workspace_members Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-members via the Bitbucket Cloud API.
---

# bitbucket_workspace_members (Data Source)

Reads Bitbucket workspace-members via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspace_members" "example" {
  member = "{member-uuid}"
  workspace = "my-workspace"
}

output "workspace_members_response" {
  value = data.bitbucket_workspace_members.example.api_response
}
```

## Schema

### Required
- `member` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
