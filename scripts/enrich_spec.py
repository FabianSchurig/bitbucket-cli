#!/usr/bin/env python3
"""
enrich_spec.py: Inject operationIds into the Bitbucket OpenAPI spec.

Strategy: slugify(summary) if present, else "{method}_{path_slug}".

operationIds are the primary key of the whole pipeline: they become CLI command
names, MCP tool names and the keys hand-written code (notably
``internal/tfprovider/crud_config.go``) refers to. Deriving them purely from the
live spec makes them unstable — the collision fallback is "first occurrence
wins", so removing one operation silently *renames* an unrelated one (when
Atlassian dropped ``GET /user/permissions/workspaces``, ``GET /user/workspaces``
inherited its ``listWorkspacesForTheCurrentUser`` id).

To make ids stable for good, assignments are persisted in a committed lockfile
(``schema/operation-ids.json``) keyed by ``"<method> <path>"``. Locked
operations always keep their id, ids of operations that disappear stay reserved
so they are never handed to a different endpoint, and only genuinely new
operations get a freshly derived id.

Usage: python3 enrich_spec.py <input.json> <output.json> [--lock <lockfile.json>]
"""

import copy
import json
import re
import sys
from pathlib import Path

# Default location of the committed operationId lockfile, relative to the
# repository root (the parent directory of scripts/).
DEFAULT_LOCK_PATH = Path(__file__).resolve().parent.parent / \
    "schema" / "operation-ids.json"


def safe_path(raw: str, allowed_extensions: set[str]) -> Path:
    """Resolve a CLI argument to a canonical path, guarding against path injection."""
    if "\0" in raw:
        raise ValueError(f"Null byte in path: {raw!r}")
    p = Path(raw).resolve()
    if p.suffix.lower() not in allowed_extensions:
        raise ValueError(
            f"Unexpected file extension {p.suffix!r}, expected one of {allowed_extensions}"
        )
    return p


def to_camel(s: str) -> str:
    """'List pull requests' -> 'listPullRequests'"""
    words = re.sub(r"[^a-zA-Z0-9 ]", "", s).split()
    if not words:
        return ""
    return words[0].lower() + "".join(w.title() for w in words[1:])


def path_slug(path: str, method: str) -> str:
    """/repositories/{workspace}/{repo_slug}/pullrequests' + 'get' -> 'getRepositoriesPullrequests'"""
    parts = [p for p in path.split("/") if p and not p.startswith("{")]
    return method.lower() + "".join(p.title() for p in parts)


# HTTP methods considered operations, in a fixed order for deterministic id
# assignment.
HTTP_METHODS = ("get", "post", "put", "patch", "delete")


# ─── Missing requestBody patches ──────────────────────────────────────────────
# Bitbucket's published OpenAPI spec omits the requestBody on a handful of write
# operations even though the endpoints accept a body. Without a requestBody the
# generators emit HasBody=false and the CLI/MCP/Terraform layers expose no typed
# body fields — only the raw `request_body`/`--body` escape hatch. Injecting the
# body here (before partitioning) is the maintainable fix: it re-applies on every
# schema-sync run, whereas hand-editing schema/*-schema.yaml is silently
# overwritten by the daily partition_spec.py --all regeneration.
#
# Keyed by (method, path). Values are OpenAPI requestBody objects; a $ref to an
# existing components/requestBodies entry is preferred so the injected body stays
# in sync with the schema Atlassian does publish for sibling operations.
REQUEST_BODY_PATCHES: dict[tuple[str, str], dict] = {
    # createAProjectInAWorkspace — POST /workspaces/{workspace}/projects has no
    # requestBody in the live spec, yet it accepts the same project body as the
    # sibling PUT updateAProjectForAWorkspace. Reuse the published
    # requestBodies/project component so key/name/description/is_private become
    # typed fields instead of requiring jsonencode(...).
    ("post", "/workspaces/{workspace}/projects"): {
        "$ref": "#/components/requestBodies/project"
    },
    # update the branching model config (repository + project). The live spec
    # omits the requestBody for both PUTs, so HasBody=false and the settings
    # (development/production/branch_types) are unreachable except via the raw
    # request_body. Reference the published branching_model_settings schema.
    ("put", "/repositories/{workspace}/{repo_slug}/branching-model/settings"): {
        "content": {
            "application/json": {
                "schema": {"$ref": "#/components/schemas/branching_model_settings"}
            }
        },
        "description": "The updated branching model configuration",
        "required": False,
    },
    ("put", "/workspaces/{workspace}/projects/{project_key}/branching-model/settings"): {
        "content": {
            "application/json": {
                "schema": {"$ref": "#/components/schemas/branching_model_settings"}
            }
        },
        "description": "The updated branching model configuration",
        "required": False,
    },
    # repository deploy keys — the add (POST) operation omits its requestBody in
    # the live spec, so HasBody=false and create sends an empty body (400). Only
    # the POST is patched: Bitbucket rejects `key` on the update PUT ("you can't
    # modify the contents of an access key"), and adding a body there would send
    # the immutable key on every update. Deploy keys are effectively immutable.
    ("post", "/repositories/{workspace}/{repo_slug}/deploy-keys"): {
        "content": {
            "application/json": {"schema": {"$ref": "#/components/schemas/deploy_key"}}
        },
        "description": "The deploy key to add",
        "required": False,
    },
}


def _refs_resolvable(node, spec: dict) -> bool:
    """Return True when every ``$ref`` in ``node`` resolves within ``spec``.

    Guards apply_request_body_patches against introducing a dangling reference
    (e.g. if Atlassian renames a component), which would break partitioning.
    """
    if isinstance(node, dict):
        ref = node.get("$ref")
        if isinstance(ref, str) and ref.startswith("#/"):
            target = spec
            for part in ref.lstrip("#/").split("/"):
                if isinstance(target, dict) and part in target:
                    target = target[part]
                else:
                    return False
        return all(_refs_resolvable(v, spec) for v in node.values())
    if isinstance(node, list):
        return all(_refs_resolvable(item, spec) for item in node)
    return True


def apply_request_body_patches(spec: dict) -> int:
    """Inject requestBody objects for operations the live spec leaves bodyless.

    Applies each entry in REQUEST_BODY_PATCHES only when the target operation
    exists, does not already declare a requestBody, and every ``$ref`` inside
    the patch resolves against the current spec — so partitioning never
    encounters a dangling reference.
    """
    applied = 0
    for (method, path), body in REQUEST_BODY_PATCHES.items():
        op = spec.get("paths", {}).get(path, {}).get(method)
        if not op or op.get("requestBody"):
            continue
        if not _refs_resolvable(body, spec):
            continue
        op["requestBody"] = copy.deepcopy(body)
        applied += 1
    return applied


def lock_key(path: str, method: str) -> str:
    """Key used in the operationId lockfile: ``"get /workspaces"``."""
    return f"{method.lower()} {path}"


def load_lock(lock_path: Path) -> dict[str, str]:
    """Load the operationId lockfile, tolerating a missing or empty file.

    A malformed lockfile is a hard error: silently falling back to derived ids
    would rename CLI commands and MCP tools without anyone noticing.
    """
    if not lock_path.exists():
        return {}
    raw = lock_path.read_text().strip()
    if not raw:
        return {}
    def reject_duplicate_keys(pairs):
        data = {}
        for key, value in pairs:
            if key in data:
                raise ValueError(
                    f"Lockfile {lock_path} contains duplicate key {key!r}"
                )
            data[key] = value
        return data

    data = json.loads(raw, object_pairs_hook=reject_duplicate_keys)
    if not isinstance(data, dict):
        raise ValueError(f"Lockfile {lock_path} must contain a JSON object")
    result: dict[str, str] = {}
    ids: set[str] = set()
    for key, value in data.items():
        if not isinstance(key, str) or not key.strip():
            raise ValueError(f"Lockfile {lock_path} contains an invalid key")
        if not isinstance(value, str) or not value.strip():
            raise ValueError(
                f"Lockfile {lock_path} has a non-empty string id for {key!r}"
            )
        if value in ids:
            raise ValueError(
                f"Lockfile {lock_path} reuses operation id {value!r}"
            )
        ids.add(value)
        result[key] = value
    return result


def _derive_id(op: dict, path: str, method: str) -> str:
    """Derive a candidate operationId for an operation without a locked id."""
    operation_id = op.get("operationId")
    if isinstance(operation_id, str) and operation_id:
        return operation_id
    summary = op.get("summary", "")
    return to_camel(summary) if isinstance(summary, str) and summary else path_slug(
        path, method
    )


def _unique_id(candidate: str, path: str, method: str, taken: set[str]) -> str:
    """Return ``candidate`` or a collision-free variant of it."""
    if candidate and candidate not in taken:
        return candidate
    slug = path_slug(path, method)
    if slug not in taken:
        return slug
    suffix = 2
    while f"{slug}{suffix}" in taken:
        suffix += 1
    return f"{slug}{suffix}"


def assign_operation_ids(spec: dict, lock: dict[str, str]) -> tuple[int, dict[str, str]]:
    """Assign stable operationIds to every operation in ``spec``.

    Locked operations keep their recorded id. Ids belonging to operations that
    vanished from the published spec stay reserved, so a brand-new or renamed
    endpoint can never inherit an id that already means something else — the
    failure mode that renamed ``getUserWorkspaces`` when Atlassian dropped a
    sibling path. New operations derive an id from their summary and fall back
    to a path-based slug on collision.

    Returns the number of operations processed and the updated lock mapping.
    """
    updated = dict(lock)
    # Every id ever handed out stays reserved, including ids of operations that
    # are no longer published.
    taken = set(lock.values())
    pending: list[tuple[str, str, dict]] = []

    count = 0
    # Iterate in a spec-order-independent order so ids do not depend on how
    # Atlassian happens to serialize the document.
    for path in sorted(spec.get("paths", {})):
        path_item = spec["paths"][path]
        if not isinstance(path_item, dict):
            continue
        for method in HTTP_METHODS:
            op = path_item.get(method)
            if not isinstance(op, dict):
                continue
            count += 1
            locked = lock.get(lock_key(path, method))
            if locked:
                op["operationId"] = locked
            else:
                pending.append((path, method, op))

    for path, method, op in pending:
        oid = _unique_id(_derive_id(op, path, method), path, method, taken)
        op["operationId"] = oid
        taken.add(oid)
        updated[lock_key(path, method)] = oid

    return count, dict(sorted(updated.items()))


def write_lock(lock: dict[str, str], lock_path: Path) -> None:
    """Persist the lockfile with stable ordering and a trailing newline."""
    lock_path.write_text(json.dumps(dict(sorted(lock.items())), indent=2) + "\n")


def main():
    args = sys.argv[1:]
    lock_path = DEFAULT_LOCK_PATH
    if "--lock" in args:
        i = args.index("--lock")
        if i + 1 >= len(args):
            print("Error: --lock requires a path", file=sys.stderr)
            sys.exit(1)
        lock_path = safe_path(args[i + 1], {".json"})
        del args[i:i + 2]

    if len(args) != 2:
        print(
            f"Usage: {sys.argv[0]} <input.json> <output.json> [--lock <lockfile.json>]",
            file=sys.stderr,
        )
        sys.exit(1)

    input_path = safe_path(args[0], {".json"})
    output_path = safe_path(args[1], {".json"})

    spec = json.loads(input_path.read_text())

    lock = load_lock(lock_path)
    count, updated_lock = assign_operation_ids(spec, lock)
    new_ids = len(updated_lock) - len(lock)

    patched = apply_request_body_patches(spec)

    output_path.write_text(json.dumps(spec, indent=2))
    write_lock(updated_lock, lock_path)
    print(f"Enriched {count} operations, wrote to {output_path}")
    print(
        f"operationId lock: {len(updated_lock)} entries "
        f"({new_ids} new), wrote to {lock_path}"
    )
    if patched:
        print(f"Injected {patched} missing requestBody object(s)")


if __name__ == "__main__":
    main()
