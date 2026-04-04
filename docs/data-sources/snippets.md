---
page_title: "bitbucket_snippets Data Source - bitbucket"
subcategory: "Snippets"
description: |-
  Reads Bitbucket snippets via the Bitbucket Cloud API.
---

# bitbucket_snippets (Data Source)

Reads Bitbucket snippets via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/snippets/{workspace}/{encoded_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-workspace-encoded-id-get) |
| List | `GET` | `/snippets` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:snippet:bitbucket` |
| List | `read:snippet:bitbucket` |

## Example Usage

```hcl
data "bitbucket_snippets" "example" {
}

output "snippets_response" {
  value = data.bitbucket_snippets.example.api_response
}
```

## Schema

### Required

### Optional
- `encoded_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.
- `workspace` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `is_private` (String) is_private
- `scm` (String) The DVCS used to store the snippet. [git]
- `title` (String) title
- `updated_on` (String) updated_on
