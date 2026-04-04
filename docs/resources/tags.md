---
page_title: "bitbucket_tags Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket tags via the Bitbucket Cloud API.
---

# bitbucket_tags (Resource)

Manages Bitbucket tags via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_tags" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `name` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `date` (String) The date that the tag was created, if available (also computed from API response)
- `message` (String) The message associated with the tag, if available. (also computed from API response)
- `name` (String) The name of the ref. (also computed from API response)
- `tagger_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket. (also computed from API response)
- `type` (String) type (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
