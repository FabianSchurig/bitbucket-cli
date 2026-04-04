---
page_title: "bitbucket_pipeline_variables Resource - bitbucket"
subcategory: "Pipelines"
description: |-
  Manages Bitbucket pipeline-variables via the Bitbucket Cloud API.
---

# bitbucket_pipeline_variables (Resource)

Manages Bitbucket pipeline-variables via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines_config/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-variables-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-variables-variable-uuid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-variables-variable-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-variables-variable-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-variables-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| Update | `admin:pipeline:bitbucket` |
| Delete | `admin:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `variable_uuid` (String) Path parameter (auto-populated from API response).
- `key` (String) The unique name of the variable. (also computed from API response)
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API. (also computed from API response)
- `uuid` (String) The UUID identifying the variable. (also computed from API response)
- `value` (String) The value of the variable. If the variable is secured, this will be empty. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
