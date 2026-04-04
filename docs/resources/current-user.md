---
page_title: "bitbucket_current_user Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket current-user via the Bitbucket Cloud API.
---

# bitbucket_current_user (Resource)

Manages Bitbucket current-user via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported

## Example Usage

```hcl
resource "bitbucket_current_user" "example" {
}
```

## Schema

### Required

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
