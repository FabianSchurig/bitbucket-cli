---
page_title: "bitbucket_workspace_hooks Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.
---

# bitbucket_workspace_hooks (Data Source)

Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-uid-get) |
| List | `GET` | `/workspaces/{workspace}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:webhook:bitbucket` |
| List | `read:webhook:bitbucket` |

## Example Usage

```hcl
data "bitbucket_workspace_hooks" "example" {
  uid = "webhook-uuid"
  workspace = "my-workspace"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
```

## Schema

### Required
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
