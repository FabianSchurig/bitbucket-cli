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
  encoded_id = "snippet-id"
  workspace = "my-workspace"
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
```

## Schema

### Required
- `encoded_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
