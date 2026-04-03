---
page_title: "bitbucket_workspace_pipeline_variables Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.
---

# bitbucket_workspace_pipeline_variables (Data Source)

Reads Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspace_pipeline_variables" "example" {
  workspace = "my-workspace"
  variable_uuid = "{variable-uuid}"
}

output "workspace_pipeline_variables_response" {
  value = data.bitbucket_workspace_pipeline_variables.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `variable_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
