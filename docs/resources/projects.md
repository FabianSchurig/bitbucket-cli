---
page_title: "bitbucket_projects Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket projects via the Bitbucket Cloud API.
---

# bitbucket_projects (Resource)

Manages Bitbucket projects via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_projects" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `description` (String) description
- `has_publicly_visible_repos` (String) 
- `is_private` (String) 
- `key` (String) The project's key.
- `name` (String) The name of the project.
- `updated_on` (String) updated_on
- `uuid` (String) The project's immutable id.
