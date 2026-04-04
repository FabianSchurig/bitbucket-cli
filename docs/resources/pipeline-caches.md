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

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `uri` (String) The uri for pipeline cache content.
