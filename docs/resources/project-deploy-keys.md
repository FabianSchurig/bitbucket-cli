---
page_title: "bitbucket_project_deploy_keys Resource - bitbucket"
subcategory: "Deployments"
description: |-
  Manages Bitbucket project-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_project_deploy_keys (Resource)

Manages Bitbucket project-deploy-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-post) |
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-key-id-get) |
| Delete | `DELETE` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-key-id-delete) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:project:bitbucket`, `write:ssh-key:bitbucket` |
| Read | `admin:project:bitbucket` |
| Delete | `admin:project:bitbucket`, `delete:ssh-key:bitbucket` |
| List | `admin:project:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_project_deploy_keys" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `key_id` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the deploy key (if present)
- `created_by` (Object) created_by
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `key` (String) The deploy key value.
- `label` (String) The user-defined label for the deploy key
- `last_used` (String) last_used
- `project` (Object) project
  Nested schema:
  - `created_on` (String) created_on
  - `description` (String) description
  - `has_publicly_visible_repos` (String) 
  - `is_private` (String) 
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `updated_on` (String) updated_on
  - `uuid` (String) The project's immutable id.

