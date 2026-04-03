---
page_title: "bitbucket_project_group_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_group_permissions (Data Source)

Reads Bitbucket project-group-permissions via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_project_group_permissions" "example" {
  group_slug = "developers"
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_group_permissions_response" {
  value = data.bitbucket_project_group_permissions.example.api_response
}
```

## Schema

### Required
- `group_slug` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
