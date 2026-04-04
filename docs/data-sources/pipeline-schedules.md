---
page_title: "bitbucket_pipeline_schedules Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-schedules via the Bitbucket Cloud API.
---

# bitbucket_pipeline_schedules (Data Source)

Reads Bitbucket pipeline-schedules via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-schedule-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/schedules` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-schedules-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_schedules_response" {
  value = data.bitbucket_pipeline_schedules.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `schedule_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the schedule was created.
- `updated_on` (String) The timestamp when the schedule was updated.
- `uuid` (String) The UUID identifying the schedule.
- `cron_pattern` (String) The cron expression with second precision (7 fields) that the schedule applies. For example, for expression: 0 0 12 * * ? *, will execute at 12pm UTC every day.
- `enabled` (String) Whether the schedule is enabled.
- `target_ref_name` (String) The name of the reference.
- `target_ref_type` (String) The type of reference (branch only). [branch]
- `target_selector_pattern` (String) The name of the matching pipeline definition.
- `target_selector_type` (String) The type of selector. [branches, tags, bookmarks, default, custom]
