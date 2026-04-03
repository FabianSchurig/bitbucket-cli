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
  workspace = "my-workspace"
  project_key = "PROJ"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
