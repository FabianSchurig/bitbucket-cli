#!/usr/bin/env python3
"""
partition_spec.py: Extract API paths by command group and recursively resolve $refs.

Produces self-contained schema YAML files with no dangling references.

Usage:
  python3 partition_spec.py <input.json> <output.yaml>          # single group (default: pr)
  python3 partition_spec.py <input.json> <output_dir> --all     # all groups
"""

import copy
import json
import sys
from pathlib import Path

import yaml


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


def safe_dir(raw: str) -> Path:
    """Resolve a CLI argument to a canonical directory path."""
    if "\0" in raw:
        raise ValueError(f"Null byte in path: {raw!r}")
    p = Path(raw).resolve()
    if not p.is_dir():
        raise ValueError(f"Not a directory: {p}")
    return p


# ─── Command group definitions ────────────────────────────────────────────────
# Each group produces a separate schema file and Cobra parent command.

COMMAND_GROUPS = {
    "pr": {
        "filename": "pr-schema.yaml",
        "title": "Bitbucket Pull Requests CLI",
        "tags": {"Pullrequests"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/pullrequests",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/merge",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/diff",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/commits",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve",
        ],
        "cli_meta": {
            "x-cli-command-name": "PR",
            "x-cli-command-use": "pr",
            "x-cli-command-short": "Manage Bitbucket pull requests",
            "x-cli-command-long": "Commands for listing, creating, reading, "
                                  "and merging Bitbucket pull requests.",
        },
    },
    "hooks": {
        "filename": "hooks-schema.yaml",
        "title": "Bitbucket Webhooks CLI",
        "tags": {"Webhooks"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/hooks",
            "/repositories/{workspace}/{repo_slug}/hooks/{uid}",
            "/workspaces/{workspace}/hooks",
            "/workspaces/{workspace}/hooks/{uid}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Hooks",
            "x-cli-command-use": "hooks",
            "x-cli-command-short": "Manage Bitbucket webhooks",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and deleting Bitbucket webhooks.",
        },
    },
    "search": {
        "filename": "search-schema.yaml",
        "title": "Bitbucket Code Search CLI",
        "tags": {"Search"},
        "paths": [
            "/users/{selected_user}/search/code",
            "/workspaces/{workspace}/search/code",
            "/teams/{username}/search/code",
        ],
        "cli_meta": {
            "x-cli-command-name": "Search",
            "x-cli-command-use": "search",
            "x-cli-command-short": "Search Bitbucket code",
            "x-cli-command-long": "Commands for searching code across "
                                  "Bitbucket repositories by user, workspace, "
                                  "or team.",
        },
    },
}

HTTP_METHODS = {"get", "post", "put", "patch", "delete"}


def collect_refs(node, spec, collected: set):
    """Walk the node tree, collecting all $ref targets recursively."""
    if isinstance(node, dict):
        if "$ref" in node:
            ref = node["$ref"]
            if ref.startswith("#/") and ref not in collected:
                collected.add(ref)
                # Resolve and recurse
                parts = ref.lstrip("#/").split("/")
                target = spec
                try:
                    for p in parts:
                        target = target[p]
                    collect_refs(target, spec, collected)
                except (KeyError, TypeError):
                    pass  # dangling ref — skip
        else:
            for v in node.values():
                collect_refs(v, spec, collected)
    elif isinstance(node, list):
        for item in node:
            collect_refs(item, spec, collected)


def extract_paths_by_tag(spec: dict, tags: set[str]) -> dict:
    """Extract paths matching the given tags using tag-based filtering."""
    out_paths = {}
    for path, path_item in spec.get("paths", {}).items():
        for method, op in path_item.items():
            if method in HTTP_METHODS and set(op.get("tags", [])) & tags:
                out_paths[path] = copy.deepcopy(path_item)
                break
    return out_paths


def extract_paths_explicit(spec: dict, paths: list[str]) -> dict:
    """Extract paths using an explicit path list."""
    out_paths = {}
    for path in paths:
        if path in spec.get("paths", {}):
            out_paths[path] = copy.deepcopy(spec["paths"][path])
    return out_paths


def build_schema(spec: dict, group: dict) -> dict:
    """Build a self-contained schema for a single command group."""
    version = spec.get("info", {}).get("version", "2.0.0")
    info = {"title": group["title"], "version": version}
    info.update(group["cli_meta"])

    out = {
        "openapi": "3.0.0",
        "info": info,
        "paths": {},
        "components": {"schemas": {}},
    }

    # Try tag-based extraction first, fall back to explicit paths
    paths_by_tag = extract_paths_by_tag(spec, group["tags"])
    paths_explicit = extract_paths_explicit(spec, group["paths"])

    # Merge both (tag-based takes priority, explicit fills gaps)
    out["paths"] = {**paths_explicit, **paths_by_tag}

    if not out["paths"]:
        print(f"Warning: no paths found for group '{group['title']}'",
              file=sys.stderr)

    # Collect all $refs referenced by those paths
    refs: set = set()
    collect_refs(out["paths"], spec, refs)

    # Copy resolved schemas, including all transitively referenced ones
    for ref in refs:
        parts = ref.lstrip("#/").split("/")
        if len(parts) >= 3 and parts[0] == "components" and parts[1] == "schemas":
            schema_name = parts[2]
            if schema_name in spec.get("components", {}).get("schemas", {}):
                out["components"]["schemas"][schema_name] = copy.deepcopy(
                    spec["components"]["schemas"][schema_name]
                )

    # Post-process schema for stable code generation
    post_process_schema(out)

    return out


def write_schema(out: dict, output_path: Path) -> None:
    """Write a schema dict to a YAML file."""
    output_path.write_text(
        yaml.dump(out, default_flow_style=False,
                  sort_keys=False, allow_unicode=True)
    )
    print(
        f"Extracted {len(out['paths'])} paths, "
        f"{len(out['components']['schemas'])} schemas, "
        f"wrote to {output_path}"
    )


def main():
    if len(sys.argv) < 3 or len(sys.argv) > 4:
        print(
            f"Usage: {sys.argv[0]} <input.json> <output.yaml>",
            file=sys.stderr,
        )
        print(
            f"       {sys.argv[0]} <input.json> <output_dir> --all",
            file=sys.stderr,
        )
        sys.exit(1)

    input_path = safe_path(sys.argv[1], {".json"})
    spec = json.loads(input_path.read_text())

    all_mode = len(sys.argv) == 4 and sys.argv[3] == "--all"

    if all_mode:
        output_dir = safe_dir(sys.argv[2])
        for group in COMMAND_GROUPS.values():
            out = build_schema(spec, group)
            write_schema(out, output_dir / group["filename"])
    else:
        output_path = safe_path(sys.argv[2], {".yaml", ".yml"})
        # Determine group from output filename, default to "pr"
        group_key = "pr"
        for key, group in COMMAND_GROUPS.items():
            if output_path.name == group["filename"]:
                group_key = key
                break
        out = build_schema(spec, COMMAND_GROUPS[group_key])
        write_schema(out, output_path)


def post_process_schema(out: dict) -> None:
    """Normalize extracted schema for stable Go code generation.

    Two transformations are applied after extraction:
    1. Ensure ``pullrequest`` always exposes ``description`` inside an allOf
       inline subschema so that oapi-codegen generates it as a Go struct field.
       When the live Bitbucket API schema places ``description`` as a top-level
       property alongside ``allOf``, oapi-codegen silently ignores it; moving
       it into the inline allOf subschema makes it visible to the generator.
    2. Lift the inline ``pullrequest_endpoint.branch`` object into a named
       schema (``pullrequest_endpoint_branch``) so that Go code can reference
       the type by name instead of using an anonymous struct literal that
       breaks whenever Bitbucket adds fields to the branch object.
    """
    schemas = out.get("components", {}).get("schemas", {})

    # 1. Ensure pullrequest.description ends up in the allOf inline subschema.
    # oapi-codegen ignores top-level `properties` when `allOf` is also present.
    # Strategy:
    #   a. If `description` exists as a top-level property alongside `allOf`,
    #      move it into the allOf inline object subschema (type: object).
    #   b. If `description` is missing entirely from both top-level properties
    #      and all allOf subschemas, inject it into the inline allOf subschema
    #      (or as a direct property if no allOf exists).
    pr = schemas.get("pullrequest")
    if pr is not None:
        all_of = pr.get("allOf", [])
        if all_of:
            # Find the inline subschema (type: object) to inject into
            inline_sub = next(
                (s for s in all_of if s.get("type") == "object"), None
            )
            # Collect all properties already visible inside allOf subschemas
            allof_props: set = set()
            for sub in all_of:
                allof_props.update(sub.get("properties", {}).keys())

            # Move top-level description (ignored by oapi-codegen) into allOf
            top_level_desc = pr.get("properties", {}).pop("description", None)
            if top_level_desc and "description" not in allof_props:
                if inline_sub is not None:
                    inline_sub.setdefault("properties", {})[
                        "description"] = top_level_desc
                # Clean up the now-empty top-level properties dict
                if not pr.get("properties"):
                    pr.pop("properties", None)
            elif "description" not in allof_props:
                # description is missing entirely — inject it
                target = inline_sub if inline_sub is not None else pr
                target.setdefault("properties", {})["description"] = {
                    "type": "string",
                    "description": "Explains what the pull request does.",
                }
        else:
            # No allOf — just ensure description exists as a direct property
            if "description" not in pr.get("properties", {}):
                pr.setdefault("properties", {})["description"] = {
                    "type": "string",
                    "description": "Explains what the pull request does.",
                }

    # 2. Lift pullrequest_endpoint.branch to a named schema
    ep = schemas.get("pullrequest_endpoint")
    if ep is not None:
        branch = ep.get("properties", {}).get("branch")
        if branch is not None and "$ref" not in branch and branch.get("type") == "object":
            # Move the inline object to a named schema
            schemas["pullrequest_endpoint_branch"] = branch
            ep["properties"]["branch"] = {
                "$ref": "#/components/schemas/pullrequest_endpoint_branch"
            }


if __name__ == "__main__":
    main()
