---
page_title: "bitbucket_project_group_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_group_permissions (Data Source)

Reads Bitbucket project-group-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-group-slug-get) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_group_permissions" "example" {
  group_slug = "developers"
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_group_permissions_response" {
  value = data.bitbucket_project_group_permissions.example.api_response
}
```

## Schema

### Required
- `group_slug` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
