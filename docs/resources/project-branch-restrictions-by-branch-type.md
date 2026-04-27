---
page_title: "bitbucket_project_branch_restrictions_by_branch_type Resource - bitbucket"
subcategory: "Internal Project Branch Restrictions"
description: |-
  Manages Bitbucket project-branch-restrictions-by-branch-type via the Bitbucket Cloud API.
---

# bitbucket_project_branch_restrictions_by_branch_type (Resource)

Manages Bitbucket project-branch-restrictions-by-branch-type via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-branch-type/{branch_type}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-by-branch-type-branch-type-put) |
| Read | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |
| Update | `PUT` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-branch-type/{branch_type}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-by-branch-type-branch-type-put) |
| Delete | `PUT` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-branch-type/{branch_type}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-by-branch-type-branch-type-put) |
| List | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | — |
| Read | — |
| Update | — |
| Delete | — |
| List | — |

## Example Usage

```hcl
resource "bitbucket_project_branch_restrictions_by_branch_type" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
  branch_type = "example-value"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.
- `branch_type` (String) Path parameter.
- `values` (List of Object) The full desired list of rules.
  Nested schema:
  - `branch_match_kind` (String) How the rule selects branches. Must be "glob" when used with
  - `branch_type` (String) Branching-model branch type (when branch_match_kind=branching_model).
  - `groups` (List of Object) Groups exempted from the restriction. Only applies to "push" and
    - `owner` (String) owner
    - `slug` (String) The group slug.
  - `kind` (String) The kind of restriction.
  - `pattern` (String) Branch glob pattern (when branch_match_kind=glob).
  - `users` (List of Object) Users exempted from the restriction. Only applies to "push" and
    - `uuid` (String) The user UUID, surrounded by curly-braces.
  - `value` (String) Numeric threshold for kinds that require one (e.g.


### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
