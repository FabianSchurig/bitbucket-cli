---
page_title: "bitbucket_deployments Resource - bitbucket"
subcategory: "Deployments"
description: |-
  Manages Bitbucket deployments via the Bitbucket Cloud API.
---

# bitbucket_deployments (Resource)

Manages Bitbucket deployments via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/environments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-environments-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/environments/{environment_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-environments-environment-uuid-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/environments/{environment_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-environments-environment-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/environments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-environments-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| Delete | `admin:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_deployments" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `environment_uuid` (String) Path parameter (auto-populated from API response).
- `name` (String) The name of the environment. (also computed from API response)
- `uuid` (String) The UUID identifying the environment. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
