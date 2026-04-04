---
page_title: "bitbucket_pipeline_variables Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-variables via the Bitbucket Cloud API.
---

# bitbucket_pipeline_variables (Data Source)

Reads Bitbucket pipeline-variables via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  variable_uuid = "{variable-uuid}"
}

output "pipeline_variables_response" {
  value = data.bitbucket_pipeline_variables.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `variable_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `key` (String) The unique name of the variable.
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.
- `uuid` (String) The UUID identifying the variable.
- `value` (String) The value of the variable. If the variable is secured, this will be empty.
