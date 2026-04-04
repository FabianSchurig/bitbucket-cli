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
- `user_created_on` (String) user.created_on
- `user_display_name` (String) user.display_name
- `user_uuid` (String) user.uuid
- `workspace_created_on` (String) workspace.created_on
- `workspace_forking_mode` (String) Controls the rules for forking repositories within this workspace.
- `workspace_is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
- `workspace_is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
- `workspace_name` (String) The name of the workspace.
- `workspace_slug` (String) The short label that identifies this workspace.
- `workspace_updated_on` (String) workspace.updated_on
- `workspace_uuid` (String) The workspace's immutable id.
