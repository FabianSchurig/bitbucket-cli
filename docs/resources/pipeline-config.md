---
page_title: "bitbucket_pipeline_config Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-config via the Bitbucket Cloud API.
---

# bitbucket_pipeline_config (Resource)

Manages Bitbucket pipeline-config via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## Example Usage

```hcl
resource "bitbucket_pipeline_config" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `enabled` (String) Whether Pipelines is enabled for the repository. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
