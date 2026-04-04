---
page_title: "bitbucket_issue_comments Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket issue-comments via the Bitbucket Cloud API.
---

# bitbucket_issue_comments (Data Source)

Reads Bitbucket issue-comments via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_issue_comments" "example" {
  comment_id = "1"
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issue_comments_response" {
  value = data.bitbucket_issue_comments.example.api_response
}
```

## Schema

### Required
- `comment_id` (String) Path parameter.
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
