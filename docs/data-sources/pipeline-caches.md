---
page_title: "bitbucket_pipeline_caches Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket pipeline-caches via the Bitbucket Cloud API.
---

# bitbucket_pipeline_caches (Data Source)

Reads Bitbucket pipeline-caches via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/caches/{cache_uuid}/content-uri` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-caches-cache-uuid-content-uri-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/caches` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-caches-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_caches" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_caches_response" {
  value = data.bitbucket_pipeline_caches.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `cache_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `uri` (String) The uri for pipeline cache content.
