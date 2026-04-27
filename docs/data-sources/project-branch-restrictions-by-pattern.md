---
page_title: "bitbucket_project_branch_restrictions_by_pattern Data Source - bitbucket"
subcategory: "Experimental"
description: |-
  Reads Bitbucket project-branch-restrictions-by-pattern via the Bitbucket Cloud API.
---

# bitbucket_project_branch_restrictions_by_pattern (Data Source)

Reads Bitbucket project-branch-restrictions-by-pattern via the Bitbucket Cloud API.

~> **Experimental — internal API.** This data source targets Bitbucket's
undocumented internal API at `https://bitbucket.org/!api/internal/`, which
**does not accept HTTP Basic Auth**. You must configure the provider with
`csrf_token` and `cloud_session_token` (or the `BITBUCKET_CSRF_TOKEN` /
`BITBUCKET_CLOUD_SESSION_TOKEN` environment variables) — see the
[provider documentation](../index.md#authenticating-against-the-internal-api)
for details.

Internal-API data sources are not auto-synced by the daily OpenAPI pipeline
and the underlying endpoint may change without notice. The
`cloud.session.token` cookie typically expires after ~1 month, so prefer
running these manually with freshly-copied browser cookies.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |
| List | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |
| List | — |

## Example Usage

```hcl
data "bitbucket_project_branch_restrictions_by_pattern" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "project_branch_restrictions_by_pattern_response" {
  value = data.bitbucket_project_branch_restrictions_by_pattern.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
