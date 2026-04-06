#!/usr/bin/env bash
set -euo pipefail

# Usage: run-destroy.sh [TEST_DIR]

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
terraform destroy -auto-approve -input=false
