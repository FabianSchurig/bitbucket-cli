---
page_title: "bitbucket_pipeline_schedules Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-schedules via the Bitbucket Cloud API.
---

# bitbucket_pipeline_schedules (Resource)

Manages Bitbucket pipeline-schedules via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-schedule-uuid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-schedule-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-schedule-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:pipeline:bitbucket`, `write:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| Update | `read:pipeline:bitbucket`, `write:pipeline:bitbucket` |
| Delete | `write:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  schedule_uuid = "{schedule-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `schedule_uuid` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
