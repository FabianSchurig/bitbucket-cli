---
page_title: "bitbucket_workspace_pipeline_variables Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.
---

# bitbucket_workspace_pipeline_variables (Resource)

Manages Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_workspace_pipeline_variables" "example" {
  workspace = "my-workspace"
  variable_uuid = "{variable-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `variable_uuid` (String) Path parameter.

### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `key` (String) The unique name of the variable.
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.
- `uuid` (String) The UUID identifying the variable.
- `value` (String) The value of the variable. If the variable is secured, this will be empty.
