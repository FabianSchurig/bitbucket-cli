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
    "refs": {
        "filename": "refs-schema.yaml",
        "title": "Bitbucket Refs CLI",
        "tags": {"Refs"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/refs",
            "/repositories/{workspace}/{repo_slug}/refs/branches",
            "/repositories/{workspace}/{repo_slug}/refs/branches/{name}",
            "/repositories/{workspace}/{repo_slug}/refs/tags",
            "/repositories/{workspace}/{repo_slug}/refs/tags/{name}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Refs",
            "x-cli-command-use": "refs",
            "x-cli-command-short": "Manage Bitbucket branches and tags",
            "x-cli-command-long": "Commands for listing, creating, and deleting "
                                  "branches and tags in Bitbucket repositories.",
        },
    },
    "commits": {
        "filename": "commits-schema.yaml",
        "title": "Bitbucket Commits CLI",
        "tags": set(),
        "paths": [
            "/repositories/{workspace}/{repo_slug}/commits",
            "/repositories/{workspace}/{repo_slug}/commits/{revision}",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/approve",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/comments",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/comments/{comment_id}",
            "/repositories/{workspace}/{repo_slug}/diff/{spec}",
            "/repositories/{workspace}/{repo_slug}/diffstat/{spec}",
            "/repositories/{workspace}/{repo_slug}/merge-base/{revspec}",
            "/repositories/{workspace}/{repo_slug}/patch/{spec}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Commits",
            "x-cli-command-use": "commits",
            "x-cli-command-short": "Manage Bitbucket commits",
            "x-cli-command-long": "Commands for listing, viewing, approving, "
                                  "and commenting on commits in Bitbucket "
                                  "repositories.",
        },
    },
    "reports": {
        "filename": "reports-schema.yaml",
        "title": "Bitbucket Reports CLI",
        "tags": set(),
        "paths": [
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/reports",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}/annotations",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}/annotations/{annotationId}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Reports",
            "x-cli-command-use": "reports",
            "x-cli-command-short": "Manage Bitbucket commit reports and annotations",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and deleting commit reports and annotations "
                                  "in Bitbucket repositories.",
        },
    },
    "repositories": {
        "filename": "repositories-schema.yaml",
        "title": "Bitbucket Repositories CLI",
        "tags": {"Repositories", "Source"},
        "paths": [
            "/repositories",
            "/repositories/{workspace}",
            "/repositories/{workspace}/{repo_slug}",
            "/repositories/{workspace}/{repo_slug}/filehistory/{commit}/{path}",
            "/repositories/{workspace}/{repo_slug}/forks",
            "/repositories/{workspace}/{repo_slug}/override-settings",
            "/repositories/{workspace}/{repo_slug}/permissions-config/groups",
            "/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}",
            "/repositories/{workspace}/{repo_slug}/permissions-config/users",
            "/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}",
            "/repositories/{workspace}/{repo_slug}/src",
            "/repositories/{workspace}/{repo_slug}/src/{commit}/{path}",
            "/repositories/{workspace}/{repo_slug}/watchers",
            "/user/permissions/repositories",
            "/user/workspaces/{workspace}/permissions/repositories",
        ],
        "cli_meta": {
            "x-cli-command-name": "Repos",
            "x-cli-command-use": "repos",
            "x-cli-command-short": "Manage Bitbucket repositories",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and managing Bitbucket repositories, "
                                  "including forks, permissions, and source files.",
        },
    },
    "workspaces": {
        "filename": "workspaces-schema.yaml",
        "title": "Bitbucket Workspaces CLI",
        "tags": {"Workspaces"},
        "paths": [
            "/user/permissions/workspaces",
            "/user/workspaces",
            "/user/workspaces/{workspace}/permission",
            "/workspaces",
            "/workspaces/{workspace}",
            "/workspaces/{workspace}/members",
            "/workspaces/{workspace}/members/{member}",
            "/workspaces/{workspace}/permissions",
            "/workspaces/{workspace}/permissions/repositories",
            "/workspaces/{workspace}/permissions/repositories/{repo_slug}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Workspaces",
            "x-cli-command-use": "workspaces",
            "x-cli-command-short": "Manage Bitbucket workspaces",
            "x-cli-command-long": "Commands for listing workspaces and managing "
                                  "workspace members and permissions.",
        },
    },
    "projects": {
        "filename": "projects-schema.yaml",
        "title": "Bitbucket Projects CLI",
        "tags": {"Projects"},
        "paths": [
            "/workspaces/{workspace}/projects",
            "/workspaces/{workspace}/projects/{project_key}",
            "/workspaces/{workspace}/projects/{project_key}/default-reviewers",
            "/workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}",
            "/workspaces/{workspace}/projects/{project_key}/permissions-config/groups",
            "/workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}",
            "/workspaces/{workspace}/projects/{project_key}/permissions-config/users",
            "/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Projects",
            "x-cli-command-use": "projects",
            "x-cli-command-short": "Manage Bitbucket projects",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and deleting projects and managing project "
                                  "default reviewers and permissions.",
        },
    },
    "pipelines": {
        "filename": "pipelines-schema.yaml",
        "title": "Bitbucket Pipelines CLI",
        "tags": {"Pipelines"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/pipelines",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/log",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/logs/{log_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/test_reports",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/test_reports/test_cases",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/steps/{step_uuid}/test_reports/test_cases/{test_case_uuid}/test_case_reasons",
            "/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}/stopPipeline",
            "/repositories/{workspace}/{repo_slug}/pipelines_config",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/build_number",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/schedules",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}/executions",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/variables",
            "/repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines-config/caches",
            "/repositories/{workspace}/{repo_slug}/pipelines-config/caches/{cache_uuid}",
            "/repositories/{workspace}/{repo_slug}/pipelines-config/caches/{cache_uuid}/content-uri",
            "/repositories/{workspace}/{repo_slug}/pipelines-config/runners",
            "/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}",
            "/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables",
            "/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}",
            "/teams/{username}/pipelines_config/variables",
            "/teams/{username}/pipelines_config/variables/{variable_uuid}",
            "/users/{selected_user}/pipelines_config/variables",
            "/users/{selected_user}/pipelines_config/variables/{variable_uuid}",
            "/workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration",
            "/workspaces/{workspace}/pipelines-config/identity/oidc/keys.json",
            "/workspaces/{workspace}/pipelines-config/runners",
            "/workspaces/{workspace}/pipelines-config/runners/{runner_uuid}",
            "/workspaces/{workspace}/pipelines-config/variables",
            "/workspaces/{workspace}/pipelines-config/variables/{variable_uuid}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Pipelines",
            "x-cli-command-use": "pipelines",
            "x-cli-command-short": "Manage Bitbucket Pipelines",
            "x-cli-command-long": "Commands for listing, triggering, and managing "
                                  "Bitbucket Pipelines, including steps, logs, "
                                  "variables, runners, caches, and schedules.",
        },
    },
    "issues": {
        "filename": "issues-schema.yaml",
        "title": "Bitbucket Issues CLI",
        "tags": {"Issue tracker"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/issues",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/attachments",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/attachments/{path}",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/changes",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/changes/{change_id}",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/vote",
            "/repositories/{workspace}/{repo_slug}/issues/{issue_id}/watch",
            "/repositories/{workspace}/{repo_slug}/issues/export",
            "/repositories/{workspace}/{repo_slug}/issues/export/{repo_name}-issues-{task_id}.zip",
            "/repositories/{workspace}/{repo_slug}/issues/import",
            "/repositories/{workspace}/{repo_slug}/components",
            "/repositories/{workspace}/{repo_slug}/components/{component_id}",
            "/repositories/{workspace}/{repo_slug}/milestones",
            "/repositories/{workspace}/{repo_slug}/milestones/{milestone_id}",
            "/repositories/{workspace}/{repo_slug}/versions",
            "/repositories/{workspace}/{repo_slug}/versions/{version_id}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Issues",
            "x-cli-command-use": "issues",
            "x-cli-command-short": "Manage Bitbucket issues",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and managing issues, comments, attachments, "
                                  "components, milestones, and versions.",
        },
    },
    "snippets": {
        "filename": "snippets-schema.yaml",
        "title": "Bitbucket Snippets CLI",
        "tags": {"Snippets"},
        "paths": [
            "/snippets",
            "/snippets/{workspace}",
            "/snippets/{workspace}/{encoded_id}",
            "/snippets/{workspace}/{encoded_id}/comments",
            "/snippets/{workspace}/{encoded_id}/comments/{comment_id}",
            "/snippets/{workspace}/{encoded_id}/commits",
            "/snippets/{workspace}/{encoded_id}/commits/{revision}",
            "/snippets/{workspace}/{encoded_id}/files/{path}",
            "/snippets/{workspace}/{encoded_id}/watch",
            "/snippets/{workspace}/{encoded_id}/watchers",
            "/snippets/{workspace}/{encoded_id}/{node_id}",
            "/snippets/{workspace}/{encoded_id}/{node_id}/files/{path}",
            "/snippets/{workspace}/{encoded_id}/{revision}/diff",
            "/snippets/{workspace}/{encoded_id}/{revision}/patch",
        ],
        "cli_meta": {
            "x-cli-command-name": "Snippets",
            "x-cli-command-use": "snippets",
            "x-cli-command-short": "Manage Bitbucket snippets",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and deleting snippets, including comments, "
                                  "commits, and file operations.",
        },
    },
    "deployments": {
        "filename": "deployments-schema.yaml",
        "title": "Bitbucket Deployments CLI",
        "tags": {"Deployments"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/deploy-keys",
            "/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}",
            "/repositories/{workspace}/{repo_slug}/deployments",
            "/repositories/{workspace}/{repo_slug}/deployments/{deployment_uuid}",
            "/repositories/{workspace}/{repo_slug}/environments",
            "/repositories/{workspace}/{repo_slug}/environments/{environment_uuid}",
            "/repositories/{workspace}/{repo_slug}/environments/{environment_uuid}/changes",
            "/workspaces/{workspace}/projects/{project_key}/deploy-keys",
            "/workspaces/{workspace}/projects/{project_key}/deploy-keys/{key_id}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Deployments",
            "x-cli-command-use": "deployments",
            "x-cli-command-short": "Manage Bitbucket deployments",
            "x-cli-command-long": "Commands for managing deploy keys, "
                                  "deployments, and environments in "
                                  "Bitbucket repositories and projects.",
        },
    },
    "branch-restrictions": {
        "filename": "branch-restrictions-schema.yaml",
        "title": "Bitbucket Branch Restrictions CLI",
        "tags": {"Branch restrictions"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/branch-restrictions",
            "/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}",
        ],
        "cli_meta": {
            "x-cli-command-name": "BranchRestrictions",
            "x-cli-command-use": "branch-restrictions",
            "x-cli-command-short": "Manage Bitbucket branch restrictions",
            "x-cli-command-long": "Commands for listing, creating, updating, "
                                  "and deleting branch restriction rules in "
                                  "Bitbucket repositories.",
        },
    },
    "branching-model": {
        "filename": "branching-model-schema.yaml",
        "title": "Bitbucket Branching Model CLI",
        "tags": {"Branching model"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/branching-model",
            "/repositories/{workspace}/{repo_slug}/branching-model/settings",
            "/repositories/{workspace}/{repo_slug}/effective-branching-model",
            "/workspaces/{workspace}/projects/{project_key}/branching-model",
            "/workspaces/{workspace}/projects/{project_key}/branching-model/settings",
        ],
        "cli_meta": {
            "x-cli-command-name": "BranchingModel",
            "x-cli-command-use": "branching-model",
            "x-cli-command-short": "Manage Bitbucket branching models",
            "x-cli-command-long": "Commands for viewing and updating branching "
                                  "model settings for repositories and projects.",
        },
    },
    "commit-statuses": {
        "filename": "commit-statuses-schema.yaml",
        "title": "Bitbucket Commit Statuses CLI",
        "tags": {"Commit statuses"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/statuses",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/statuses/build",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/statuses/build/{key}",
        ],
        "cli_meta": {
            "x-cli-command-name": "CommitStatuses",
            "x-cli-command-use": "commit-statuses",
            "x-cli-command-short": "Manage Bitbucket commit statuses",
            "x-cli-command-long": "Commands for listing, creating, and updating "
                                  "build statuses on commits in Bitbucket "
                                  "repositories.",
        },
    },
    "downloads": {
        "filename": "downloads-schema.yaml",
        "title": "Bitbucket Downloads CLI",
        "tags": {"Downloads"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/downloads",
            "/repositories/{workspace}/{repo_slug}/downloads/{filename}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Downloads",
            "x-cli-command-use": "downloads",
            "x-cli-command-short": "Manage Bitbucket repository downloads",
            "x-cli-command-long": "Commands for listing, uploading, and deleting "
                                  "download artifacts in Bitbucket repositories.",
        },
    },
    "users": {
        "filename": "users-schema.yaml",
        "title": "Bitbucket Users CLI",
        "tags": {"Users", "SSH", "GPG"},
        "paths": [
            "/user",
            "/user/emails",
            "/user/emails/{email}",
            "/users/{selected_user}",
            "/users/{selected_user}/ssh-keys",
            "/users/{selected_user}/ssh-keys/{key_id}",
            "/users/{selected_user}/gpg-keys",
            "/users/{selected_user}/gpg-keys/{fingerprint}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Users",
            "x-cli-command-use": "users",
            "x-cli-command-short": "Manage Bitbucket users",
            "x-cli-command-long": "Commands for viewing user profiles, "
                                  "managing email addresses, SSH keys, "
                                  "and GPG keys.",
        },
    },
    "properties": {
        "filename": "properties-schema.yaml",
        "title": "Bitbucket Properties CLI",
        "tags": {"properties"},
        "paths": [
            "/repositories/{workspace}/{repo_slug}/properties/{app_key}/{property_name}",
            "/repositories/{workspace}/{repo_slug}/commit/{commit}/properties/{app_key}/{property_name}",
            "/repositories/{workspace}/{repo_slug}/pullrequests/{pullrequest_id}/properties/{app_key}/{property_name}",
            "/users/{selected_user}/properties/{app_key}/{property_name}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Properties",
            "x-cli-command-use": "properties",
            "x-cli-command-short": "Manage Bitbucket application properties",
            "x-cli-command-long": "Commands for getting, updating, and deleting "
                                  "application properties on repositories, "
                                  "commits, pull requests, and users.",
        },
    },
    "addon": {
        "filename": "addon-schema.yaml",
        "title": "Bitbucket Addon CLI",
        "tags": {"Addon"},
        "paths": [
            "/addon",
            "/addon/linkers",
            "/addon/linkers/{linker_key}",
            "/addon/linkers/{linker_key}/values",
            "/addon/linkers/{linker_key}/values/{value_id}",
        ],
        "cli_meta": {
            "x-cli-command-name": "Addon",
            "x-cli-command-use": "addon",
            "x-cli-command-short": "Manage Bitbucket Connect addons",
            "x-cli-command-long": "Commands for managing Bitbucket Connect "
                                  "addon installations and linker values.",
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
    # Sort paths alphabetically for deterministic output across runs.
    merged = {**paths_explicit, **paths_by_tag}
    out["paths"] = dict(sorted(merged.items()))

    if not out["paths"]:
        print(f"Warning: no paths found for group '{group['title']}'",
              file=sys.stderr)

    # Collect all $refs referenced by those paths
    refs: set = set()
    collect_refs(out["paths"], spec, refs)

    # Copy resolved schemas, including all transitively referenced ones.
    # Sort refs for deterministic insertion order (Python sets are unordered).
    for ref in sorted(refs):
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
