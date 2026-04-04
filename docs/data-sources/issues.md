---
page_title: "bitbucket_issues Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket issues via the Bitbucket Cloud API.
---

# bitbucket_issues (Data Source)

Reads Bitbucket issues via the Bitbucket Cloud API.

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
- `component_id` (String) component.id
- `created_on` (String) created_on
- `milestone_id` (String) milestone.id
- `updated_on` (String) updated_on
- `version_id` (String) version.id
- `component_name` (String) component.name
- `content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `content_raw` (String) The text as it was typed by a user.
- `edited_on` (String) edited_on
- `kind` (String) [bug, enhancement, proposal, task]
- `milestone_name` (String) milestone.name
- `priority` (String) [trivial, minor, major, critical, blocker]
- `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
- `title` (String) title
- `version_name` (String) version.name
- `votes` (String) votes
