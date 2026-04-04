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
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "projects_response" {
  value = data.bitbucket_projects.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `description` (String) description
- `has_publicly_visible_repos` (String) 
- `is_private` (String) 
- `key` (String) The project's key.
- `name` (String) The name of the project.
- `updated_on` (String) updated_on
- `uuid` (String) The project's immutable id.
