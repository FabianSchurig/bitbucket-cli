---
page_title: "bitbucket_pr Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pr via the Bitbucket Cloud API.
---

# bitbucket_pr (Resource)

Manages Bitbucket pr via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_pr" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pull_request_id = "1"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `pull_request_id` (String) Path parameter.

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
