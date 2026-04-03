---
page_title: "bitbucket_repo_user_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_user_permissions (Data Source)

Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repo_user_permissions" "example" {
  repo_slug = "my-repo"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}

output "repo_user_permissions_response" {
  value = data.bitbucket_repo_user_permissions.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `selected_user_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
