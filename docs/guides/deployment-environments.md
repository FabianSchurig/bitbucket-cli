---
page_title: "Managing deployment environments across repositories"
subcategory: "Guides"
description: |-
  Manage, import and parameterize Bitbucket deployment environments with one bitbucket_deployments resource per environment.
---

# Managing deployment environments across repositories

`bitbucket_deployments` manages **one deployment environment** per resource
instance. You do not have to describe every environment of a repository in one
place: give each environment its own resource (or `for_each` instance) and keep
its variables next to it.

| Attribute | Notes |
|-----------|-------|
| `name` | Display name. Changed in place. |
| `environment_type.name` | `Test`, `Staging` or `Production`. Required on create. Bitbucket cannot change it, so a change replaces the environment. |
| `restrictions.admin_only` | Only admins may deploy. Changed in place. |
| `uuid` | Computed. Use it for `bitbucket_deployment_variables.environment_uuid`. |

Fields that have no dedicated attribute (for example `rank` or `hidden`) can be
read from `api_response` with `jsondecode(...)`.

## Create an environment

```hcl
resource "bitbucket_deployments" "dt_production" {
  workspace        = "my-workspace"
  repo_slug        = "my-repo"
  name             = "DT Production"
  environment_type = { name = "Production" }
}

resource "bitbucket_deployment_variables" "region" {
  workspace        = "my-workspace"
  repo_slug        = "my-repo"
  environment_uuid = bitbucket_deployments.dt_production.uuid
  key              = "AWS_REGION"
  value            = "eu-west-1"
}
```

## Adopt existing environments

Every repository starts with the environments `Test`, `Staging` and
`Production`; they are hidden until a pipeline deploys to them. Creating an
environment whose name is already taken fails with `409 Conflict`, so import
existing environments instead.

The import ID is `workspace/repo_slug/environment_uuid`. The
`bitbucket_deployments` data source lists every environment of a repository, so
the UUIDs never have to be copied by hand. The example below keys environments
by slug, imports the ones that already exist and creates the rest
(`for_each` in `import` blocks requires Terraform 1.7+):

```hcl
locals {
  workspace = "my-workspace"
  repo_slug = "my-repo"

  # Keyed by environment slug, so existing environments are matched by slug.
  environments = {
    "staging"       = { name = "Staging", type = "Staging" }
    "production"    = { name = "Production", type = "Production" }
    "dt-staging"    = { name = "DT Staging", type = "Staging" }
    "dt-production" = { name = "DT Production", type = "Production" }
  }
}

data "bitbucket_deployments" "existing" {
  workspace = local.workspace
  repo_slug = local.repo_slug
}

locals {
  existing_uuids = {
    for env in jsondecode(data.bitbucket_deployments.existing.api_response) : env.slug => env.uuid
  }
}

import {
  for_each = {
    for slug, env in local.environments : slug => local.existing_uuids[slug]
    if contains(keys(local.existing_uuids), slug)
  }
  to = bitbucket_deployments.env[each.key]
  id = "${local.workspace}/${local.repo_slug}/${each.value}"
}

resource "bitbucket_deployments" "env" {
  for_each = local.environments

  workspace        = local.workspace
  repo_slug        = local.repo_slug
  name             = each.value.name
  environment_type = { name = each.value.type }
}
```

The import blocks can stay in place: once an environment is in state, Terraform
skips its import.

To see what a repository currently has before writing the map, output the
environments in a compact shape:

```hcl
output "environments" {
  value = {
    for env in jsondecode(data.bitbucket_deployments.existing.api_response) :
    env.slug => { name = env.name, type = env.environment_type.name, uuid = env.uuid }
  }
}
```

## Many repositories

Put the configuration above into a module with `workspace`, `repo_slug` and
`environments` variables and call it once per repository:

```hcl
module "environments" {
  source   = "./modules/repo-environments"
  for_each = local.repositories # e.g. { "dt/inf-devops" = { workspace = "dt", repository = "inf-devops", environments = {...} } }

  workspace    = each.value.workspace
  repo_slug    = each.value.repository
  environments = each.value.environments
}
```

Repositories that do not list an environment simply do not manage it; an
environment that is in state but removed from the map is deleted on the next
apply.
