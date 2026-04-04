---
page_title: "bitbucket_addon Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket addon via the Bitbucket Cloud API.
---

# bitbucket_addon (Data Source)

Reads Bitbucket addon via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_addon" "example" {
}

output "addon_response" {
  value = data.bitbucket_addon.example.api_response
}
```

## Schema

### Required

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
