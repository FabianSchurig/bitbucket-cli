---
page_title: "bitbucket_workspaces Data Source - bitbucket"
subcategory: ""
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
  workspace = "my-workspace"
}

output "workspaces_response" {
  value = data.bitbucket_workspaces.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
