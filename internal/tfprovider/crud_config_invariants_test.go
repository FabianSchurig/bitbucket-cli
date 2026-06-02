package tfprovider_test

import (
	"sort"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

// These invariant tests are the guardrails that would have caught issue #100,
// where `bitbucket_pipeline_config` was a writable Terraform resource (it
// declared Read+Update) but had no Create mapping, so `terraform apply` failed
// at runtime with "Create not supported". The existing mock tests only asserted
// hand-maintained expectations that mirrored CRUDConfig, so a *missing* slot was
// invisible. These tests assert structural invariants about CRUDConfig itself,
// independently of any per-group expectation table, so the whole class of bug is
// caught at `go test` time instead of during a real `terraform apply`.

// allCRUDConfigOperationIDs collects every operation ID referenced by any slot
// of every CRUDConfig entry.
func allCRUDConfigOperationIDs() map[string][]string {
	refs := map[string][]string{}
	for name, m := range tfprovider.CRUDConfig {
		for _, id := range []string{m.Create, m.Read, m.Update, m.Delete, m.List} {
			if id != "" {
				refs[name] = append(refs[name], id)
			}
		}
	}
	return refs
}

// realOperationIDs returns the set of operation IDs that actually exist across
// all registered resource groups' generated operations.
func realOperationIDs() map[string]bool {
	valid := map[string]bool{}
	for _, g := range tfprovider.RegisteredGroups() {
		for _, op := range g.AllOps {
			valid[op.OperationID] = true
		}
	}
	return valid
}

// TestCRUDConfig_AllOperationIDsResolve ensures every operation ID referenced in
// CRUDConfig resolves to a real generated operation. MapCRUDOps silently drops
// unknown IDs (the lookup returns nil), so a typo or a renamed-away operation in
// the OpenAPI spec would otherwise vanish without error — producing exactly the
// "Create not supported"-style runtime failures from #100. This test makes such
// dead references fail loudly at build/test time.
func TestCRUDConfig_AllOperationIDsResolve(t *testing.T) {
	valid := realOperationIDs()
	for name, ids := range allCRUDConfigOperationIDs() {
		for _, id := range ids {
			if !valid[id] {
				t.Errorf("CRUDConfig[%q] references operation %q which does not exist in any "+
					"registered resource group; it would be silently dropped by MapCRUDOps. "+
					"Fix the operation ID or the CRUDConfig entry.", name, id)
			}
		}
	}
}

// writableNotCreatableAllowlist lists resource groups that are intentionally
// writable (they declare Update and/or Delete) yet deliberately have no Create
// mapping, because the underlying Bitbucket API offers no way to create the
// entity:
//
//   - addon:                   apps are installed out-of-band, not created via a PUT/POST here.
//   - pipeline-caches:         caches are produced as a side effect of pipeline runs; the API
//     only supports list/read-URI/delete.
//   - project-branch-restrictions: the "grouped by branch" read/list view; creation is handled
//     by the dedicated project-branch-restrictions-by-pattern /
//     -by-branch-type sub-resources, which do have Create.
//
// Adding a group here is a conscious, reviewed decision — which is the point:
// the test forces that decision instead of letting a missing Create slip
// through silently (issue #100).
var writableNotCreatableAllowlist = map[string]bool{
	"addon":                       true, // apps are installed out-of-band, not created via a PUT/POST here
	"pipeline-caches":             true, // caches are a side effect of pipeline runs; API only lists/reads-URI/deletes
	"project-branch-restrictions": true, // read/list grouping; creation handled by the -by-pattern/-by-branch-type sub-resources
}

// TestCRUDConfig_WritableResourcesAreCreatable enforces the invariant violated
// by #100: any resource group that is writable (declares Update or Delete) must
// also declare Create, so Terraform can actually create what it can manage.
// Groups that are genuinely not creatable must be added to
// writableNotCreatableAllowlist with a justification.
func TestCRUDConfig_WritableResourcesAreCreatable(t *testing.T) {
	var names []string
	for name := range tfprovider.CRUDConfig {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		m := tfprovider.CRUDConfig[name]
		writable := m.Update != "" || m.Delete != ""
		if !writable || m.Create != "" {
			continue
		}
		if writableNotCreatableAllowlist[name] {
			continue
		}
		t.Errorf("CRUDConfig[%q] is writable (update=%q delete=%q) but has no Create mapping: "+
			"`terraform apply` to create this resource would fail with \"Create not supported\" "+
			"(the #100 bug). Map Create to the upsert PUT operation, or add %q to "+
			"writableNotCreatableAllowlist with a justification.", name, m.Update, m.Delete, name)
	}
}
