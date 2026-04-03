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
  workspace = "my-workspace"
  repo_slug = "my-repo"
  selected_user_id = "{user-uuid}"
}

output "repo_user_permissions_response" {
  value = data.bitbucket_repo_user_permissions.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `selected_user_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
