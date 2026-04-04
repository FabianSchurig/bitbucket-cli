---
page_title: "bitbucket_repo_runners Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-runners via the Bitbucket Cloud API.
---

# bitbucket_repo_runners (Data Source)

Reads Bitbucket repo-runners via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:runner:bitbucket` |
| List | `read:runner:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "repo_runners_response" {
  value = data.bitbucket_repo_runners.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `runner_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
- `labels` (List of String) Labels assigned to the runner for identification and routing.
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
