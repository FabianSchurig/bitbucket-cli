# Reality-checks

Single combined reality-check that creates a repository and performs multiple checks (branch restriction, webhook, pipeline variable).

Location:

- `terraform/reality-checks/01_create_repo` — creates a repo and then adds a branch restriction, a webhook, and a pipeline variable.

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
