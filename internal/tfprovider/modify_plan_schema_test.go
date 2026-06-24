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
