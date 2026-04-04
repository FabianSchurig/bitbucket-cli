---
page_title: "bitbucket_search Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket search via the Bitbucket Cloud API.
---

# bitbucket_search (Resource)

Manages Bitbucket search via the Bitbucket Cloud API.

## CRUD Operations
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_search" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
