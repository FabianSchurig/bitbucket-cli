---
page_title: "bitbucket_issues Data Source - bitbucket"
subcategory: "Issues"
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
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issues_response" {
  value = data.bitbucket_issues.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `issue_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
- `assignee` (Object) assignee
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `component` (Object) component
  Nested schema:
  - `name` (String) name
  - `id` (String) id

- `content` (Object) content
  Nested schema:
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
  - `raw` (String) The text as it was typed by a user.

- `edited_on` (String) edited_on
- `kind` (String) [bug, enhancement, proposal, task]
- `milestone` (Object) milestone
  Nested schema:
  - `name` (String) name
  - `id` (String) id

- `priority` (String) [trivial, minor, major, critical, blocker]
- `reporter` (Object) reporter
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `repository` (Object) repository
  Nested schema:
  - `description` (String) description
  - `fork_policy` (String) 
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `has_issues` (String) 
  - `has_wiki` (String) 
  - `is_private` (String) is_private
  - `language` (String) language
  - `name` (String) name
  - `scm` (String) [git]
  - `size` (String) size
  - `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.

- `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
- `title` (String) title
- `version` (Object) version
  Nested schema:
  - `name` (String) name
  - `id` (String) id

- `votes` (String) votes
