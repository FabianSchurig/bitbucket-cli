---
page_title: "bitbucket_project_default_reviewers Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_project_default_reviewers (Resource)

Manages Bitbucket project-default-reviewers via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-selected-user-put) |
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-selected-user-get) |
| Delete | `DELETE` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-selected-user-delete) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:project:bitbucket` |
| Read | `read:pullrequest:bitbucket` |
| Delete | `admin:project:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_project_default_reviewers" "example" {
  project_key = "PROJ"
  selected_user = "jdoe"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `selected_user` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
