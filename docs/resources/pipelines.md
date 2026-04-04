---
page_title: "bitbucket_pipelines Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipelines via the Bitbucket Cloud API.
---

# bitbucket_pipelines (Resource)

Manages Bitbucket pipelines via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-pipeline-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:pipeline:bitbucket`, `write:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `pipeline_uuid` (String) Path parameter (auto-populated from API response).
- `build_number` (String) The build number of the pipeline. (also computed from API response)
- `build_seconds_used` (String) The number of build seconds used by this pipeline. (also computed from API response)
- `completed_on` (String) The timestamp when the Pipeline was completed. This is not set if the pipeline is still in progress. (also computed from API response)
- `configuration_sources` (List of Object) An ordered list of sources of the pipeline configuration (also computed from API response)
  Nested schema:
  - `source` (String) Identifier of the configuration source
  - `uri` (String) Link to the configuration source view or its immediate content

- `uuid` (String) The UUID identifying the pipeline. (also computed from API response)
- `variables` (List of Object) The variables for the pipeline. (also computed from API response)
  Nested schema:
  - `uuid` (String) The UUID identifying the variable.
  - `key` (String) The unique name of the variable.
  - `value` (String) The value of the variable. If the variable is secured, this will be empty.
  - `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the pipeline was created.
