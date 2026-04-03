---
page_title: "bitbucket_downloads Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket downloads via the Bitbucket Cloud API.
---

# bitbucket_downloads (Resource)

Manages Bitbucket downloads via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_downloads" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  filename = "artifact.zip"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `filename` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
