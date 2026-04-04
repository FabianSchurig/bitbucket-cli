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

## Example Usage

```hcl
resource "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  runner_uuid = "{runner-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
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
