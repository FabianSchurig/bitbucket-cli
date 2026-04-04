---
page_title: "bitbucket_repo_runners Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket repo-runners via the Bitbucket Cloud API.
---

# bitbucket_repo_runners (Resource)

Manages Bitbucket repo-runners via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:runner:bitbucket`, `read:runner:bitbucket` |
| Read | `read:runner:bitbucket` |
| Update | `read:runner:bitbucket`, `write:runner:bitbucket` |
| Delete | `write:runner:bitbucket` |
| List | `read:runner:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `runner_uuid` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
- `labels` (String) Labels assigned to the runner for identification and routing. (JSON array)
- `name` (String) The name of the runner.
- `oauth_client_audience` (String) The intended audience for the OAuth token.
- `oauth_client_id` (String) The OAuth client ID.
- `oauth_client_secret` (String) The OAuth client secret. This is an optional element that is only provided once.
- `oauth_client_token_endpoint` (String) The OAuth token endpoint URL.
- `state_cordoned` (String) Whether the runner is cordoned (prevented from accepting new steps).
- `state_status` (String) The current status of the runner. [UNREGISTERED, ONLINE, OFFLINE, DISABLED, ENABLED, UNHEALTHY]
- `state_updated_on` (String) The timestamp when the runner state was last updated.
- `state_version_current` (String) The current recommended version of the runner.
- `state_version_version` (String) The currently installed version of the runner.
- `updated_on` (String) The timestamp when the runner was last updated.
- `uuid` (String) The UUID identifying the runner.
