---
page_title: "bitbucket_project_user_permissions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_user_permissions (Data Source)

Reads Bitbucket project-user-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-selected-user-id-get) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_user_permissions" "example" {
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}

output "project_user_permissions_response" {
  value = data.bitbucket_project_user_permissions.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `selected_user_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
