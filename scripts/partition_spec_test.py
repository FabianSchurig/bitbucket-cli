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

    def test_writes_empty_schema_when_existing_schema_is_missing_or_unusable(self):
        for existing in (None, "paths: [", b"\xff"):
            with self.subTest(existing=existing):
                with tempfile.TemporaryDirectory() as temp_dir:
                    output_path = Path(temp_dir) / "empty-schema.yaml"
                    if existing is not None:
                        if isinstance(existing, bytes):
                            output_path.write_bytes(existing)
                        else:
                            output_path.write_text(existing)

                    partition_spec.write_schema(
                        {"paths": {}, "components": {"schemas": {"kept": {}}}},
                        output_path,
                    )

                    self.assertIn("kept: {}", output_path.read_text())

    def test_writes_empty_schema_when_existing_schema_has_no_paths(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            output_path = Path(temp_dir) / "empty-schema.yaml"
            output_path.write_text("openapi: 3.0.0\npaths: {}\n")

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


class RetainedOperationsTests(unittest.TestCase):
    """Operations Atlassian stops publishing must not vanish from the CLI."""

    def _existing(self) -> dict:
        return {
            "paths": {
                "/snippets": {
                    "get": {
                        "operationId": "listSnippets",
                        "responses": {
                            "200": {
                                "content": {
                                    "application/json": {
                                        "schema": {
                                            "$ref": "#/components/schemas/paginated_snippets"
                                        }
                                    }
                                }
                            }
                        },
                    },
                    "post": {"operationId": "createASnippet"},
                    "parameters": [],
                }
            },
            "components": {
                "schemas": {
                    "paginated_snippets": {
                        "type": "object",
                        "properties": {
                            "values": {"$ref": "#/components/schemas/snippet"}
                        },
                    },
                    "snippet": {"type": "object"},
                }
            },
        }

    def _fresh(self) -> dict:
        return {
            "paths": {"/snippets": {"post": {"operationId": "createASnippet"}}},
            "components": {"schemas": {}},
        }

    def test_retains_operation_dropped_from_live_spec(self):
        out = self._fresh()
        retained = partition_spec.merge_retained_operations(
            out, self._existing(), partition_spec.published_operations([out])
        )

        self.assertEqual(retained, ["listSnippets"])
        op = out["paths"]["/snippets"]["get"]
        self.assertEqual(op["operationId"], "listSnippets")
        self.assertTrue(op["x-bb-cli-retained"])
        # The operation is carried over verbatim: a transient spec omission
        # must not flip a working endpoint to deprecated (which would hide the
        # generated command from `bb-cli --help`).
        self.assertNotIn("deprecated", op)
        # Transitively referenced schemas travel with the retained operation so
        # the partitioned file stays self-contained.
        self.assertIn("paginated_snippets", out["components"]["schemas"])
        self.assertIn("snippet", out["components"]["schemas"])
        # Methods keep their canonical order ahead of path-level keys.
        self.assertEqual(
            list(out["paths"]["/snippets"]), ["get", "post", "parameters"]
        )

    def test_published_operations_are_not_retained(self):
        existing = self._existing()
        out = {
            "paths": {"/snippets": {"get": {"operationId": "listSnippets"}}},
            "components": {"schemas": {}},
        }
        retained = partition_spec.merge_retained_operations(
            out, existing, partition_spec.published_operations([out])
        )

        self.assertEqual(retained, ["createASnippet"])
        self.assertEqual(out["paths"]["/snippets"]["get"], {"operationId": "listSnippets"})

    def test_operation_moved_to_another_group_is_not_duplicated(self):
        # published_operations covers every group, so an endpoint that merely
        # moved between schema files is not resurrected in its old file.
        old_group = {"paths": {}, "components": {"schemas": {}}}
        new_group = {
            "paths": {"/snippets": {"get": {"operationId": "listSnippets"}}},
            "components": {"schemas": {}},
        }
        published = partition_spec.published_operations([old_group, new_group])

        retained = partition_spec.merge_retained_operations(
            old_group,
            {"paths": {"/snippets": {"get": {"operationId": "listSnippets"}}}},
            published,
        )

        self.assertEqual(retained, [])
        self.assertEqual(old_group["paths"], {})

    def test_retention_is_idempotent(self):
        first = self._fresh()
        partition_spec.merge_retained_operations(
            first, self._existing(), partition_spec.published_operations([first])
        )

        second = self._fresh()
        partition_spec.merge_retained_operations(
            second, first, partition_spec.published_operations([second])
        )

        self.assertEqual(first, second)

    def test_deleting_from_committed_schema_retires_operation(self):
        # The committed schema is the only baseline: removing an endpoint there
        # is how a maintainer retires it permanently.
        out = self._fresh()
        retained = partition_spec.merge_retained_operations(
            out,
            {"paths": {"/snippets": {"post": {"operationId": "createASnippet"}}}},
            partition_spec.published_operations([out]),
        )

        self.assertEqual(retained, [])
        self.assertNotIn("get", out["paths"]["/snippets"])

    def test_missing_or_unreadable_existing_schema_is_ignored(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            missing = Path(temp_dir) / "absent-schema.yaml"
            self.assertEqual(partition_spec.load_existing_schema(missing), {})

            broken = Path(temp_dir) / "broken-schema.yaml"
            broken.write_text("paths: [unbalanced\n")
            self.assertEqual(partition_spec.load_existing_schema(broken), {})


if __name__ == "__main__":
    unittest.main()
