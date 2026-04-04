---
page_title: "bitbucket_project_branching_model Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-branching-model via the Bitbucket Cloud API.
---

# bitbucket_project_branching_model (Data Source)

Reads Bitbucket project-branching-model via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_branching_model_response" {
  value = data.bitbucket_project_branching_model.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
