---
page_title: "bitbucket_issues Resource - bitbucket"
subcategory: "Issues"
description: |-
  Manages Bitbucket issues via the Bitbucket Cloud API.
---

# bitbucket_issues (Resource)

Manages Bitbucket issues via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/issues` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/issues` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:issue:bitbucket`, `write:issue:bitbucket` |
| Read | `read:issue:bitbucket` |
| Update | `read:issue:bitbucket`, `write:issue:bitbucket` |
| Delete | `delete:issue:bitbucket` |
| List | `read:issue:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_issues" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `issue_id` (String) Path parameter (auto-populated from API response).
- `assignee` (Object) assignee (also computed from API response)
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `component` (Object) component (also computed from API response)
  Nested schema:
  - `name` (String) name

- `content` (Object) content (also computed from API response)
  Nested schema:
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
  - `raw` (String) The text as it was typed by a user.

- `edited_on` (String) edited_on (also computed from API response)
- `kind` (String) [bug, enhancement, proposal, task] (also computed from API response)
- `milestone` (Object) milestone (also computed from API response)
  Nested schema:
  - `name` (String) name

- `priority` (String) [trivial, minor, major, critical, blocker] (also computed from API response)
- `reporter` (Object) reporter (also computed from API response)
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `repository` (Object) repository (also computed from API response)
  Nested schema:
  - `is_private` (String) is_private
  - `description` (String) description
  - `fork_policy` (String) 
  - `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `name` (String) name
  - `size` (String) size
  - `has_issues` (String) 
  - `has_wiki` (String) 
  - `scm` (String) [git]
  - `language` (String) language

- `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed] (also computed from API response)
- `title` (String) title (also computed from API response)
- `version` (Object) version (also computed from API response)
  Nested schema:
  - `name` (String) name

- `votes` (String) votes (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
