import tempfile
import unittest
from pathlib import Path

import scripts.partition_spec as partition_spec


class InlineRequestBodyRefsTests(unittest.TestCase):
    def _spec(self) -> dict:
        return {
            "components": {
                "requestBodies": {
                    "project": {
                        "required": True,
                        "content": {
                            "application/json": {
                                "schema": {"$ref": "#/components/schemas/project"}
                            }
                        },
                    }
                },
                "schemas": {
                    "project": {
                        "type": "object",
                        "properties": {"key": {"type": "string"}},
                    }
                },
            }
        }

    def test_inlines_request_body_ref(self):
        spec = self._spec()
        paths = {
            "/workspaces/{workspace}/projects": {
                "post": {
                    "operationId": "createAProjectInAWorkspace",
                    "requestBody": {"$ref": "#/components/requestBodies/project"},
                }
            }
        }
        inlined = partition_spec.inline_request_body_refs(paths, spec)
        self.assertEqual(inlined, 1)
        rb = paths["/workspaces/{workspace}/projects"]["post"]["requestBody"]
        self.assertNotIn("$ref", rb)
        self.assertEqual(
            rb["content"]["application/json"]["schema"]["$ref"],
            "#/components/schemas/project",
        )

    def test_leaves_inline_body_untouched(self):
        spec = self._spec()
        inline_body = {
            "content": {"application/json": {"schema": {"$ref": "#/components/schemas/project"}}}
        }
        paths = {"/x": {"post": {"operationId": "op", "requestBody": inline_body}}}
        inlined = partition_spec.inline_request_body_refs(paths, spec)
        self.assertEqual(inlined, 0)
        self.assertEqual(paths["/x"]["post"]["requestBody"], inline_body)

    def test_skips_missing_component(self):
        spec = {"components": {"requestBodies": {}, "schemas": {}}}
        paths = {
            "/x": {
                "post": {
                    "operationId": "op",
                    "requestBody": {"$ref": "#/components/requestBodies/missing"},
                }
            }
        }
        inlined = partition_spec.inline_request_body_refs(paths, spec)
        self.assertEqual(inlined, 0)
        # Reference left as-is (unresolvable), never silently corrupted.
        self.assertEqual(
            paths["/x"]["post"]["requestBody"],
            {"$ref": "#/components/requestBodies/missing"},
        )

    def test_build_schema_inlines_and_copies_schema(self):
        # End-to-end: build_schema should inline the body and copy the nested
        # project schema so the output is self-contained (no dangling refs).
        spec = self._spec()
        spec["info"] = {"version": "2.0.0"}
        spec["paths"] = {
            "/workspaces/{workspace}/projects": {
                "post": {
                    "operationId": "createAProjectInAWorkspace",
                    "tags": ["Projects"],
                    "requestBody": {"$ref": "#/components/requestBodies/project"},
                }
            }
        }
        group = {
            "title": "T",
            "tags": {"Projects"},
            "paths": ["/workspaces/{workspace}/projects"],
            "cli_meta": {},
        }
        out = partition_spec.build_schema(spec, group)
        post = out["paths"]["/workspaces/{workspace}/projects"]["post"]
        self.assertNotIn("$ref", post["requestBody"])
        self.assertIn("project", out["components"]["schemas"])


class WriteSchemaTests(unittest.TestCase):
    def test_preserves_existing_schema_when_new_schema_has_no_paths(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            output_path = Path(temp_dir) / "schema.yaml"
            existing = "openapi: 3.0.0\npaths:\n  /existing:\n    get: {}\n"
            output_path.write_text(existing)

            partition_spec.write_schema(
                {"paths": {}, "components": {"schemas": {}}}, output_path
            )

            self.assertEqual(output_path.read_text(), existing)

    def test_writes_empty_schema_without_populated_existing_schema(self):
        for existing in (None, "openapi: 3.0.0\npaths: {}\n", "paths: ["):
            with self.subTest(existing=existing):
                with tempfile.TemporaryDirectory() as temp_dir:
                    output_path = Path(temp_dir) / "empty-schema.yaml"
                    if existing is not None:
                        output_path.write_text(existing)

                    partition_spec.write_schema(
                        {"paths": {}, "components": {"schemas": {"kept": {}}}},
                        output_path,
                    )

                    self.assertIn("kept: {}", output_path.read_text())

    def test_overwrites_existing_schema_when_new_schema_has_paths(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            output_path = Path(temp_dir) / "schema.yaml"
            output_path.write_text("openapi: 3.0.0\npaths:\n  /old:\n    get: {}\n")

            partition_spec.write_schema(
                {
                    "paths": {"/new": {"post": {}}},
                    "components": {"schemas": {}},
                },
                output_path,
            )

            written = output_path.read_text()
            self.assertIn("/new:", written)
            self.assertNotIn("/old:", written)


if __name__ == "__main__":
    unittest.main()
