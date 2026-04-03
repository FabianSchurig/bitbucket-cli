---
page_title: "bitbucket_issue_comments Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket issue-comments via the Bitbucket Cloud API.
---

# bitbucket_issue_comments (Resource)

Manages Bitbucket issue-comments via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_issue_comments" "example" {
  comment_id = "1"
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `comment_id` (String) Path parameter.
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
