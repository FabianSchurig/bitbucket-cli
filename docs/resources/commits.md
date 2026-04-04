---
page_title: "bitbucket_commits Resource - bitbucket"
subcategory: "Commits"
description: |-
  Manages Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Resource)

Manages Bitbucket commits via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commits` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commits-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_commits" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `author_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.
- `committer_raw` (String) The raw committer value from the repository. This may be the only value available if the committer does not match a user in Bitbucket.
- `date` (String) date
- `hash` (String) hash
- `message` (String) message
- `parents` (String) parents (JSON array)
- `participants` (List of Object) participants
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `repository_created_on` (String) repository.created_on
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
- `repository_updated_on` (String) repository.updated_on
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
