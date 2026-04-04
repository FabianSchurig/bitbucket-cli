---
page_title: "bitbucket_project_default_reviewers Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_project_default_reviewers (Resource)

Manages Bitbucket project-default-reviewers via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_project_default_reviewers" "example" {
  project_key = "PROJ"
  selected_user = "jdoe"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `selected_user` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `account_id` (String) The user's Atlassian account ID.
- `account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `has_2fa_enabled` (String) has_2fa_enabled
- `is_staff` (String) is_staff
- `nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be use...
- `uuid` (String) uuid
