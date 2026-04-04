---
page_title: "bitbucket_pipeline_schedules Resource - bitbucket"
subcategory: "Pipelines"
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
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `schedule_uuid` (String) Path parameter (auto-populated from API response).
- `cron_pattern` (String) The cron expression with second precision (7 fields) that the schedule applies. For example, for expression: 0 0 12 * * ? *, will execute at 12pm UTC every day. (also computed from API response)
- `enabled` (String) Whether the schedule is enabled. (also computed from API response)
- `target` (Object) The target on which the schedule will be executed. (also computed from API response)
  Nested schema:
  - `selector` (Object) selector
    - `type` (String) The type of selector. [branches, tags, bookmarks, default, custom]
    - `pattern` (String) The name of the matching pipeline definition.
  - `ref_name` (String) The name of the reference.
  - `ref_type` (String) The type of reference (branch only). [branch]

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the schedule was created.
- `updated_on` (String) The timestamp when the schedule was updated.
- `uuid` (String) The UUID identifying the schedule.
