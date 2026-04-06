#!/usr/bin/env python3
"""Generate a migration guide from the legacy terraform-provider-bitbucket."""

import argparse
import html
import re
import ssl
import sys
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path
from typing import Iterable
from urllib.error import HTTPError, URLError

LEGACY_BASE = "https://raw.githubusercontent.com/DrFaust92/terraform-provider-bitbucket/master"

CURRENT_KIND_PATH = {
    "resource": "resources",
    "data-source": "data-sources",
}

LEGACY_KIND_PREFIX = {
    "resource": "resource_",
    "data-source": "data_",
}

CURRENT_ALIAS = {
    ("resource", "bitbucket_branch_restriction"): ["bitbucket_branch_restrictions"],
    ("resource", "bitbucket_branching_model"): ["bitbucket_branching_model"],
    ("resource", "bitbucket_commit_file"): ["bitbucket_commit_file"],
    ("resource", "bitbucket_default_reviewers"): ["bitbucket_default_reviewers"],
    ("resource", "bitbucket_deploy_key"): ["bitbucket_repo_deploy_keys"],
    ("resource", "bitbucket_deployment"): ["bitbucket_deployments"],
    ("resource", "bitbucket_deployment_variable"): ["bitbucket_deployment_variables"],
    ("resource", "bitbucket_forked_repository"): ["bitbucket_forked_repository"],
    ("resource", "bitbucket_group"): [],
    ("resource", "bitbucket_group_membership"): [],
    ("resource", "bitbucket_hook"): ["bitbucket_hooks"],
    ("resource", "bitbucket_pipeline_schedule"): ["bitbucket_pipeline_schedules"],
    ("resource", "bitbucket_pipeline_ssh_key"): ["bitbucket_pipeline_ssh_keys"],
    ("resource", "bitbucket_pipeline_ssh_known_host"): ["bitbucket_pipeline_known_hosts"],
    ("resource", "bitbucket_project"): ["bitbucket_projects"],
    ("resource", "bitbucket_project_branching_model"): ["bitbucket_project_branching_model"],
    ("resource", "bitbucket_project_default_reviewers"): ["bitbucket_project_default_reviewers"],
    ("resource", "bitbucket_project_group_permission"): ["bitbucket_project_group_permissions"],
    ("resource", "bitbucket_project_user_permission"): ["bitbucket_project_user_permissions"],
    ("resource", "bitbucket_repository"): [
        "bitbucket_repos",
        "bitbucket_repo_settings",
        "bitbucket_pipeline_config",
    ],
    ("resource", "bitbucket_repository_group_permission"): ["bitbucket_repo_group_permissions"],
    ("resource", "bitbucket_repository_user_permission"): ["bitbucket_repo_user_permissions"],
    ("resource", "bitbucket_repository_variable"): ["bitbucket_pipeline_variables"],
    ("resource", "bitbucket_ssh_key"): ["bitbucket_ssh_keys"],
    ("resource", "bitbucket_workspace_hook"): ["bitbucket_workspace_hooks"],
    ("resource", "bitbucket_workspace_variable"): ["bitbucket_workspace_pipeline_variables"],
    ("data-source", "bitbucket_current_user"): ["bitbucket_current_user"],
    ("data-source", "bitbucket_deployment"): ["bitbucket_deployments"],
    ("data-source", "bitbucket_deployments"): ["bitbucket_deployments"],
    ("data-source", "bitbucket_file"): ["bitbucket_commit_file"],
    ("data-source", "bitbucket_group"): [],
    ("data-source", "bitbucket_group_members"): [],
    ("data-source", "bitbucket_groups"): [],
    ("data-source", "bitbucket_hook_types"): ["bitbucket_hook_types"],
    ("data-source", "bitbucket_ip_ranges"): [],
    ("data-source", "bitbucket_pipeline_oidc_config"): ["bitbucket_pipeline_oidc"],
    ("data-source", "bitbucket_pipeline_oidc_config_keys"): ["bitbucket_pipeline_oidc_keys"],
    ("data-source", "bitbucket_project"): ["bitbucket_projects"],
    ("data-source", "bitbucket_repository"): ["bitbucket_repos"],
    ("data-source", "bitbucket_user"): ["bitbucket_users"],
    ("data-source", "bitbucket_workspace"): ["bitbucket_workspaces"],
    ("data-source", "bitbucket_workspace_members"): ["bitbucket_workspace_members"],
}

OBJECT_NOTES = {
    ("resource", "bitbucket_repository"): (
        "The legacy repository resource bundled core repository CRUD, pipeline "
        "enablement, and override-settings flags. In the new provider, core CRUD "
        "stays on `bitbucket_repos`, pipeline enablement moves to "
        "`bitbucket_pipeline_config`, and repository settings have their own "
        "`bitbucket_repo_settings` resource."
    ),
    ("resource", "bitbucket_repository_variable"): (
        "Legacy repository variables map to the pipelines variable API. Use "
        "`bitbucket_pipeline_variables` and rename `owner`/`repository` to "
        "`workspace`/`repo_slug`."
    ),
    ("resource", "bitbucket_workspace_variable"): (
        "Workspace variables now live under the pipelines API as "
        "`bitbucket_workspace_pipeline_variables`."
    ),
    ("resource", "bitbucket_deploy_key"): (
        "The generated provider exposes deploy keys as `bitbucket_repo_deploy_keys` "
        "and also has separate project-level deploy key resources."
    ),
    ("data-source", "bitbucket_file"): (
        "The legacy `bitbucket_file` data source maps most closely to "
        "`bitbucket_commit_file`, which reads file content via the commit-file "
        "endpoint."
    ),
    ("data-source", "bitbucket_current_user"): (
        "The legacy data source also fetched `/user/emails`. The generated provider "
        "splits that into `bitbucket_current_user` plus `bitbucket_user_emails` when "
        "you need email addresses."
    ),
    ("data-source", "bitbucket_deployment"): (
        "Use `bitbucket_deployments` with the identifying path parameters for a "
        "single deployment; omit the single-resource expectation and treat the "
        "response as the generic deployment payload."
    ),
    ("resource", "bitbucket_group"): (
        "Workspace group management is not currently exposed by the generated "
        "provider because those endpoints are not represented in the generated "
        "Terraform docs."
    ),
    ("resource", "bitbucket_group_membership"): (
        "Group membership management is not currently exposed by the generated "
        "provider."
    ),
    ("data-source", "bitbucket_group"): (
        "Group lookup is not currently exposed by the generated provider."
    ),
    ("data-source", "bitbucket_group_members"): (
        "Group member lookup is not currently exposed by the generated provider."
    ),
    ("data-source", "bitbucket_groups"): (
        "Group listing is not currently exposed by the generated provider."
    ),
    ("data-source", "bitbucket_ip_ranges"): (
        "The generated provider does not currently expose Bitbucket IP ranges as a "
        "Terraform data source."
    ),
}

COMMON_RENAMES = [
    (
        "Provider `password`",
        "Provider `token`",
        "The new provider only accepts `token`; `BITBUCKET_PASSWORD` is replaced by "
        "`BITBUCKET_TOKEN`.",
    ),
    (
        "Provider `oauth_client_id`, `oauth_client_secret`, `oauth_token`",
        "No direct equivalent",
        "The generated provider currently supports API tokens and "
        "workspace/repository access tokens only.",
    ),
    (
        "`owner`",
        "`workspace`",
        "Most repository/project scoped resources renamed the workspace path "
        "parameter to match Bitbucket Cloud OpenAPI naming.",
    ),
    (
        "`repository` or legacy repository name/slug fields",
        "`repo_slug`",
        "The generated provider consistently uses the Bitbucket path parameter name "
        "`repo_slug`.",
    ),
    (
        "Singular resource names like `bitbucket_repository`",
        "Plural/group-based names like `bitbucket_repos`",
        "Generated resources follow API operation groups instead of the legacy "
        "hand-written naming scheme.",
    ),
]

PARAM_RENAMES = {
    "owner": "workspace",
    "repository": "repo_slug",
}

PLACEHOLDER_TOKENS = {
    ("Workspace",): "{workspace}",
    ("Repo", "Slug"): "{repo_slug}",
    ("Project", "Key"): "{project_key}",
    ("Username",): "{username}",
    ("Selected", "User"): "{selected_user}",
    ("Node",): "{node}",
    ("Id",): "{id}",
    ("Key",): "{key}",
    ("Target",): "{target}",
    ("Path",): "{path}",
    ("Commit",): "{commit}",
    ("Environment", "Uuid"): "{environment_uuid}",
    ("Variable", "Uuid"): "{variable_uuid}",
    ("Hook", "Uuid"): "{hook_uuid}",
    ("Schedule", "Uuid"): "{schedule_uuid}",
    ("Runner", "Uuid"): "{runner_uuid}",
    ("Email",): "{email}",
}

HTTP_VERBS = {
    "Get": "GET",
    "Post": "POST",
    "Put": "PUT",
    "Delete": "DELETE",
    "Patch": "PATCH",
}

CAMEL_RE = re.compile(r"[A-Z]+(?=$|[A-Z][a-z0-9])|[A-Z]?[a-z0-9]+")


@dataclass
class DocObject:
    kind: str
    name: str
    title: str = ""
    inputs_required: list[str] = field(default_factory=list)
    inputs_optional: list[str] = field(default_factory=list)
    read_only: list[str] = field(default_factory=list)
    endpoints: list[str] = field(default_factory=list)

    @property
    def inputs(self) -> list[str]:
        return dedupe(self.inputs_required + self.inputs_optional)


def dedupe(items: Iterable[str]) -> list[str]:
    seen = set()
    result = []
    for item in items:
        if item and item not in seen:
            seen.add(item)
            result.append(item)
    return result


def fetch(url: str, *, required: bool = True, timeout: int = 20) -> str | None:
    request = urllib.request.Request(url, headers={"User-Agent": "bitbucket-cli-migration-generator"})
    try:
        with urllib.request.urlopen(
            request,
            timeout=timeout,
            context=ssl.create_default_context(),
        ) as response:
            return response.read().decode("utf-8")
    except (HTTPError, URLError, TimeoutError) as error:
        if required:
            raise RuntimeError(f"failed to fetch {url}: {error}") from error
        print(f"warning: skipping unavailable URL {url}: {error}", file=sys.stderr)
        return None


def parse_bullets(section: str, current: bool) -> list[str]:
    """Parse top-level markdown bullets and intentionally ignore nested schema bullets."""
    items = []
    for line in section.splitlines():
        line = line.rstrip()
        if current:
            match = re.match(r"- `([^`]+)`", line)
        else:
            match = re.match(r"\* `([^`]+)`", line.strip())
        if match:
            items.append(match.group(1))
    return items


def split_section(text: str, start: str, stops: list[str]) -> str:
    if start not in text:
        return ""
    tail = text.split(start, 1)[1]
    end = len(tail)
    for stop in stops:
        index = tail.find(stop)
        if index != -1:
            end = min(end, index)
    return tail[:end]


def parse_current_doc(path: Path, kind: str) -> DocObject:
    text = path.read_text()
    title_match = re.search(r"^#\s+(.+)$", text, re.M)
    if title_match is None:
        raise ValueError(f"missing Terraform object heading in {path}")
    title = title_match.group(1)
    name_match = re.search(r"(bitbucket_[^\s]+)", title)
    if name_match is None:
        raise ValueError(f"missing Terraform object name in heading for {path}")
    name = name_match.group(1)
    doc = DocObject(kind=kind, name=name, title=title)
    doc.inputs_required = parse_bullets(
        split_section(text, "### Required", ["### Optional", "### Read-Only", "##"]),
        True,
    )
    doc.inputs_optional = parse_bullets(
        split_section(text, "### Optional", ["### Read-Only", "##"]),
        True,
    )
    doc.read_only = parse_bullets(split_section(text, "### Read-Only", ["##"]), True)
    for line in text.splitlines():
        if (
            line.startswith("| ")
            and "`" in line
            and not line.startswith("| Operation")
            and not line.startswith("|-----------")
        ):
            cols = [col.strip() for col in line.strip("|").split("|")]
            if len(cols) >= 3 and cols[1].startswith("`"):
                doc.endpoints.append(
                    f"{cols[0]} {cols[1].strip('`')} {cols[2].strip('`')}"
                )
    doc.endpoints = dedupe(doc.endpoints)
    return doc


def parse_legacy_doc(kind: str, name: str) -> DocObject:
    base = name.removeprefix("bitbucket_")
    doc_url = f"{LEGACY_BASE}/docs/{CURRENT_KIND_PATH[kind]}/{base}.md"
    text = fetch(doc_url, required=False)
    if text is None:
        return DocObject(kind=kind, name=name, title=name)
    title_match = re.search(r"^#\s+([^\n]+)", text, re.M)
    title = name
    if title_match:
        title = html.unescape(title_match.group(1).replace("\\_", "_"))
    doc = DocObject(kind=kind, name=name, title=title)
    args = split_section(
        text,
        "## Argument Reference",
        ["## Attributes Reference", "## Import", "## Attributes", "### "],
    )
    attrs = split_section(text, "## Attributes Reference", ["## Import", "### "])
    for line in args.splitlines():
        match = re.match(r"\* `([^`]+)` - \((Required|Optional)\)", line.strip())
        if not match:
            continue
        if match.group(2) == "Required":
            doc.inputs_required.append(match.group(1))
        else:
            doc.inputs_optional.append(match.group(1))
    doc.read_only = parse_bullets(attrs, False)
    return doc


def source_filename(kind: str, name: str) -> str:
    return f"bitbucket/{LEGACY_KIND_PREFIX[kind]}{name.removeprefix('bitbucket_')}.go"


def tokens_from_method(method_name: str) -> tuple[str, str] | None:
    parts = CAMEL_RE.findall(method_name)
    if not parts or parts[-1] not in HTTP_VERBS:
        return None
    verb = HTTP_VERBS[parts[-1]]
    parts = parts[:-1]
    segments = []
    static = []
    index = 0
    while index < len(parts):
        matched = False
        for size in (2, 1):
            key = tuple(parts[index : index + size])
            if key in PLACEHOLDER_TOKENS:
                if static:
                    segments.append("-".join(token.lower() for token in static))
                    static = []
                segments.append(PLACEHOLDER_TOKENS[key])
                index += size
                matched = True
                break
        if matched:
            continue
        static.append(parts[index])
        index += 1
    if static:
        segments.append("-".join(token.lower() for token in static))
    path = "/" + "/".join(segment for segment in segments if segment)
    path = path.replace("/-/", "/")
    path = re.sub(r"/2-0", "", path)
    path = re.sub(r"//+", "/", path)
    return verb, path


def normalize_raw_path(path: str) -> str:
    path = path.strip().strip('"')
    path = path.replace("2.0/", "/")
    path = path.replace("%s", "{param}")
    path = re.sub(r"//+", "/", path)
    if not path.startswith("/"):
        path = "/" + path
    return path


def parse_legacy_endpoints(kind: str, name: str, mapped_current: list[DocObject]) -> list[str]:
    text = fetch(f"{LEGACY_BASE}/{source_filename(kind, name)}", required=False)
    if text is None:
        return dedupe(sum((doc.endpoints for doc in mapped_current), []))
    endpoints = []
    for match in re.finditer(r"\.([A-Z][A-Za-z0-9]+(?:Get|Post|Put|Delete|Patch))\(", text):
        converted = tokens_from_method(match.group(1))
        if converted:
            endpoints.append(f"{converted[0]} {converted[1]}")
    raw_patterns = [
        ("GET", r'\.Get\("(2\.0/[^"]+)"\)'),
        ("PUT", r'\.Put\(fmt\.Sprintf\("(2\.0/[^"]+)"'),
        ("GET", r'\.Get\(fmt\.Sprintf\("(2\.0/[^"]+)"'),
        ("POST", r'\.Post\(fmt\.Sprintf\("(2\.0/[^"]+)"'),
        ("DELETE", r'\.Delete\(fmt\.Sprintf\("(2\.0/[^"]+)"'),
    ]
    for verb, pattern in raw_patterns:
        for match in re.finditer(pattern, text):
            endpoints.append(f"{verb} {normalize_raw_path(match.group(1))}")
    endpoints = dedupe(endpoints)
    if mapped_current:
        mapped = dedupe(sum((doc.endpoints for doc in mapped_current), []))
        if any("{param}" in endpoint for endpoint in endpoints):
            verbs = {
                endpoint.split(" ", 1)[0]
                for endpoint in endpoints
                if "{param}" in endpoint
            }
            endpoints = [
                endpoint
                for endpoint in endpoints
                if "{param}" not in endpoint or endpoint.split(" ", 1)[0] not in verbs
            ]
            endpoints = dedupe(
                endpoints
                + [endpoint for endpoint in mapped if endpoint.split(" ", 1)[0] in verbs]
            )
        if not endpoints:
            return mapped
    return endpoints


def current_objects(repo_root: Path) -> dict[tuple[str, str], DocObject]:
    docs = {}
    for kind, subdir in CURRENT_KIND_PATH.items():
        for path in sorted((repo_root / "docs" / subdir).glob("*.md")):
            doc = parse_current_doc(path, kind)
            docs[(kind, doc.name)] = doc
    return docs


def legacy_names() -> dict[str, list[str]]:
    provider = fetch(f"{LEGACY_BASE}/bitbucket/provider.go")
    return {
        "resource": re.findall(r'"(bitbucket_[^"]+)":\s+resource', provider),
        "data-source": re.findall(r'"(bitbucket_[^"]+)":\s+data', provider),
    }


def mapped_current_objects(
    kind: str, name: str, current: dict[tuple[str, str], DocObject]
) -> list[DocObject]:
    mapped_names = CURRENT_ALIAS[(kind, name)]
    return [
        current[(kind, mapped_name)]
        for mapped_name in mapped_names
        if (kind, mapped_name) in current
    ]


def format_params(required: list[str], optional: list[str]) -> str:
    parts = []
    if required:
        parts.append("required: " + ", ".join(f"`{item}`" for item in required))
    if optional:
        parts.append("optional: " + ", ".join(f"`{item}`" for item in optional))
    return "; ".join(parts) if parts else "none"


def format_endpoints(items: list[str]) -> str:
    return "<br>".join(f"`{item}`" for item in items) if items else "none"


def normalized_params(params: Iterable[str]) -> set[str]:
    normalized = set()
    for param in params:
        normalized.add(PARAM_RENAMES.get(param, param))
    return normalized


def diff_summary(legacy: DocObject, currents: list[DocObject]) -> str:
    current_inputs = dedupe(sum((doc.inputs for doc in currents), []))
    legacy_norm = normalized_params(legacy.inputs)
    current_set = set(current_inputs)
    renamed = []
    for source, target in PARAM_RENAMES.items():
        if source in legacy.inputs and target in current_set:
            renamed.append(f"`{source}` → `{target}`")
    dropped = sorted(
        param for param in legacy.inputs if PARAM_RENAMES.get(param, param) not in current_set
    )
    added = sorted(param for param in current_inputs if param not in legacy_norm)
    parts = []
    if renamed:
        parts.append("renamed: " + ", ".join(renamed))
    if dropped:
        parts.append("legacy-only inputs: " + ", ".join(f"`{item}`" for item in dropped))
    if added:
        parts.append("new-only inputs: " + ", ".join(f"`{item}`" for item in added))
    return "; ".join(parts) if parts else "input names are effectively unchanged"


def make_overview(
    kind: str, legacy: list[str], current: dict[tuple[str, str], DocObject]
) -> tuple[list[str], list[str], list[str]]:
    matched = []
    legacy_only = []
    mapped_current = set()
    for name in legacy:
        targets = CURRENT_ALIAS[(kind, name)]
        if targets:
            matched.append(name)
            mapped_current.update(targets)
        else:
            legacy_only.append(name)
    current_names = sorted(name for doc_kind, name in current if doc_kind == kind)
    new_only = [name for name in current_names if name not in mapped_current]
    return matched, legacy_only, new_only


def render(repo_root: Path) -> str:
    current = current_objects(repo_root)
    legacy = legacy_names()
    res_matched, res_legacy_only, res_new_only = make_overview(
        "resource", legacy["resource"], current
    )
    ds_matched, ds_legacy_only, ds_new_only = make_overview(
        "data-source", legacy["data-source"], current
    )

    lines = []
    lines.append("# Migration from `DrFaust92/terraform-provider-bitbucket`")
    lines.append("")
    lines.append(
        "This guide compares the legacy hand-written provider with the generated "
        "`FabianSchurig/bitbucket` provider in this repository."
    )
    lines.append("")
    lines.append(
        "It is intentionally a best-effort migration baseline: the generated docs "
        "sometimes list optional fields that are also computed by the API, so subtle "
        "cases still need manual review."
    )
    lines.append("")
    lines.append(
        "It was generated with `python3 scripts/gen_migration.py --output MIGRATION.md`, "
        "using:"
    )
    lines.append("")
    lines.append("- current docs from `./docs/`")
    lines.append(
        "- legacy docs and source from "
        "`https://github.com/DrFaust92/terraform-provider-bitbucket/tree/master`"
    )
    lines.append("")
    lines.append("## What changes first")
    lines.append("")
    lines.append("1. Switch the provider source to `FabianSchurig/bitbucket`.")
    lines.append("2. Update provider authentication fields.")
    lines.append("3. Rename legacy resources/data sources to the generated equivalents below.")
    lines.append(
        "4. Rename common path inputs like `owner` → `workspace` and "
        "`repository` → `repo_slug`."
    )
    lines.append(
        "5. Review objects that split into multiple generated resources, especially "
        "repositories and variables."
    )
    lines.append("")
    lines.append("## Provider block changes")
    lines.append("")
    lines.append("### Example")
    lines.append("")
    lines.append("```hcl")
    lines.append("terraform {")
    lines.append("  required_providers {")
    lines.append("    bitbucket = {")
    lines.append('      source = "FabianSchurig/bitbucket"')
    lines.append("    }")
    lines.append("  }")
    lines.append("}")
    lines.append("")
    lines.append('provider "bitbucket" {')
    lines.append(
        "  username = var.bitbucket_username # optional for workspace/repo access tokens"
    )
    lines.append("  token    = var.bitbucket_token")
    lines.append("}")
    lines.append("```")
    lines.append("")
    lines.append("### Provider-level renames and removals")
    lines.append("")
    lines.append("| Legacy | New | Notes |")
    lines.append("|---|---|---|")
    for old, new, note in COMMON_RENAMES:
        lines.append(f"| {old} | {new} | {note} |")
    lines.append("")
    lines.append("## Coverage summary")
    lines.append("")
    lines.append(f"- Matched legacy resources: **{len(res_matched)} / {len(legacy['resource'])}**")
    lines.append(f"- Legacy-only resources: **{len(res_legacy_only)}**")
    lines.append(f"- New-only resources: **{len(res_new_only)}**")
    lines.append(
        f"- Matched legacy data sources: **{len(ds_matched)} / {len(legacy['data-source'])}**"
    )
    lines.append(f"- Legacy-only data sources: **{len(ds_legacy_only)}**")
    lines.append(f"- New-only data sources: **{len(ds_new_only)}**")
    lines.append("")
    lines.append("## Quick rename table for matched resources")
    lines.append("")
    lines.append("| Legacy resource | New resource(s) |")
    lines.append("|---|---|")
    for name in res_matched:
        targets = ", ".join(f"`{target}`" for target in CURRENT_ALIAS[("resource", name)])
        lines.append(f"| `{name}` | {targets} |")
    lines.append("")
    lines.append("## Quick rename table for matched data sources")
    lines.append("")
    lines.append("| Legacy data source | New data source(s) |")
    lines.append("|---|---|")
    for name in ds_matched:
        targets = ", ".join(
            f"`{target}`" for target in CURRENT_ALIAS[("data-source", name)]
        )
        lines.append(f"| `{name}` | {targets} |")
    lines.append("")

    for section_title, kind, names in [
        ("Matched legacy resources", "resource", res_matched),
        ("Legacy-only resources", "resource", res_legacy_only),
        ("Matched legacy data sources", "data-source", ds_matched),
        ("Legacy-only data sources", "data-source", ds_legacy_only),
    ]:
        lines.append(f"## {section_title}")
        lines.append("")
        for name in names:
            legacy_doc = parse_legacy_doc(kind, name)
            current_docs = mapped_current_objects(kind, name, current)
            legacy_doc.endpoints = parse_legacy_endpoints(kind, name, current_docs)
            lines.append(f"### `{name}`")
            lines.append("")
            if current_docs:
                lines.append(
                    "- New equivalent(s): "
                    + ", ".join(f"`{doc.name}`" for doc in current_docs)
                )
            else:
                lines.append("- New equivalent(s): none")
            lines.append(
                f"- Legacy inputs: {format_params(legacy_doc.inputs_required, legacy_doc.inputs_optional)}"
            )
            lines.append(f"- Legacy endpoints: {format_endpoints(legacy_doc.endpoints)}")
            if current_docs:
                lines.append(
                    "- New inputs: "
                    + format_params(
                        dedupe(sum((doc.inputs_required for doc in current_docs), [])),
                        dedupe(sum((doc.inputs_optional for doc in current_docs), [])),
                    )
                )
                new_endpoints = dedupe(sum((doc.endpoints for doc in current_docs), []))
                lines.append(f"- New operations: {format_endpoints(new_endpoints)}")
                lines.append(f"- Diff summary: {diff_summary(legacy_doc, current_docs)}")
            note = OBJECT_NOTES.get((kind, name))
            if note:
                lines.append(f"- Notes: {note}")
            lines.append("")

    lines.append("## New provider-only resources")
    lines.append("")
    for name in res_new_only:
        lines.append(f"- `{name}`")
    lines.append("")
    lines.append("## New provider-only data sources")
    lines.append("")
    for name in ds_new_only:
        lines.append(f"- `{name}`")
    lines.append("")
    lines.append("## Can this be automated?")
    lines.append("")
    lines.append(
        "Partly. A comparison script is practical today, but a fully automatic HCL "
        "rewrite is only safe for the straightforward cases."
    )
    lines.append("")
    lines.append("Good candidates for an automated rewrite later:")
    lines.append("")
    lines.append("- provider source replacement")
    lines.append("- provider auth field rename (`password` → `token`)")
    lines.append("- direct resource/data source renames where there is a 1:1 mapping")
    lines.append(
        "- path argument renames like `owner` → `workspace` and `repository` → `repo_slug`"
    )
    lines.append("")
    lines.append("Cases that still need manual review:")
    lines.append("")
    lines.append("- legacy objects that split into multiple generated resources")
    lines.append("- objects missing from one provider or the other")
    lines.append("- fields whose semantics changed even when the name looks similar")
    lines.append("- places where the generated provider expects `request_body` for uncommon fields")
    lines.append("")
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Generate a migration guide from the legacy Terraform provider."
    )
    parser.add_argument(
        "--repo-root",
        default=Path(__file__).resolve().parents[1],
        type=Path,
    )
    parser.add_argument(
        "--output",
        type=Path,
        help="Write markdown to this file instead of stdout.",
    )
    args = parser.parse_args()

    text = render(args.repo_root)
    if args.output:
        args.output.write_text(text)
    else:
        sys.stdout.write(text)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
