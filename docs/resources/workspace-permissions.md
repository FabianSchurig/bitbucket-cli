---
page_title: "bitbucket_workspace_permissions Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket workspace-permissions via the Bitbucket Cloud API.
---

# bitbucket_workspace_permissions (Resource)

Manages Bitbucket workspace-permissions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

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
resource "bitbucket_workspace_permissions" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
