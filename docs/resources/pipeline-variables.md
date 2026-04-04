---
page_title: "bitbucket_pipeline_variables Resource - bitbucket"
subcategory: ""
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

## Example Usage

```hcl
resource "bitbucket_pipeline_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  variable_uuid = "{variable-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `variable_uuid` (String) Path parameter.

### Optional
- `key` (String) The unique name of the variable. (also computed from API response)
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API. (also computed from API response)
- `uuid` (String) The UUID identifying the variable. (also computed from API response)
- `value` (String) The value of the variable. If the variable is secured, this will be empty. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
