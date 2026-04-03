---
page_title: "bitbucket_commit_statuses Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket commit-statuses via the Bitbucket Cloud API.
---

# bitbucket_commit_statuses (Resource)

Manages Bitbucket commit-statuses via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_commit_statuses" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  key = "build-key"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `key` (String) Path parameter.

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
