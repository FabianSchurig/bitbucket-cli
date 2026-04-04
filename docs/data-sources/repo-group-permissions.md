---
page_title: "bitbucket_repo_group_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_group_permissions (Data Source)

Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-group-slug-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

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
