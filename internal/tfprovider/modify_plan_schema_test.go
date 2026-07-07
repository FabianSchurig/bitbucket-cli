package tfprovider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestModifyPlan_RebuildsSchemaAttrsWhenCacheEmpty is a regression test for
// https://github.com/FabianSchurig/bitbucket-cli/issues/111.
//
// terraform-plugin-framework does not guarantee that Schema() runs on the
// same resource instance it uses for PlanResourceChange, so r.schemaAttrs
// can be empty when ModifyPlan fires. Previously the plan walk then
// iterated zero attributes, vacuously concluded "all equal", and replaced
// the plan with prior state — nulling a configured request_body and
// silently dropping body-field updates. ModifyPlan must rebuild the
// attribute set on demand so a real change is detected.
func TestModifyPlan_RebuildsSchemaAttrsWhenCacheEmpty(t *testing.T) {
	ctx := context.Background()

	// Build the schema once to derive the object type, but hand ModifyPlan
	// a *fresh* resource whose schemaAttrs cache is empty — exactly the
	// state the framework leaves the PlanResourceChange instance in.
	var sresp resource.SchemaResponse
	(&GenericResource{group: ReposResourceGroup}).Schema(ctx, resource.SchemaRequest{}, &sresp)
	objType := sresp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	stateVals := map[string]tftypes.Value{}
	cfgVals := map[string]tftypes.Value{}
	for name, ty := range objType.AttributeTypes {
		stateVals[name] = tftypes.NewValue(ty, nil)
		cfgVals[name] = tftypes.NewValue(ty, nil)
	}
	// request_body: null in prior state, set in config (the issue scenario).
	stateVals["request_body"] = tftypes.NewValue(tftypes.String, nil)
	cfgVals["request_body"] = tftypes.NewValue(tftypes.String, `{"description":"X"}`)

	state := tftypes.NewValue(objType, stateVals)
	cfg := tftypes.NewValue(objType, cfgVals)

	r := &GenericResource{group: ReposResourceGroup} // empty schemaAttrs cache
	req := resource.ModifyPlanRequest{
		State:  tfsdk.State{Schema: sresp.Schema, Raw: state},
		Config: tfsdk.Config{Schema: sresp.Schema, Raw: cfg},
		Plan:   tfsdk.Plan{Schema: sresp.Schema, Raw: cfg},
	}
	resp := resource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: sresp.Schema, Raw: cfg}}
	r.ModifyPlan(ctx, req, &resp)

	planMap := map[string]tftypes.Value{}
	if err := resp.Plan.Raw.As(&planMap); err != nil {
		t.Fatalf("decode plan: %v", err)
	}
	if planMap["request_body"].IsNull() {
		t.Fatal("request_body was nulled in the plan; ModifyPlan substituted prior state despite a real change")
	}
}

// TestModifyPlan_UnsetComputedAttrPreservesState is a regression test for the
// branching-model perpetual no-op diff. When an Optional+Computed attribute is
// left unset (null) in config but populated in prior state, it must not be
// treated as a change. Otherwise ModifyPlan bails, sibling Computed attributes
// get promoted to Unknown, and Terraform reports a spurious in-place update on
// every plan even though nothing changed.
func TestModifyPlan_UnsetComputedAttrPreservesState(t *testing.T) {
	ctx := context.Background()

	var sresp resource.SchemaResponse
	(&GenericResource{group: ReposResourceGroup}).Schema(ctx, resource.SchemaRequest{}, &sresp)
	objType := sresp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	stateVals := map[string]tftypes.Value{}
	cfgVals := map[string]tftypes.Value{}
	for name, ty := range objType.AttributeTypes {
		stateVals[name] = tftypes.NewValue(ty, nil)
		cfgVals[name] = tftypes.NewValue(ty, nil)
	}
	// description: managed and unchanged (same in config and state).
	stateVals["description"] = tftypes.NewValue(tftypes.String, "foo")
	cfgVals["description"] = tftypes.NewValue(tftypes.String, "foo")
	// language: Optional+Computed, populated in state, left unset in config.
	stateVals["language"] = tftypes.NewValue(tftypes.String, "go")
	cfgVals["language"] = tftypes.NewValue(tftypes.String, nil)

	state := tftypes.NewValue(objType, stateVals)
	cfg := tftypes.NewValue(objType, cfgVals)

	r := &GenericResource{group: ReposResourceGroup}
	req := resource.ModifyPlanRequest{
		State:  tfsdk.State{Schema: sresp.Schema, Raw: state},
		Config: tfsdk.Config{Schema: sresp.Schema, Raw: cfg},
		Plan:   tfsdk.Plan{Schema: sresp.Schema, Raw: cfg},
	}
	resp := resource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: sresp.Schema, Raw: cfg}}
	r.ModifyPlan(ctx, req, &resp)

	planMap := map[string]tftypes.Value{}
	if err := resp.Plan.Raw.As(&planMap); err != nil {
		t.Fatalf("decode plan: %v", err)
	}
	var lang *string
	if err := planMap["language"].As(&lang); err != nil {
		t.Fatalf("decode language: %v", err)
	}
	if lang == nil || *lang != "go" {
		t.Fatalf("expected language preserved from prior state (%q), got %v", "go", lang)
	}
}
