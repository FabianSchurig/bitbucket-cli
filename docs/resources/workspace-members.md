---
page_title: "bitbucket_workspace_members Resource - bitbucket"
subcategory: "Workspaces"
description: |-
  Manages Bitbucket workspace-members via the Bitbucket Cloud API.
---

# bitbucket_workspace_members (Resource)

Manages Bitbucket workspace-members via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/members/{member}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-members-member-get) |
| List | `GET` | `/workspaces/{workspace}/members` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-members-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:workspace:bitbucket` |
| List | `read:workspace:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_workspace_members" "example" {
  member = "{member-uuid}"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `member` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

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

