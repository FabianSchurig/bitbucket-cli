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

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
