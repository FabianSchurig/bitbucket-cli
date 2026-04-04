---
page_title: "bitbucket_search Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket search via the Bitbucket Cloud API.
---

# bitbucket_search (Data Source)

Reads Bitbucket search via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_search" "example" {
  workspace = "my-workspace"
}

output "search_response" {
  value = data.bitbucket_search.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
