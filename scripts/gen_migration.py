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

@dataclass
class DocObject:
    kind: str
    name: str
    title: str = ""
    doc_link: str = ""
    doc_available: bool = True
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
    """Fetch UTF-8 text from a URL with timeout and optional missing-file tolerance."""
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


def legacy_doc_url(kind: str, name: str) -> str:
    base = name.removeprefix("bitbucket_")
    return (
        "https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/"
        f"{CURRENT_KIND_PATH[kind]}/{base}.md"
    )


def current_doc_url(kind: str, name: str) -> str:
    return f"./docs/{CURRENT_KIND_PATH[kind]}/{name}.md"


def parse_current_doc(path: Path, kind: str) -> DocObject:
    """Parse one generated doc in ./docs/resources or ./docs/data-sources."""
    text = path.read_text()
    title_match = re.search(r"^#\s+(.+)$", text, re.M)
    if title_match is None:
        raise ValueError(f"missing Terraform object heading in {path}")
    title = title_match.group(1)
    name_match = re.search(r"(bitbucket_[^\s]+)", title)
    if name_match is None:
        raise ValueError(f"missing Terraform object name in heading for {path}")
    name = name_match.group(1)
    doc = DocObject(kind=kind, name=name, title=title, doc_link=current_doc_url(kind, name))
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
    """Parse one legacy markdown doc from the DrFaust92 provider, if it exists."""
    base = name.removeprefix("bitbucket_")
    doc_url = f"{LEGACY_BASE}/docs/{CURRENT_KIND_PATH[kind]}/{base}.md"
    link = legacy_doc_url(kind, name)
    text = fetch(doc_url, required=False)
    if text is None:
        return DocObject(kind=kind, name=name, title=name, doc_link=link, doc_available=False)
    title_match = re.search(r"^#\s+([^\n]+)", text, re.M)
    title = name
    if title_match:
        title = html.unescape(title_match.group(1).replace("\\_", "_"))
    doc = DocObject(kind=kind, name=name, title=title, doc_link=link)
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


def format_doc_link(doc: DocObject) -> str:
    if doc.doc_available:
        return f"[`{doc.name}`]({doc.doc_link})"
    return f"`{doc.name}` (legacy doc not available)"


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
        "It intentionally avoids generated field-by-field or HCL diffs. Nested fields, "
        "computed attributes, and generated doc structure can otherwise produce "
        "misleading migration advice."
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
            lines.append(f"### `{name}`")
            lines.append("")
            if current_docs:
                lines.append(
                    "- New equivalent(s): "
                    + ", ".join(f"`{doc.name}`" for doc in current_docs)
                )
            else:
                lines.append("- New equivalent(s): none")
            lines.append(f"- Legacy docs: {format_doc_link(legacy_doc)}")
            if current_docs:
                lines.append(
                    "- New docs: "
                    + ", ".join(format_doc_link(doc) for doc in current_docs)
                )
            note = OBJECT_NOTES.get((kind, name))
            if note:
                lines.append(f"- Notes: {note}")
            lines.append("")

    lines.append("## New provider-only resources")
    lines.append("")
    for name in res_new_only:
        lines.append(f"- {format_doc_link(current[('resource', name)])}")
    lines.append("")
    lines.append("## New provider-only data sources")
    lines.append("")
    for name in ds_new_only:
        lines.append(f"- {format_doc_link(current[('data-source', name)])}")
    lines.append("")
    lines.append("## Can this be automated?")
    lines.append("")
    lines.append(
        "Only partly. The rename tables are useful, but the actual migration still "
        "needs a human review against the authoritative docs."
    )
    lines.append("")
    lines.append("Good candidates for an automated rewrite later:")
    lines.append("")
    lines.append("- provider source replacement")
    lines.append("- provider auth field rename (`password` → `token`)")
    lines.append("- direct resource/data source renames where there is a 1:1 mapping")
    lines.append("")
    lines.append("Cases that still need manual review:")
    lines.append("")
    lines.append("- legacy objects that split into multiple generated resources")
    lines.append("- objects missing from one provider or the other")
    lines.append("- nested or computed fields")
    lines.append("- fields whose semantics changed even when the name looks similar")
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
