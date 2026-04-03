---
page_title: "bitbucket_projects Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket projects via the Bitbucket Cloud API.
---

# bitbucket_projects (Data Source)

Reads Bitbucket projects via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_projects" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "projects_response" {
  value = data.bitbucket_projects.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
