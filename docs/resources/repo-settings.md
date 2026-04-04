---
page_title: "bitbucket_repo_settings Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket repo-settings via the Bitbucket Cloud API.
---

# bitbucket_repo_settings (Resource)

Manages Bitbucket repo-settings via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## Example Usage

```hcl
resource "bitbucket_repo_settings" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `type` (String) type
