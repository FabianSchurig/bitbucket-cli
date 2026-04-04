---
page_title: "bitbucket_pr Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pr via the Bitbucket Cloud API.
---

# bitbucket_pr (Resource)

Manages Bitbucket pr via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pullrequests` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-put) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| Read | `read:pullrequest:bitbucket` |
| Update | `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pr" "example" {
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `pull_request_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
