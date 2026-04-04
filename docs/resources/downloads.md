---
page_title: "bitbucket_downloads Resource - bitbucket"
subcategory: "Downloads"
description: |-
  Manages Bitbucket downloads via the Bitbucket Cloud API.
---

# bitbucket_downloads (Resource)

Manages Bitbucket downloads via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/downloads` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/downloads/{filename}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-filename-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/downloads/{filename}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-filename-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/downloads` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Delete | `write:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_downloads" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `filename` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
