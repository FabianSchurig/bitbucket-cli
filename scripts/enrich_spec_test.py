import tempfile
import unittest
from pathlib import Path

import scripts.enrich_spec as enrich_spec


class ApplyRequestBodyPatchesTests(unittest.TestCase):
    def _spec_with_component(self, post_op: dict) -> dict:
        return {
            "paths": {
                "/workspaces/{workspace}/projects": {"post": post_op},
            },
            "components": {
                "requestBodies": {
                    "project": {
                        "content": {
                            "application/json": {
                                "schema": {"$ref": "#/components/schemas/project"}
                            }
                        }
                    }
                }
            },
        }

    def test_injects_missing_body_from_component(self):
        spec = self._spec_with_component({"operationId": "createAProjectInAWorkspace"})
        applied = enrich_spec.apply_request_body_patches(spec)
        self.assertEqual(applied, 1)
        rb = spec["paths"]["/workspaces/{workspace}/projects"]["post"]["requestBody"]
        self.assertEqual(rb, {"$ref": "#/components/requestBodies/project"})

    def test_does_not_overwrite_existing_body(self):
        existing = {"content": {"application/json": {"schema": {"type": "object"}}}}
        spec = self._spec_with_component(
            {"operationId": "createAProjectInAWorkspace", "requestBody": existing}
        )
        applied = enrich_spec.apply_request_body_patches(spec)
        self.assertEqual(applied, 0)
        self.assertEqual(
            spec["paths"]["/workspaces/{workspace}/projects"]["post"]["requestBody"],
            existing,
        )

    def test_skips_when_referenced_component_absent(self):
        # No components/requestBodies/project — must not create a dangling ref.
        spec = {
            "paths": {
                "/workspaces/{workspace}/projects": {
                    "post": {"operationId": "createAProjectInAWorkspace"}
                }
            },
            "components": {"requestBodies": {}},
        }
        applied = enrich_spec.apply_request_body_patches(spec)
        self.assertEqual(applied, 0)
        self.assertNotIn(
            "requestBody",
            spec["paths"]["/workspaces/{workspace}/projects"]["post"],
        )

    def test_injects_inline_body_when_schema_ref_resolves(self):
        # Branching-model style: inline requestBody referencing a schema that
        # exists -> injected.
        path = "/repositories/{workspace}/{repo_slug}/branching-model/settings"
        spec = {
            "paths": {path: {"put": {"operationId": "updateTheBranchingModelConfigForARepository"}}},
            "components": {
                "schemas": {"branching_model_settings": {"type": "object"}},
                "requestBodies": {},
            },
        }
        applied = enrich_spec.apply_request_body_patches(spec)
        self.assertEqual(applied, 1)
        rb = spec["paths"][path]["put"]["requestBody"]
        self.assertEqual(
            rb["content"]["application/json"]["schema"]["$ref"],
            "#/components/schemas/branching_model_settings",
        )

    def test_skips_inline_body_when_schema_ref_missing(self):
        # Same operation but the referenced schema is absent -> skipped, no
        # dangling reference introduced.
        path = "/repositories/{workspace}/{repo_slug}/branching-model/settings"
        spec = {
            "paths": {path: {"put": {"operationId": "updateTheBranchingModelConfigForARepository"}}},
            "components": {"schemas": {}, "requestBodies": {}},
        }
        applied = enrich_spec.apply_request_body_patches(spec)
        self.assertEqual(applied, 0)
        self.assertNotIn("requestBody", spec["paths"][path]["put"])


class OperationIdLockTests(unittest.TestCase):
    """The lockfile is what keeps CLI command and MCP tool names stable."""

    def _spec(self) -> dict:
        return {
            "paths": {
                "/workspaces": {"get": {"summary": "List workspaces for user"}},
                "/user/workspaces": {
                    "get": {"summary": "List workspaces for the current user"}
                },
            }
        }

    def test_locked_ids_win_over_derived_ones(self):
        spec = self._spec()
        lock = {"get /user/workspaces": "getUserWorkspaces"}
        count, updated = enrich_spec.assign_operation_ids(spec, lock)

        self.assertEqual(count, 2)
        self.assertEqual(
            spec["paths"]["/user/workspaces"]["get"]["operationId"],
            "getUserWorkspaces",
        )
        self.assertEqual(updated["get /workspaces"], "listWorkspacesForUser")

    def test_vanished_operation_keeps_its_id_reserved(self):
        # Regression for the schema-sync failure: Atlassian dropped
        # GET /user/permissions/workspaces, and GET /user/workspaces inherited
        # its id, renaming a shipped CLI command. The reserved lock entry must
        # prevent that even though the operation is gone from the spec.
        spec = {
            "paths": {
                "/user/workspaces": {
                    "get": {"summary": "List workspaces for the current user"}
                }
            }
        }
        lock = {
            "get /user/permissions/workspaces": "listWorkspacesForTheCurrentUser",
        }
        _, updated = enrich_spec.assign_operation_ids(spec, lock)

        oid = spec["paths"]["/user/workspaces"]["get"]["operationId"]
        self.assertNotEqual(oid, "listWorkspacesForTheCurrentUser")
        self.assertEqual(oid, "getUserWorkspaces")
        self.assertEqual(
            updated["get /user/permissions/workspaces"],
            "listWorkspacesForTheCurrentUser",
        )

    def test_ids_are_independent_of_spec_ordering(self):
        forward = self._spec()
        reverse = {"paths": dict(reversed(list(self._spec()["paths"].items())))}

        _, lock_forward = enrich_spec.assign_operation_ids(forward, {})
        _, lock_reverse = enrich_spec.assign_operation_ids(reverse, {})

        self.assertEqual(lock_forward, lock_reverse)

    def test_colliding_summaries_fall_back_to_path_slug(self):
        spec = {
            "paths": {
                "/a/activity": {"get": {"summary": "List activity"}},
                "/b/activity": {"get": {"summary": "List activity"}},
            }
        }
        _, updated = enrich_spec.assign_operation_ids(spec, {})

        self.assertEqual(len(set(updated.values())), 2)
        self.assertIn("getBActivity", updated.values())

    def test_lock_roundtrip(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            lock_path = Path(temp_dir) / "operation-ids.json"
            self.assertEqual(enrich_spec.load_lock(lock_path), {})

            enrich_spec.write_lock({"b": "2", "a": "1"}, lock_path)
            self.assertEqual(enrich_spec.load_lock(lock_path), {"a": "1", "b": "2"})
            # Sorted on disk so the committed lockfile stays diff-stable.
            self.assertLess(
                lock_path.read_text().index('"a"'),
                lock_path.read_text().index('"b"'),
            )

    def test_malformed_lock_is_rejected(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            lock_path = Path(temp_dir) / "operation-ids.json"
            lock_path.write_text("[]")
            with self.assertRaises(ValueError):
                enrich_spec.load_lock(lock_path)


if __name__ == "__main__":
    unittest.main()
