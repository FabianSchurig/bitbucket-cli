---
page_title: "bitbucket_default_reviewers Resource - bitbucket"
subcategory: "Pull Requests"
description: |-
  Manages Bitbucket default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_default_reviewers (Resource)

Manages Bitbucket default-reviewers via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-target-username-put) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-target-username-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-target-username-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/default-reviewers` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket` |
| Read | `read:pullrequest:bitbucket` |
| Delete | `admin:repository:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_default_reviewers" "example" {
  repo_slug = "my-repo"
  target_username = "jdoe"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `target_username` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
