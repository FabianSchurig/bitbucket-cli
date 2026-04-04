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
- `assignee_display_name` (String) assignee.display_name (also computed from API response)
- `assignee_uuid` (String) assignee.uuid (also computed from API response)
- `component_name` (String) component.name (also computed from API response)
- `content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext] (also computed from API response)
- `content_raw` (String) The text as it was typed by a user. (also computed from API response)
- `edited_on` (String) edited_on (also computed from API response)
- `kind` (String) [bug, enhancement, proposal, task] (also computed from API response)
- `milestone_name` (String) milestone.name (also computed from API response)
- `priority` (String) [trivial, minor, major, critical, blocker] (also computed from API response)
- `reporter_display_name` (String) reporter.display_name (also computed from API response)
- `reporter_uuid` (String) reporter.uuid (also computed from API response)
- `repository_description` (String) repository.description (also computed from API response)
- `repository_fork_policy` (String)  (also computed from API response)
- `repository_full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs. (also computed from API response)
- `repository_has_issues` (String)  (also computed from API response)
- `repository_has_wiki` (String)  (also computed from API response)
- `repository_is_private` (String) repository.is_private (also computed from API response)
- `repository_language` (String) repository.language (also computed from API response)
- `repository_name` (String) repository.name (also computed from API response)
- `repository_scm` (String) [git] (also computed from API response)
- `repository_size` (String) repository.size (also computed from API response)
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user. (also computed from API response)
- `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed] (also computed from API response)
- `title` (String) title (also computed from API response)
- `version_name` (String) version.name (also computed from API response)
- `votes` (String) votes (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `assignee_created_on` (String) assignee.created_on
- `component_id` (String) component.id
- `created_on` (String) created_on
- `milestone_id` (String) milestone.id
- `reporter_created_on` (String) reporter.created_on
- `repository_created_on` (String) repository.created_on
- `repository_updated_on` (String) repository.updated_on
- `updated_on` (String) updated_on
- `version_id` (String) version.id
