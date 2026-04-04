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
- `assignee_created_on` (String) assignee.created_on
- `component_id` (String) component.id
- `created_on` (String) created_on
- `milestone_id` (String) milestone.id
- `reporter_created_on` (String) reporter.created_on
- `repository_created_on` (String) repository.created_on
- `repository_updated_on` (String) repository.updated_on
- `updated_on` (String) updated_on
- `version_id` (String) version.id
- `assignee_display_name` (String) assignee.display_name
- `assignee_uuid` (String) assignee.uuid
- `component_name` (String) component.name
- `content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `content_raw` (String) The text as it was typed by a user.
- `edited_on` (String) edited_on
- `kind` (String) [bug, enhancement, proposal, task]
- `milestone_name` (String) milestone.name
- `priority` (String) [trivial, minor, major, critical, blocker]
- `reporter_display_name` (String) reporter.display_name
- `reporter_uuid` (String) reporter.uuid
- `repository_description` (String) repository.description
- `repository_fork_policy` (String) 
- `repository_full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
- `repository_has_issues` (String) 
- `repository_has_wiki` (String) 
- `repository_is_private` (String) repository.is_private
- `repository_language` (String) repository.language
- `repository_name` (String) repository.name
- `repository_scm` (String) [git]
- `repository_size` (String) repository.size
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
- `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
- `title` (String) title
- `version_name` (String) version.name
- `votes` (String) votes
