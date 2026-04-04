---
page_title: "bitbucket_refs Data Source - bitbucket"
subcategory: "Refs"
description: |-
  Reads Bitbucket refs via the Bitbucket Cloud API.
---

# bitbucket_refs (Data Source)

Reads Bitbucket refs via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/branches/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-branches-name-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_refs" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `name` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `merge_strategies` (List of String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
- `target_date` (String) target.date
- `target_hash` (String) target.hash
- `target_message` (String) target.message
- `target_parents` (String) parents (JSON array)
- `target_participants` (List of Object) participants
  Nested schema:
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.
  - `role` (String) [PARTICIPANT, REVIEWER]

- `target_summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `target_summary_raw` (String) The text as it was typed by a user.
- `type` (String) type
