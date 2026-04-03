---
page_title: "bitbucket_user_emails Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket user-emails via the Bitbucket Cloud API.
---

# bitbucket_user_emails (Data Source)

Reads Bitbucket user-emails via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_user_emails" "example" {
  email = "user@example.com"
}

output "user_emails_response" {
  value = data.bitbucket_user_emails.example.api_response
}
```

## Schema

### Required
- `email` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
