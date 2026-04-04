---
page_title: "bitbucket_workspace_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-permissions via the Bitbucket Cloud API.
---

# bitbucket_workspace_permissions (Data Source)

Reads Bitbucket workspace-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/user/workspaces/{workspace}/permission` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-user-workspaces-workspace-permission-get) |
| List | `GET` | `/workspaces/{workspace}/permissions` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-permissions-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:workspace:bitbucket` |
| List | `read:workspace:bitbucket` |

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
