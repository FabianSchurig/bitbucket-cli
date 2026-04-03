---
page_title: "bitbucket_reports Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket reports via the Bitbucket Cloud API.
---

# bitbucket_reports (Resource)

Manages Bitbucket reports via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  reportId = "report-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `reportId` (String) Path parameter.

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
