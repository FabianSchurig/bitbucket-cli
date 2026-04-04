---
page_title: "bitbucket_properties Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket properties via the Bitbucket Cloud API.
---

# bitbucket_properties (Resource)

Manages Bitbucket properties via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/properties/{app_key}/{property_name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-properties/#api-repositories-workspace-repo-slug-properties-app-key-property-name-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/properties/{app_key}/{property_name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-properties/#api-repositories-workspace-repo-slug-properties-app-key-property-name-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/properties/{app_key}/{property_name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-properties/#api-repositories-workspace-repo-slug-properties-app-key-property-name-delete) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |
| Update | — |
| Delete | — |

## Example Usage

```hcl
resource "bitbucket_properties" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  app_key = "my-app"
  property_name = "my-property"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `app_key` (String) Path parameter.
- `property_name` (String) Path parameter.

### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `_attributes` (String) _attributes (JSON array)
