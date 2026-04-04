---
page_title: "bitbucket_workspace_permissions Data Source - bitbucket"
subcategory: "Workspaces"
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
- `user` (Object) user
  Nested schema:
  - `uuid` (String) uuid
  - `created_on` (String) created_on
  - `display_name` (String) display_name

- `workspace` (Object) workspace
  Nested schema:
  - `uuid` (String) The workspace's immutable id.
  - `is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
  - `created_on` (String) created_on
  - `name` (String) The name of the workspace.
  - `slug` (String) The short label that identifies this workspace.
  - `is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
  - `forking_mode` (String) Controls the rules for forking repositories within this workspace.
  - `updated_on` (String) updated_on

