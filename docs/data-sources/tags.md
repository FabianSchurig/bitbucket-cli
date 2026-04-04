---
page_title: "bitbucket_tags Data Source - bitbucket"
subcategory: "Refs"
description: |-
  Reads Bitbucket tags via the Bitbucket Cloud API.
---

# bitbucket_tags (Data Source)

Reads Bitbucket tags via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-name-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_tags" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "tags_response" {
  value = data.bitbucket_tags.example.api_response
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
- `target_participants` (List of Object) participants
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `target_summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `target_summary_raw` (String) The text as it was typed by a user.
- `date` (String) The date that the tag was created, if available
- `message` (String) The message associated with the tag, if available.
- `tagger_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.
- `target_date` (String) target.date
- `target_hash` (String) target.hash
- `target_message` (String) target.message
- `target_parents` (String) parents (JSON array)
- `type` (String) type
