---
page_title: "bitbucket_workspace_members Data Source - bitbucket"
subcategory: "Workspaces"
description: |-
  Reads Bitbucket workspace-members via the Bitbucket Cloud API.
---

# bitbucket_workspace_members (Data Source)

Reads Bitbucket workspace-members via the Bitbucket Cloud API.

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
data "bitbucket_workspace_members" "example" {
  workspace = "my-workspace"
}

output "workspace_members_response" {
  value = data.bitbucket_workspace_members.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `member` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
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

