---
page_title: "bitbucket_issues Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket issues via the Bitbucket Cloud API.
---

# bitbucket_issues (Resource)

Manages Bitbucket issues via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_issues" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  issue_id = "1"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `issue_id` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
