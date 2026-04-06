import io
import tempfile
import unittest
from pathlib import Path
from unittest import mock
from urllib.error import HTTPError, URLError

import scripts.gen_migration as gen_migration


class DummyResponse:
    def __init__(self, payload: bytes):
        self.payload = payload

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc, tb):
        return False

    def read(self):
        return self.payload


class GenMigrationTests(unittest.TestCase):
    def test_fetch_returns_text(self):
        with mock.patch.object(
            gen_migration.urllib.request,
            "urlopen",
            return_value=DummyResponse(b"hello"),
        ) as urlopen:
            result = gen_migration.fetch("https://example.com/test")

        self.assertEqual(result, "hello")
        _, kwargs = urlopen.call_args
        self.assertEqual(kwargs["timeout"], 20)
        self.assertIsNotNone(kwargs["context"])

    def test_fetch_required_raises_runtime_error(self):
        with mock.patch.object(
            gen_migration.urllib.request,
            "urlopen",
            side_effect=URLError("boom"),
        ):
            with self.assertRaisesRegex(RuntimeError, "failed to fetch"):
                gen_migration.fetch("https://example.com/test")

    def test_fetch_optional_returns_none_and_warns(self):
        error = HTTPError("https://example.com/test", 404, "missing", None, None)
        with mock.patch.object(
            gen_migration.urllib.request,
            "urlopen",
            side_effect=error,
        ):
            stderr = io.StringIO()
            with mock.patch.object(gen_migration.sys, "stderr", stderr):
                result = gen_migration.fetch("https://example.com/test", required=False)

        self.assertIsNone(result)
        self.assertIn("warning: skipping unavailable URL", stderr.getvalue())

    def test_parse_current_doc_rejects_missing_heading(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "broken.md"
            path.write_text("no heading here\n")

            with self.assertRaisesRegex(ValueError, "missing Terraform object heading"):
                gen_migration.parse_current_doc(path, "resource")

    def test_parse_current_doc_rejects_heading_without_terraform_object_name(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "broken.md"
            path.write_text("# not-a-terraform-object\n")

            with self.assertRaisesRegex(ValueError, "missing Terraform object name"):
                gen_migration.parse_current_doc(path, "resource")

    def test_parse_bullets_filters_nested_current_bullets(self):
        section = "- `top_level`\n  - `nested`\n- `another_top_level`\n"

        parsed = gen_migration.parse_bullets(section, current=True)

        self.assertEqual(parsed, ["top_level", "another_top_level"])

    def test_parse_current_doc_ignores_nested_schema_bullets(self):
        markdown = """# bitbucket_example (Resource)

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `request_body` (String) Body.
  Nested schema:
  - `nested_only` (String) Nested field.

### Read-Only
- `api_response` (String) Response.
"""
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "example.md"
            path.write_text(markdown)
            doc = gen_migration.parse_current_doc(path, "resource")

        self.assertEqual(doc.inputs_required, ["workspace"])
        self.assertEqual(doc.inputs_optional, ["request_body"])
        self.assertNotIn("nested_only", doc.inputs_optional)
        self.assertEqual(doc.read_only, ["api_response"])

    def test_parse_legacy_doc_missing_file_returns_empty_doc(self):
        with mock.patch.object(gen_migration, "fetch", return_value=None):
            doc = gen_migration.parse_legacy_doc("resource", "bitbucket_missing")

        self.assertEqual(doc.kind, "resource")
        self.assertEqual(doc.name, "bitbucket_missing")
        self.assertEqual(doc.inputs, [])

    def test_parse_legacy_endpoints_missing_source_falls_back_to_current(self):
        current = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_branch_restrictions",
            endpoints=["Create POST /repositories/{workspace}/{repo_slug}/branch-restrictions"],
        )

        with mock.patch.object(gen_migration, "fetch", return_value=None):
            endpoints = gen_migration.parse_legacy_endpoints(
                "resource",
                "bitbucket_branch_restriction",
                [current],
            )

        self.assertEqual(endpoints, current.endpoints)

    def test_parse_legacy_endpoints_replaces_param_placeholders_with_mapped_endpoints(self):
        source = """
client.Put(fmt.Sprintf("2.0/repositories/%s/%s/override-settings", workspace, repoSlug))
client.Get(fmt.Sprintf("2.0/repositories/%s/%s/override-settings", workspace, repoSlug))
"""
        current = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_repo_settings",
            endpoints=[
                "Read GET /repositories/{workspace}/{repo_slug}/override-settings",
                "Update PUT /repositories/{workspace}/{repo_slug}/override-settings",
            ],
        )

        with mock.patch.object(gen_migration, "fetch", return_value=source):
            endpoints = gen_migration.parse_legacy_endpoints(
                "resource",
                "bitbucket_repository",
                [current],
            )

        self.assertEqual(
            endpoints,
            [
                "Read GET /repositories/{workspace}/{repo_slug}/override-settings",
                "Update PUT /repositories/{workspace}/{repo_slug}/override-settings",
            ],
        )

    def test_diff_summary_reports_renames_dropped_and_added_fields(self):
        legacy = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_repository",
            inputs_required=["owner", "repository"],
            inputs_optional=["legacy_only"],
        )
        current = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_repos",
            inputs_required=["workspace", "repo_slug"],
            inputs_optional=["request_body"],
        )

        summary = gen_migration.diff_summary(legacy, [current])

        self.assertIn("`owner` → `workspace`", summary)
        self.assertIn("`repository` → `repo_slug`", summary)
        self.assertIn("`legacy_only`", summary)
        self.assertIn("`request_body`", summary)

    def test_build_legacy_hcl_renders_required_and_optional_fields(self):
        legacy = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_branch_restriction",
            inputs_required=["repository", "kind", "owner"],
            inputs_optional=["pattern"],
        )

        lines = gen_migration.build_legacy_hcl(legacy)
        rendered = "\n".join(lines)

        self.assertIn('resource "bitbucket_branch_restriction" "legacy" {', rendered)
        self.assertLess(rendered.index('kind = "push"'), rendered.index('owner = "my-workspace"'))
        self.assertLess(rendered.index('owner = "my-workspace"'), rendered.index('repository = "my-repo"'))
        self.assertIn('owner = "my-workspace"', rendered)
        self.assertIn('repository = "my-repo"', rendered)
        self.assertIn('# pattern = "main"  # optional', rendered)

    def test_build_current_hcl_shows_renamed_and_legacy_only_fields(self):
        legacy = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_default_reviewers",
            inputs_required=["owner", "repository", "reviewers"],
        )
        current = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_default_reviewers",
            inputs_required=["workspace", "repo_slug", "target_username"],
        )

        lines = gen_migration.build_current_hcl(current, legacy)
        rendered = "\n".join(lines)

        self.assertIn('resource "bitbucket_default_reviewers" "migrated" {', rendered)
        self.assertIn('repo_slug = "my-repo"', rendered)
        self.assertIn('target_username = "example-user"', rendered)
        self.assertIn('workspace = "my-workspace"', rendered)
        self.assertLess(rendered.index('repo_slug = "my-repo"'), rendered.index('target_username = "example-user"'))
        self.assertLess(rendered.index('target_username = "example-user"'), rendered.index('workspace = "my-workspace"'))
        self.assertIn('# reviewers = ["example-user"]  # legacy-only', rendered)

    def test_render_uses_relative_docs_path(self):
        current = {
            ("resource", "bitbucket_branch_restrictions"): gen_migration.DocObject(
                kind="resource",
                name="bitbucket_branch_restrictions",
                inputs_required=["workspace"],
                endpoints=["Create POST /repositories/{workspace}/{repo_slug}/branch-restrictions"],
            ),
            ("data-source", "bitbucket_current_user"): gen_migration.DocObject(
                kind="data-source",
                name="bitbucket_current_user",
                endpoints=["Read GET /user"],
            ),
        }

        def fake_parse_legacy_doc(kind, name):
            return gen_migration.DocObject(kind=kind, name=name)

        def fake_parse_legacy_endpoints(kind, name, mapped_current):
            return ["GET /legacy"]

        with mock.patch.object(gen_migration, "current_objects", return_value=current), mock.patch.object(
            gen_migration,
            "legacy_names",
            return_value={
                "resource": ["bitbucket_branch_restriction"],
                "data-source": ["bitbucket_current_user"],
            },
        ), mock.patch.object(
            gen_migration,
            "parse_legacy_doc",
            side_effect=fake_parse_legacy_doc,
        ), mock.patch.object(
            gen_migration,
            "parse_legacy_endpoints",
            side_effect=fake_parse_legacy_endpoints,
        ):
            rendered = gen_migration.render(Path("/repo/root"))

        self.assertIn("- current docs from `./docs/`", rendered)
        self.assertIn("best-effort migration baseline", rendered)
        self.assertIn("`bitbucket_branch_restriction`", rendered)
        self.assertIn("`bitbucket_current_user`", rendered)
        self.assertIn("#### Legacy HCL", rendered)
        self.assertIn("#### New HCL", rendered)
        self.assertIn('resource "bitbucket_branch_restrictions" "migrated" {', rendered)


if __name__ == "__main__":
    unittest.main()
