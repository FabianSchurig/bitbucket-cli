---
page_title: "bitbucket_project_deploy_keys Resource - bitbucket"
subcategory: ""
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
  key_id = "123"
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
