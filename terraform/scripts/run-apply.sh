#!/usr/bin/env bash
set -euo pipefail

# Usage: run-apply.sh [TEST_DIR]
# Example: terraform/scripts/run-apply.sh terraform/reality-checks/01_create_repo

WORKSPACE_FOLDER=${WORKSPACE_FOLDER:-"${PWD}"}
ENV_FILE="${WORKSPACE_FOLDER}/.env"

if [ -f "${ENV_FILE}" ]; then
  echo "Sourcing ${ENV_FILE}"
  set -o allexport
  . "${ENV_FILE}"
  set +o allexport
else
  echo "Warning: ${ENV_FILE} not found. Continuing without loading env."
fi

DIR=${1:-$(pwd)}
cd "${DIR}"

terraform init -input=false
terraform apply -auto-approve -input=false
