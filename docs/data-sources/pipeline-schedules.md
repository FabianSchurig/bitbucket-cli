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
  schedule_uuid = "{schedule-uuid}"
}

output "pipeline_schedules_response" {
  value = data.bitbucket_pipeline_schedules.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `schedule_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
