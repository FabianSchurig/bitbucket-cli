---
page_title: "bitbucket_workspace_members Data Source - bitbucket"
subcategory: ""
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
  member = "{member-uuid}"
  workspace = "my-workspace"
}

output "workspace_members_response" {
  value = data.bitbucket_workspace_members.example.api_response
}
```

## Schema

### Required
- `member` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
