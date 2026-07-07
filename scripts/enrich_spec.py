#!/usr/bin/env python3
"""
enrich_spec.py: Inject operationIds into the Bitbucket OpenAPI spec.

Strategy: slugify(summary) if present, else "{method}_{path_slug}".
This is deterministic and stable as long as Atlassian doesn't rename summaries
(rare) — and the CI diff check will catch it if they do.

Usage: python3 enrich_spec.py <input.json> <output.json>
"""

import copy
import json
import re
import sys
from pathlib import Path


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


def main():
    if len(sys.argv) != 3:
        print(
            f"Usage: {sys.argv[0]} <input.json> <output.json>", file=sys.stderr)
        sys.exit(1)

    input_path = safe_path(sys.argv[1], {".json"})
    output_path = safe_path(sys.argv[2], {".json"})

    spec = json.loads(input_path.read_text())

    http_methods = ["get", "post", "put", "patch", "delete"]
    count = 0

    for path, path_item in spec.get("paths", {}).items():
        for method in http_methods:
            op = path_item.get(method)
            if not op:
                continue
            if "operationId" not in op:
                summary = op.get("summary", "")
                op["operationId"] = to_camel(
                    summary) if summary else path_slug(path, method)
            count += 1

    # Deduplicate operationIds: when two different operations share the same
    # operationId (e.g. two "List a pull request activity log" endpoints at
    # different paths), regenerate the duplicate using path_slug which
    # incorporates the full URL path and is therefore unique.
    seen: dict[str, tuple[str, str]] = {}  # operationId → (path, method)
    for path, path_item in spec.get("paths", {}).items():
        for method in http_methods:
            op = path_item.get(method)
            if not op or "operationId" not in op:
                continue
            oid = op["operationId"]
            if oid in seen:
                # Collision — regenerate this one using path-based slug
                op["operationId"] = path_slug(path, method)
            else:
                seen[oid] = (path, method)

    patched = apply_request_body_patches(spec)

    output_path.write_text(json.dumps(spec, indent=2))
    print(f"Enriched {count} operations, wrote to {output_path}")
    if patched:
        print(f"Injected {patched} missing requestBody object(s)")


if __name__ == "__main__":
    main()
