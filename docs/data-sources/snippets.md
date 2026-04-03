---
page_title: "bitbucket_snippets Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket snippets via the Bitbucket Cloud API.
---

# bitbucket_snippets (Data Source)

Reads Bitbucket snippets via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_snippets" "example" {
  workspace = "my-workspace"
  encoded_id = "snippet-id"
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `encoded_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
