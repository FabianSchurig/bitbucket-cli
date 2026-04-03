---
page_title: "bitbucket_repo_group_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_group_permissions (Data Source)

Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repo_group_permissions" "example" {
  group_slug = "developers"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_group_permissions_response" {
  value = data.bitbucket_repo_group_permissions.example.api_response
}
```

## Schema

### Required
- `group_slug` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
