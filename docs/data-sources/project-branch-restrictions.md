---
page_title: "bitbucket_project_branch_restrictions Data Source - bitbucket"
subcategory: "Internal Project Branch Restrictions"
description: |-
  Reads Bitbucket project-branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_project_branch_restrictions (Data Source)

Reads Bitbucket project-branch-restrictions via the Bitbucket Cloud API.

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
data "bitbucket_project_branch_restrictions" "example" {
  workspace = "my-workspace"
  project_key = "PROJ"
}

output "project_branch_restrictions_response" {
  value = data.bitbucket_project_branch_restrictions.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `project_key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
