/*
  Reality-check 03 — branching-model typed fields.

  Verifies that the repository and project branching-model resources expose
  typed body fields (development / production / branch_types) instead of only
  the raw request_body escape hatch. Bitbucket's live spec omits the
  requestBody on the two branching-model PUT operations, so HasBody was false
  and these fields were unreachable; scripts/enrich_spec.py now injects the
  branching_model_settings body (regression-proofed against the daily
  schema-sync that previously overwrote the hand-edited schema YAML).

  Change `-var feature_prefix=feature2/` on a second apply to confirm the
  branching-model update is planned (not silently dropped).
*/

terraform {
  required_version = ">= 1.0.0"
  required_providers {
    bitbucket = { source = "FabianSchurig/bitbucket" }
    random    = { source = "hashicorp/random" }
  }
}

provider "bitbucket" {}
provider "random" {}

resource "random_pet" "repo" {
  length = 2
}

resource "random_string" "project_key" {
  length  = 6
  upper   = true
  lower   = false
  numeric = true
  special = false
}

# Bitbucket's branching model always returns all four branch-type kinds, so the
# configuration must enumerate all of them or Terraform reports drift
# ("new element has appeared"). Only feature's prefix is parameterised so the
# update-detection test has something to change.
locals {
  branch_types = [
    { kind = "feature", prefix = var.feature_prefix, enabled = true },
    { kind = "bugfix", prefix = "bugfix/", enabled = true },
    { kind = "release", prefix = "release/", enabled = true },
    { kind = "hotfix", prefix = "hotfix/", enabled = true },
  ]
}

# A repository to host the branching model.
resource "bitbucket_repos" "repo" {
  workspace  = var.workspace
  repo_slug  = "tf-reality-bm-${random_pet.repo.id}"
  is_private = true
}

# Repository branching model — configured with TYPED fields (no request_body).
#
# Note: `development`/`production` are intentionally left unmanaged here. They
# are objects whose nested `is_valid` field is read-only ("ignored when
# updating") but Bitbucket's spec does not mark it `readOnly`, so managing them
# yields a benign perpetual no-op re-plan (is_valid -> known after apply). That
# is a separate framework limitation, unrelated to the requestBody/typed-fields
# fix this check verifies. `branch_types` fields are all writable and stable.
resource "bitbucket_branching_model" "repo_model" {
  workspace = var.workspace
  repo_slug = bitbucket_repos.repo.repo_slug

  branch_types = local.branch_types
}

# Optional: a project + its branching model, exercising the project PUT body.
resource "bitbucket_projects" "proj" {
  count       = var.manage_project_model ? 1 : 0
  workspace   = var.workspace
  key         = "BM${random_string.project_key.result}"
  name        = "TF Branching Model BM${random_string.project_key.result}"
  description = "reality-check 03"
  is_private  = true
}

resource "bitbucket_project_branching_model" "project_model" {
  count       = var.manage_project_model ? 1 : 0
  workspace   = var.workspace
  project_key = bitbucket_projects.proj[0].key

  branch_types = local.branch_types
}

output "repo_slug" {
  value = bitbucket_repos.repo.repo_slug
}

output "repo_feature_prefix" {
  value = one([for bt in bitbucket_branching_model.repo_model.branch_types : bt.prefix if bt.kind == "feature"])
}

output "project_key" {
  value = var.manage_project_model ? bitbucket_projects.proj[0].key : null
}
