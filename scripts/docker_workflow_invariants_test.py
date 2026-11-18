#!/usr/bin/env python3
"""Guard the publishing invariants of .github/workflows/docker.yml.

Two classes of regression are covered:

Release-permission regressions (observed in run 32467532833 for tag v0.20.0):
  1) sbom-action tried to attach release assets from a job with contents: read
     -> "Resource not accessible by integration"
  2) the failed SBOM step skipped scan-action entirely
  3) upload-sarif then ran with an empty sarif_file input

Consumer-visibility regressions: the image SBOM and CVE reports must keep
reaching image users, i.e. a signed SBOM attestation on the GHCR digest plus
SPDX/Grype files attached to the GitHub release.

Requires PyYAML. Runs with the other pipeline tests:
    python3 -m unittest discover -s scripts -p '*_test.py'
"""

from __future__ import annotations

import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
WORKFLOW = ROOT / ".github/workflows/docker.yml"

SBOM_STEP = "Generate SBOM"
ATTEST_STEP = "Attest SBOM"
SCAN_STEP = "Scan "
SARIF_STEP = "Upload CVE scan report"
CVE_REPORT_STEP = "Build CVE report"
CVE_ARTIFACT_STEP = "Upload CVE reports"
PUBLISH_JOB = "publish-image-security"

# scan-action is allowed to run even when an earlier step failed only through
# one of these guards; anything else silently drops the CVE scan.
SCAN_RUNS_AFTER_FAILURE = {
    "${{ !cancelled() }}",
    "!cancelled()",
    "${{ success() || failure() }}",
    "success() || failure()",
}


def as_text(value: object) -> str:
    """Normalize a YAML scalar so that `false` and "false" compare equal."""
    if isinstance(value, bool):
        return "true" if value else "false"
    return str(value)


def find_step(steps: list[dict], name_prefix: str) -> dict | None:
    for step in steps:
        if str(step.get("name", "")).startswith(name_prefix):
            return step
    return None


def find_steps(steps: list[dict], name_prefix: str) -> list[dict]:
    return [s for s in steps if str(s.get("name", "")).startswith(name_prefix)]


def check_release_permissions(docker: dict, failures: list[str]) -> None:
    steps = docker["steps"]
    perms = docker.get("permissions") or {}
    sbom = find_step(steps, SBOM_STEP)
    scan = find_step(steps, SCAN_STEP)
    sarif = find_step(steps, SARIF_STEP)

    if sbom is None or scan is None or sarif is None:
        failures.append("docker job is missing the SBOM, scan, or upload-sarif step")
        return

    upload_release_assets = as_text(sbom.get("with", {}).get("upload-release-assets", True))
    if perms.get("contents") == "read" and upload_release_assets != "false":
        failures.append(
            "sbom-action would attach release assets from a job with contents: read"
        )

    scan_if = str(scan.get("if", ""))
    if not sbom.get("continue-on-error") and scan_if not in SCAN_RUNS_AFTER_FAILURE:
        failures.append("a failing SBOM step would skip the CVE scan")

    sarif_if = str(sarif.get("if", ""))
    uses_scan_output = "steps.scan.outputs.sarif" in str(
        sarif.get("with", {}).get("sarif_file", "")
    )
    guarded = "steps.scan.outputs.sarif !=" in sarif_if or "steps.scan.outcome" in sarif_if
    if "always()" in sarif_if and uses_scan_output and not guarded:
        failures.append("upload-sarif would run with an empty sarif_file when the scan is skipped")


def check_consumer_visibility(docker: dict, workflow: dict, failures: list[str]) -> None:
    steps = docker["steps"]
    perms = docker.get("permissions") or {}
    sbom = find_step(steps, SBOM_STEP)
    attest = find_step(steps, ATTEST_STEP)

    if sbom is None or not sbom.get("with", {}).get("output-file"):
        failures.append("SBOM step must write output-file so it can be attested and published")

    if attest is None:
        failures.append("image SBOM is not attested to the GHCR digest")
    else:
        with_ = attest.get("with", {})
        if "steps.build.outputs.digest" not in str(with_.get("subject-digest", "")):
            failures.append("SBOM attestation is not bound to the pushed image digest")
        if not with_.get("sbom-path"):
            failures.append("SBOM attestation has no sbom-path")
        if perms.get("attestations") != "write":
            failures.append("docker job needs attestations: write to attest the SBOM")

    formats = {
        str(step.get("with", {}).get("output-format", ""))
        for step in find_steps(steps, CVE_REPORT_STEP)
    }
    missing = {"json", "table"} - formats
    if missing:
        failures.append(f"missing consumer CVE report format(s): {sorted(missing)}")

    if find_step(steps, CVE_ARTIFACT_STEP) is None:
        failures.append("CVE reports are never uploaded as workflow artifacts")

    publish = workflow["jobs"].get(PUBLISH_JOB)
    if publish is None:
        failures.append(f"missing {PUBLISH_JOB} job that attaches reports to the release")
        return

    publish_perms = publish.get("permissions") or {}
    if publish_perms.get("contents") != "write":
        failures.append(f"{PUBLISH_JOB} needs contents: write to upload release assets")
    if set(publish_perms) != {"contents"}:
        failures.append(f"{PUBLISH_JOB} should only request contents: write")
    if perms.get("contents") != "read":
        failures.append("docker build job must keep contents: read")

    if not any("gh release upload" in str(step.get("run", "")) for step in publish["steps"]):
        failures.append(f"{PUBLISH_JOB} does not upload anything to the release")


class DockerWorkflowInvariantsTest(unittest.TestCase):
    def test_publishing_invariants(self):
        workflow = yaml.safe_load(WORKFLOW.read_text(encoding="utf-8"))
        docker = workflow["jobs"]["docker"]
        failures: list[str] = []
        check_release_permissions(docker, failures)
        check_consumer_visibility(docker, workflow, failures)
        self.assertEqual(failures, [])


if __name__ == "__main__":
    unittest.main()
