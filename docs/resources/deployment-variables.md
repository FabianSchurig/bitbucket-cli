---
page_title: "bitbucket_deployment_variables Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket deployment-variables via the Bitbucket Cloud API.
---

# bitbucket_deployment_variables (Resource)

Manages Bitbucket deployment-variables via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-deployments-config-environments-environment-uuid-variables-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-deployments-config-environments-environment-uuid-variables-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-deployments-config-environments-environment-uuid-variables-variable-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-deployments-config-environments-environment-uuid-variables-variable-uuid-delete) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| Update | `admin:pipeline:bitbucket` |
| Delete | `admin:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_deployment_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `environment_uuid` (String) Path parameter.

### Optional
- `key` (String) The unique name of the variable. (also computed from API response)
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API. (also computed from API response)
- `uuid` (String) The UUID identifying the variable. (also computed from API response)
- `value` (String) The value of the variable. If the variable is secured, this will be empty. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
