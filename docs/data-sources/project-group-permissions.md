---
page_title: "bitbucket_project_group_permissions Data Source - bitbucket"
subcategory: "Projects"
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
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_group_permissions_response" {
  value = data.bitbucket_project_group_permissions.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `group_slug` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `group` (Object) group
  Nested schema:
  - `name` (String) name
  - `slug` (String) The "sluggified" version of the group's name. This contains only ASCII
  - `full_slug` (String) The concatenation of the workspace's slug and the group's slug,

- `permission` (String) [read, write, create-repo, admin, none]
- `project` (Object) project
  Nested schema:
  - `created_on` (String) created_on
  - `has_publicly_visible_repos` (String) 
  - `uuid` (String) The project's immutable id.
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `description` (String) description
  - `is_private` (String) 
  - `updated_on` (String) updated_on

- `type` (String) type
