---
page_title: "bitbucket_project_branch_restrictions Resource - bitbucket"
subcategory: "Internal Project Branch Restrictions"
description: |-
  Manages Bitbucket project-branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_project_branch_restrictions (Resource)

Manages Bitbucket project-branch-restrictions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |
| Update | `PUT` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-pattern/{pattern}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-by-pattern-pattern-put) |
| List | `GET` | `https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-project-branch-restrictions-(internal)/#api-https:--bitbucket.org-!api-internal-workspaces-workspace-projects-project-key-branch-restrictions-group-by-branch--get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |
| Update | — |
| List | — |

## Example Usage

```hcl
resource "bitbucket_project_branch_restrictions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.
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
- `pattern` (String) Path parameter (auto-populated from API response).
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
