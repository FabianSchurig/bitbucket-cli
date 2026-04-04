---
page_title: "bitbucket_pipelines Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipelines via the Bitbucket Cloud API.
---

# bitbucket_pipelines (Data Source)

Reads Bitbucket pipelines via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-pipeline-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pipeline_uuid = "pipeline-uuid"
}

output "pipelines_response" {
  value = data.bitbucket_pipelines.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `pipeline_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
