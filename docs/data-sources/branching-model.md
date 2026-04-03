---
page_title: "bitbucket_branching_model Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Data Source)

Reads Bitbucket branching-model via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_branching_model" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branching_model_response" {
  value = data.bitbucket_branching_model.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
