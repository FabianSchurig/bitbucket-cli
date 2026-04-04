---
page_title: "bitbucket_repo_user_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_user_permissions (Data Source)

Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-selected-user-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

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
