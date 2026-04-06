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

    def test_get_legacy_doc_url(self):
        self.assertEqual(
            gen_migration.get_legacy_doc_url("resource", "bitbucket_repository"),
            "https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/repository.md",
        )

    def test_get_current_doc_url(self):
        self.assertEqual(
            gen_migration.get_current_doc_url("data-source", "bitbucket_current_user"),
            "https://github.com/FabianSchurig/bitbucket-cli/blob/main/docs/data-sources/current-user.md",
        )

    def test_format_doc_link_handles_missing_legacy_docs(self):
        doc = gen_migration.DocObject(
            kind="resource",
            name="bitbucket_missing",
            doc_link="https://example.com/missing",
            doc_available=False,
        )

        self.assertEqual(
            gen_migration.format_doc_link(doc),
            "`bitbucket_missing` (legacy doc not available)",
        )

    def test_render_uses_main_branch_doc_links(self):
        current = {
            ("resource", "bitbucket_branch_restrictions"): gen_migration.DocObject(
                kind="resource",
                name="bitbucket_branch_restrictions",
                doc_link="https://github.com/FabianSchurig/bitbucket-cli/blob/main/docs/resources/branch-restrictions.md",
            ),
            ("data-source", "bitbucket_current_user"): gen_migration.DocObject(
                kind="data-source",
                name="bitbucket_current_user",
                doc_link="https://github.com/FabianSchurig/bitbucket-cli/blob/main/docs/data-sources/current-user.md",
            ),
        }

        def fake_parse_legacy_doc(kind, name):
            return gen_migration.DocObject(
                kind=kind,
                name=name,
                doc_link=f"https://legacy.example/{kind}/{name}.md",
            )

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
        ):
            rendered = gen_migration.render(Path("/repo/root"))

        self.assertIn("- current docs from `./docs/`", rendered)
        self.assertIn("It intentionally avoids generated field-by-field or HCL diffs", rendered)
        self.assertIn("`bitbucket_branch_restriction`", rendered)
        self.assertIn("`bitbucket_current_user`", rendered)
        self.assertIn("Legacy docs:", rendered)
        self.assertIn("New docs:", rendered)
        self.assertIn("- Legacy docs: [`bitbucket_branch_restriction`](https://legacy.example/resource/bitbucket_branch_restriction.md)", rendered)
        self.assertIn(
            "- New docs: [`bitbucket_branch_restrictions`](https://github.com/FabianSchurig/bitbucket-cli/blob/main/docs/resources/branch-restrictions.md)",
            rendered,
        )
        self.assertNotIn("#### Legacy HCL", rendered)
        self.assertNotIn("- Diff summary:", rendered)


if __name__ == "__main__":
    unittest.main()
