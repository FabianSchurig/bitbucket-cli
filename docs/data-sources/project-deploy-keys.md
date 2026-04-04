---
page_title: "bitbucket_project_deploy_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_project_deploy_keys (Data Source)

Reads Bitbucket project-deploy-keys via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-key-id-get) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-workspaces-workspace-projects-project-key-deploy-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:project:bitbucket` |
| List | `admin:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_deploy_keys" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
  key_id = "123"
}

output "project_deploy_keys_response" {
  value = data.bitbucket_project_deploy_keys.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.
- `key_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the deploy key (if present)
- `key` (String) The deploy key value.
- `label` (String) The user-defined label for the deploy key
- `last_used` (String) last_used
- `project_created_on` (String) project.created_on
- `project_description` (String) project.description
- `project_has_publicly_visible_repos` (String) 
- `project_is_private` (String) 
- `project_key` (String) The project's key.
- `project_name` (String) The name of the project.
- `project_updated_on` (String) project.updated_on
- `project_uuid` (String) The project's immutable id.
