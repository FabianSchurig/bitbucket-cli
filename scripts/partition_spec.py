#!/usr/bin/env python3
"""
partition_spec.py: Extract PR paths and recursively resolve all $refs.

Produces a self-contained pr-schema.yaml with no dangling references.

Usage: python3 partition_spec.py <input.json> <output.yaml>
"""

import copy
import json
import sys
from pathlib import Path

import yaml

# Target tags to extract — set to {"Pullrequests"} for initial scope
TARGET_TAGS = {"Pullrequests"}

# Explicit PR paths (fallback if tag-based filtering is insufficient)
PR_PATHS = [
    "/repositories/{workspace}/{repo_slug}/pullrequests",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/merge",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/diff",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/commits",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments",
    "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve",
]

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


def extract_paths_by_tag(spec: dict) -> dict:
    """Extract paths matching TARGET_TAGS using tag-based filtering."""
    out_paths = {}
    for path, path_item in spec.get("paths", {}).items():
        for method, op in path_item.items():
            if method in HTTP_METHODS and set(op.get("tags", [])) & TARGET_TAGS:
                out_paths[path] = copy.deepcopy(path_item)
                break
    return out_paths


def extract_paths_explicit(spec: dict) -> dict:
    """Extract paths using the explicit PR_PATHS list."""
    out_paths = {}
    for path in PR_PATHS:
        if path in spec.get("paths", {}):
            out_paths[path] = copy.deepcopy(spec["paths"][path])
    return out_paths


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input.json> <output.yaml>", file=sys.stderr)
        sys.exit(1)

    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])

    spec = json.loads(input_path.read_text())

    # Build output spec
    version = spec.get("info", {}).get("version", "2.0.0")
    out = {
        "openapi": "3.0.0",
        "info": {"title": "Bitbucket Pull Requests CLI", "version": version},
        "paths": {},
        "components": {"schemas": {}},
    }

    # Try tag-based extraction first, fall back to explicit paths
    paths_by_tag = extract_paths_by_tag(spec)
    paths_explicit = extract_paths_explicit(spec)

    # Merge both (tag-based takes priority, explicit fills gaps)
    out["paths"] = {**paths_explicit, **paths_by_tag}

    if not out["paths"]:
        print("Warning: no PR paths found in spec", file=sys.stderr)

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

    output_path.write_text(
        yaml.dump(out, default_flow_style=False, sort_keys=False, allow_unicode=True)
    )
    print(
        f"Extracted {len(out['paths'])} paths, "
        f"{len(out['components']['schemas'])} schemas, "
        f"wrote to {output_path}"
    )


if __name__ == "__main__":
    main()
