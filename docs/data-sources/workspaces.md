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
