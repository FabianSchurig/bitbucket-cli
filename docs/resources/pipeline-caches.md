---
page_title: "bitbucket_pipeline_caches Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-caches via the Bitbucket Cloud API.
---

# bitbucket_pipeline_caches (Resource)

Manages Bitbucket pipeline-caches via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/caches/{cache_uuid}/content-uri` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-caches-cache-uuid-content-uri-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines-config/caches/{cache_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-caches-cache-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/caches` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-caches-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| Delete | `write:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_caches" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  cache_uuid = "{cache-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `cache_uuid` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
