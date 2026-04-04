---
page_title: "bitbucket_project_default_reviewers Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_project_default_reviewers (Data Source)

Reads Bitbucket project-default-reviewers via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-selected-user-get) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/default-reviewers` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-default-reviewers-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_default_reviewers" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_default_reviewers_response" {
  value = data.bitbucket_project_default_reviewers.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `selected_user` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `account_id` (String) The user's Atlassian account ID.
- `account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `has_2fa_enabled` (String) has_2fa_enabled
- `is_staff` (String) is_staff
- `nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
- `uuid` (String) uuid
