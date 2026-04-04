---
page_title: "bitbucket_issues Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket issues via the Bitbucket Cloud API.
---

# bitbucket_issues (Data Source)

Reads Bitbucket issues via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/issues` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:issue:bitbucket` |
| List | `read:issue:bitbucket` |

## Example Usage

```hcl
data "bitbucket_issues" "example" {
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issues_response" {
  value = data.bitbucket_issues.example.api_response
}
```

## Schema

### Required
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
