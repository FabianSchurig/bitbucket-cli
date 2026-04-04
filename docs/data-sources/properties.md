---
page_title: "bitbucket_properties Data Source - bitbucket"
subcategory: "Properties"
description: |-
  Reads Bitbucket properties via the Bitbucket Cloud API.
---

# bitbucket_properties (Data Source)

Reads Bitbucket properties via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/properties/{app_key}/{property_name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-properties/#api-repositories-workspace-repo-slug-properties-app-key-property-name-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |

## Example Usage

```hcl
data "bitbucket_properties" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  app_key = "my-app"
  property_name = "my-property"
}

output "properties_response" {
  value = data.bitbucket_properties.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `app_key` (String) Path parameter.
- `property_name` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `_attributes` (List of String) _attributes [public, read_only]
