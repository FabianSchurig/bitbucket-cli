---
page_title: "bitbucket_workspaces Data Source - bitbucket"
subcategory: "Workspaces"
description: |-
  Reads Bitbucket workspaces via the Bitbucket Cloud API.
---

# bitbucket_workspaces (Data Source)

Reads Bitbucket workspaces via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-get) |
| List | `GET` | `/workspaces` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:workspace:bitbucket` |
| List | `read:user:bitbucket`, `read:workspace:bitbucket` |

## Example Usage

```hcl
data "bitbucket_workspaces" "example" {
}

output "workspaces_response" {
  value = data.bitbucket_workspaces.example.api_response
}
```

## Schema

### Required

### Optional
- `workspace` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `forking_mode` (String) Controls the rules for forking repositories within this workspace.
- `is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
- `is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
- `name` (String) The name of the workspace.
- `slug` (String) The short label that identifies this workspace.
- `updated_on` (String) updated_on
- `uuid` (String) The workspace's immutable id.
