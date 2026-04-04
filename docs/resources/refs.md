---
page_title: "bitbucket_refs Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket refs via the Bitbucket Cloud API.
---

# bitbucket_refs (Resource)

Manages Bitbucket refs via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/refs/branches` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-branches-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/branches/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-branches-name-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/refs/branches/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-branches-name-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:repository:bitbucket`, `read:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Delete | `write:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_refs" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `name` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `name` (String) The name of the ref.
- `type` (String) type
