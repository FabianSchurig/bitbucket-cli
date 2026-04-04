---
page_title: "bitbucket_reports Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket reports via the Bitbucket Cloud API.
---

# bitbucket_reports (Resource)

Manages Bitbucket reports via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-put) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Delete | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `report_id` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
