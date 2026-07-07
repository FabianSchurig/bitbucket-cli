# Reality-checks

Combined reality-checks that create real Bitbucket resources and validate provider behaviour end-to-end.

Locations:

- `terraform/reality-checks/01_create_repo` — creates a repo and then adds a branch restriction, a webhook, and a pipeline variable.
- `terraform/reality-checks/02_import_and_manage` — issue #111 coverage: reads a workspace (data source), creates a project via `request_body = jsonencode(...)`, and creates repositories using both typed fields (`description`, `project.key`) and the raw `request_body` escape hatch.
- `terraform/reality-checks/03_branching_model` — verifies the repository and project branching-model resources expose typed body fields (`branch_types`) and are idempotent + update-detecting.

## Reality-check 02 — issue #111 (import & manage projects/repos)

This check reproduces the two failure modes reported in issue #111 on affected
provider versions:

1. **Update detection.** Apply once, then re-apply appending a character to
   every description via `-var description_suffix=.` — the plan must show a real
   update. Pre-fix, the provider wrongly reported "no changes".
2. **Raw `request_body`.** Adding `request_body = jsonencode(...)` to an
   existing repo (`-var typed_use_request_body=true`) must plan cleanly. Pre-fix
   it failed with "planned value cty.NullVal ... does not match config value".

Toggles (`variables.tf`):

- `description_suffix` — text appended to every description (update-detection test).
- `typed_use_request_body` — adds a raw body to the typed repo (symptom 2).
- `enable_request_body` — the standalone `request_body`-driven repo.
- `create_project` — create a project via `request_body` and attach the repo.

### Validation results (2026-07-07)

Validated against the registry release **v0.17.4** and the PR #117 fix branch
(`copilot/modify-existing-resources`, built locally via `dev_overrides`):

| Scenario | v0.17.4 (released) | PR #117 fix |
|----------|--------------------|-------------|
| Change typed `description` (repo & project) | ❌ "No changes" (update dropped) | ✅ update planned |
| Add `request_body` to existing repo | ❌ "invalid plan … cty.NullVal" | ✅ clean plan |
| Create repo (typed / `request_body`) | ✅ works | ✅ works |
| Create/manage project with typed fields | ❌ post-create 404 (project_key=UUID) | ✅ works |
| Read workspace (data source) | ✅ works | ✅ works |

Two follow-up fixes landed alongside #117 so the full "workspaces → projects →
repositories" workflow works with plain typed config:

- **Projects gained typed fields.** `key`/`name`/`description`/`is_private` are
  now first-class attributes (no `jsonencode` required). Fixed in the codegen
  pipeline by inlining `requestBody` `$ref`s in `scripts/partition_spec.py`
  (with a safety-net patch in `scripts/enrich_spec.py`), so it survives every
  daily schema-sync run instead of being overwritten.
- **Project identifier bug fixed.** The provider populated the `{project_key}`
  path parameter with the project **UUID** instead of its **key**, 404ing the
  post-create read (and refresh/destroy). Fixed in
  `internal/tfprovider/resource.go` (`responseParamValue` now resolves a
  `{resource}_key` param to the response `key` field).

```bash
# 1) create the full stack (typed project + repos) with a few lines of config
cd terraform/reality-checks/02_import_and_manage
terraform apply -auto-approve

# 2) re-apply with a modified description to prove update detection
#    (project + both repos plan a real in-place update)
terraform apply -auto-approve -var description_suffix=.

# 3) prove the request_body escape hatch on an existing repo
terraform plan -var typed_use_request_body=true

# 4) clean up
terraform destroy -auto-approve
```

## Reality-check 03 — branching-model typed fields

Verifies the repository and project branching-model resources expose typed
body fields instead of only the raw `request_body`. Bitbucket's live spec omits
the `requestBody` on both branching-model PUT operations, so `HasBody` was false
and the settings were unreachable. `scripts/enrich_spec.py` now injects the
`branching_model_settings` body (regression-proofed against the daily
schema-sync that previously overwrote the hand-edited schema YAML in `#110`).

```bash
cd terraform/reality-checks/03_branching_model
terraform apply -auto-approve                      # create repo/project + models
terraform plan -var feature_prefix=feature2/       # update detected (branch_types)
terraform plan                                     # No changes (idempotent)
terraform destroy -auto-approve
```

Notes:

- `branch_types` must enumerate all four kinds (feature/bugfix/release/hotfix) —
  Bitbucket always returns the complete set.
- `development`/`production` are left unmanaged: their nested `is_valid` field is
  read-only but not marked `readOnly` in the spec, so managing them yields a
  benign no-op re-plan. The resource-level `ModifyPlan` now preserves unset
  Optional+Computed attributes from prior state, so leaving them unset is stable.

Running the combined reality-check

1. Copy `terraform/.env.example` to `${workspaceFolder}/.env` and fill in values.
2. From repo root, source the env file:

```bash
source ${PWD}/.env
```

3. Apply the combined test (helper script or manual):

```bash
terraform/scripts/run-apply.sh terraform/reality-checks/01_create_repo
# when done
terraform/scripts/run-destroy.sh terraform/reality-checks/01_create_repo
```

Notes

- The helper scripts source `${workspaceFolder}/.env` and then run Terraform; set `TF_VAR_workspace` in that `.env` file so tests pick it up automatically.
- CI should gate real runs with `TF_ACC=1` and use a dedicated test workspace and token with appropriate scopes.
