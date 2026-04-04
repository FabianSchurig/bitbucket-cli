---
page_title: "bitbucket_tags Resource - bitbucket"
subcategory: "Refs"
description: |-
  Manages Bitbucket tags via the Bitbucket Cloud API.
---

# bitbucket_tags (Resource)

Manages Bitbucket tags via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/refs/tags` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-name-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/refs/tags/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-name-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:repository:bitbucket`, `read:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Delete | `write:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_tags" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `name` (String) Path parameter (auto-populated from API response).
- `date` (String) The date that the tag was created, if available (also computed from API response)
- `message` (String) The message associated with the tag, if available. (also computed from API response)
- `tagger` (Object) tagger (also computed from API response)
  Nested schema:
  - `raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.

- `target` (Object) target (also computed from API response)
  Nested schema:
  - `message` (String) message
  - `parents` (String) parents (JSON array)
  - `hash` (String) hash
  - `date` (String) date

- `type` (String) type (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
