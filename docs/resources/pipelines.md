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
  pipeline_uuid = "pipeline-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `pipeline_uuid` (String) Path parameter.

### Optional
- `build_number` (String) The build number of the pipeline. (also computed from API response)
- `build_seconds_used` (String) The number of build seconds used by this pipeline. (also computed from API response)
- `completed_on` (String) The timestamp when the Pipeline was completed. This is not set if the pipeline is still in progress. (also computed from API response)
- `uuid` (String) The UUID identifying the pipeline. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the pipeline was created.
