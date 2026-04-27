package spec

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestInternalBranchRestrictionsSchemaLoads validates the hand-authored
// internal-branch-restrictions schema. Because this schema is not produced by
// the OpenAPI sync pipeline (the internal API is not in Bitbucket's public
// spec), regressions can only be caught by an explicit test.
//
// The test asserts the full set of operations the generators rely on:
//   - operationIds and HTTP methods
//   - absolute URL paths (the dispatcher passes these through to resty
//     unchanged, which is what makes a different host work without runtime
//     changes)
//   - presence of request body fields for the PUT operations
//   - CLI parent-command metadata.
func TestInternalBranchRestrictionsSchemaLoads(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	schemaPath := filepath.Join(filepath.Dir(thisFile),
		"..", "..", "..", "schema", "internal-branch-restrictions-schema.yaml")

	schema, err := LoadSchema(schemaPath)
	if err != nil {
		t.Fatalf("LoadSchema(%s): %v", schemaPath, err)
	}

	name, use, short, _ := CommandMeta(schema.Info)
	if name != "ProjectBranchRestrictions" {
		t.Errorf("CommandMeta name = %q, want ProjectBranchRestrictions", name)
	}
	if use != "project-branch-restrictions" {
		t.Errorf("CommandMeta use = %q, want project-branch-restrictions", use)
	}
	if short == "" {
		t.Error("CommandMeta short is empty")
	}

	ops := BuildOperations(schema)

	type wantOp struct {
		method     string
		pathPrefix string
		hasBody    bool
		minBodyLen int
		mustHaveOp bool
	}
	want := map[string]wantOp{
		"getProjectBranchRestrictionsGroupedByBranch": {
			method:     "GET",
			pathPrefix: "https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/group-by-branch/",
			hasBody:    false,
			mustHaveOp: true,
		},
		"replaceProjectBranchRestrictionsByPattern": {
			method:     "PUT",
			pathPrefix: "https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-pattern/{pattern}",
			hasBody:    true,
			minBodyLen: 1,
			mustHaveOp: true,
		},
		"replaceProjectBranchRestrictionsByBranchType": {
			method:     "PUT",
			pathPrefix: "https://bitbucket.org/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-branch-type/{branch_type}",
			hasBody:    true,
			minBodyLen: 1,
			mustHaveOp: true,
		},
	}

	got := make(map[string]OperationDef, len(ops))
	for _, op := range ops {
		got[op.OperationID] = op
	}

	for id, w := range want {
		op, ok := got[id]
		if !ok {
			t.Errorf("BuildOperations missing operationId %q", id)
			continue
		}
		if op.Method != w.method {
			t.Errorf("op %s: Method = %q, want %q", id, op.Method, w.method)
		}
		if !strings.HasPrefix(op.Path, w.pathPrefix) && op.Path != w.pathPrefix {
			t.Errorf("op %s: Path = %q, want prefix %q", id, op.Path, w.pathPrefix)
		}
		if op.HasBody != w.hasBody {
			t.Errorf("op %s: HasBody = %v, want %v", id, op.HasBody, w.hasBody)
		}
		if w.hasBody && len(op.BodyFields) < w.minBodyLen {
			t.Errorf("op %s: len(BodyFields) = %d, want at least %d", id, len(op.BodyFields), w.minBodyLen)
		}
		// All three operations must require workspace and project_key path
		// params — they're the keys used by handlers.InferRepoContext and
		// the CLI required-flag logic.
		hasWorkspace := false
		hasProjectKey := false
		for _, p := range op.Params {
			if p.In == "path" && p.Name == "workspace" {
				hasWorkspace = true
			}
			if p.In == "path" && p.Name == "project_key" {
				hasProjectKey = true
			}
		}
		if !hasWorkspace {
			t.Errorf("op %s: missing required path param 'workspace'", id)
		}
		if !hasProjectKey {
			t.Errorf("op %s: missing required path param 'project_key'", id)
		}
	}
}
