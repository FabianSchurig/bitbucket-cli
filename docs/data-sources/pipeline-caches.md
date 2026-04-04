---
page_title: "bitbucket_pipeline_caches Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-caches via the Bitbucket Cloud API.
---

# bitbucket_pipeline_caches (Data Source)

Reads Bitbucket pipeline-caches via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_caches" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  cache_uuid = "{cache-uuid}"
}

output "pipeline_caches_response" {
  value = data.bitbucket_pipeline_caches.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `cache_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
