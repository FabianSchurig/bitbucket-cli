---
page_title: "bitbucket_repo_runners Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket repo-runners via the Bitbucket Cloud API.
---

# bitbucket_repo_runners (Resource)

Manages Bitbucket repo-runners via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  runner_uuid = "{runner-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
