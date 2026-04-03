---
page_title: "bitbucket_addon Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket addon via the Bitbucket Cloud API.
---

# bitbucket_addon (Resource)

Manages Bitbucket addon via the Bitbucket Cloud API.

## CRUD Operations
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_addon" "example" {
}
```

## Schema

### Required

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
