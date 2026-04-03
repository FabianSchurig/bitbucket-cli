package tfprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
)

// Ensure the implementation satisfies the resource interface.
var _ resource.Resource = &GenericResource{}
var _ resource.ResourceWithConfigure = &GenericResource{}

// ─── Resource group metadata (shared with generators) ─────────────────────────

// OperationDef holds metadata for a single Bitbucket API operation.
type OperationDef struct {
	OperationID string
	Method      string
	Path        string
	Summary     string
	Description string
	Params      []ParamDef
	BodyFields  []BodyFieldDef
	HasBody     bool
	Paginated   bool
}

// ParamDef describes a single API parameter.
type ParamDef struct {
	Name     string
	In       string // "path" or "query"
	Type     string // "string", "integer", "boolean"
	Required bool
}

// BodyFieldDef describes a flattened request body field.
type BodyFieldDef struct {
	Path string // dot-separated path (e.g., "source.branch.name")
	Type string // "string", "integer", "boolean"
	Desc string
}

// CRUDOps maps CRUD operations to their OperationDef. All fields are optional —
// a resource only needs a Read operation at minimum.
type CRUDOps struct {
	Create *OperationDef
	Read   *OperationDef
	Update *OperationDef
	Delete *OperationDef
	List   *OperationDef
}

// ResourceGroup holds metadata for a Terraform resource generated from a
// Bitbucket API group (e.g., repositories, pull requests).
type ResourceGroup struct {
	TypeName    string // e.g., "bitbucket_repository"
	Description string
	Ops         CRUDOps
	AllOps      []OperationDef // all operations in the group
}

// ─── Generic resource implementation ──────────────────────────────────────────

// GenericResource implements a Terraform resource backed by Bitbucket API operations.
type GenericResource struct {
	group  ResourceGroup
	client *client.BBClient
}

// Metadata returns the resource type name.
func (r *GenericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.group.TypeName
}

// Schema builds the resource schema dynamically from the operation definitions.
func (r *GenericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Resource identifier.",
			Computed:    true,
		},
		"operation": schema.StringAttribute{
			Description: "The API operation to use. Set this to override the default CRUD operation selection.",
			Optional:    true,
		},
		"api_response": schema.StringAttribute{
			Description: "The raw JSON response from the Bitbucket API.",
			Computed:    true,
		},
	}

	// Build a set of params that are required in the primary Create or Read op.
	// Params from other ops (Update, Delete, List) are always Optional.
	primaryRequired := map[string]bool{}
	for _, op := range []*OperationDef{r.group.Ops.Create, r.group.Ops.Read} {
		if op == nil {
			continue
		}
		for _, p := range op.Params {
			if p.Required && p.In == "path" {
				primaryRequired[p.Name] = true
			}
		}
	}

	// Collect params from all CRUD ops.
	paramSeen := map[string]bool{}
	for _, op := range r.crudOps() {
		for _, p := range op.Params {
			if paramSeen[p.Name] {
				continue
			}
			paramSeen[p.Name] = true
			attrName := toSnakeCase(p.Name)
			isRequired := primaryRequired[p.Name]
			attrs[attrName] = schema.StringAttribute{
				Description: fmt.Sprintf("%s parameter (%s)", p.Name, p.In),
				Required:    isRequired,
				Optional:    !isRequired,
			}
		}
	}

	// Collect body fields from create/update ops.
	bodyFieldSeen := map[string]bool{}
	for _, op := range r.crudOps() {
		for _, bf := range op.BodyFields {
			key := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
			if bodyFieldSeen[key] {
				continue
			}
			bodyFieldSeen[key] = true
			desc := bf.Desc
			if desc == "" {
				desc = bf.Path
			}
			attrs[key] = schema.StringAttribute{
				Description: desc,
				Optional:    true,
			}
		}
	}

	resp.Schema = schema.Schema{
		Description: r.group.Description,
		Attributes:  attrs,
	}
}

// Configure receives the provider-configured client.
func (r *GenericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.BBClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.BBClient, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Create calls the create API operation and stores the result in state.
func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	op := r.group.Ops.Create
	if op == nil {
		resp.Diagnostics.AddError("Create not supported", fmt.Sprintf("Resource %s does not support create", r.group.TypeName))
		return
	}
	r.dispatch(ctx, op, &req.Plan, &resp.State, &resp.Diagnostics)
}

// Read calls the read API operation and refreshes state.
func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	op := r.group.Ops.Read
	if op == nil {
		// If no read operation, preserve existing state.
		return
	}
	r.dispatch(ctx, op, &req.State, &resp.State, &resp.Diagnostics)
}

// Update calls the update API operation and updates state.
func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	op := r.group.Ops.Update
	if op == nil {
		resp.Diagnostics.AddError("Update not supported", fmt.Sprintf("Resource %s does not support update", r.group.TypeName))
		return
	}
	r.dispatch(ctx, op, &req.Plan, &resp.State, &resp.Diagnostics)
}

// Delete calls the delete API operation.
func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	op := r.group.Ops.Delete
	if op == nil {
		// If no delete, just remove from state.
		return
	}

	plan := &req.State
	pathParams, queryParams := r.extractParams(ctx, op, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := handlers.DispatchRaw(ctx, r.client, handlers.Request{
		Method:      op.Method,
		URLTemplate: op.Path,
		PathParams:  pathParams,
		QueryParams: queryParams,
		All:         false,
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Delete failed: %v", err))
	}
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// dispatch executes an API operation, reading params from source and writing results to target.
func (r *GenericResource) dispatch(ctx context.Context, op *OperationDef, source stateAccessor, target stateAccessor, diags *diag.Diagnostics) {
	pathParams, queryParams := r.extractParams(ctx, op, source, diags)
	if diags.HasError() {
		return
	}

	body := r.buildBody(ctx, op, source, diags)
	if diags.HasError() {
		return
	}

	result, err := handlers.DispatchRaw(ctx, r.client, handlers.Request{
		Method:      op.Method,
		URLTemplate: op.Path,
		PathParams:  pathParams,
		QueryParams: queryParams,
		Body:        body,
		All:         op.Paginated,
	})
	if err != nil {
		diags.AddError("API Error", fmt.Sprintf("Operation %s failed: %v", op.OperationID, err))
		return
	}

	// Store the API response as JSON.
	if result != nil {
		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			diags.AddError("Response Error", fmt.Sprintf("Failed to marshal response: %v", err))
			return
		}
		diags.Append(target.SetAttribute(ctx, attrPath("api_response"), types.StringValue(string(jsonBytes)))...)

		// Try to extract an ID from the response.
		if m, ok := result.(map[string]any); ok {
			if id := extractID(m); id != "" {
				diags.Append(target.SetAttribute(ctx, attrPath("id"), types.StringValue(id))...)
			}
		}
	} else {
		diags.Append(target.SetAttribute(ctx, attrPath("api_response"), types.StringValue(""))...)
	}

	// Copy source attributes to target for params and body fields.
	r.copyAttributes(ctx, op, source, target, diags)
}

// extractParams reads path and query params from the source state/plan.
func (r *GenericResource) extractParams(ctx context.Context, op *OperationDef, source stateAccessor, diags *diag.Diagnostics) (map[string]string, map[string]string) {
	pathParams := map[string]string{}
	queryParams := map[string]string{}

	for _, p := range op.Params {
		attrName := toSnakeCase(p.Name)
		var val types.String
		d := source.GetAttribute(ctx, attrPath(attrName), &val)
		diags.Append(d...)
		if d.HasError() {
			continue
		}
		if val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
			if p.Required {
				diags.AddError("Missing Required Parameter",
					fmt.Sprintf("Parameter %q is required for operation %s", p.Name, op.OperationID))
			}
			continue
		}
		switch p.In {
		case "path":
			pathParams[p.Name] = val.ValueString()
		case "query":
			queryParams[p.Name] = val.ValueString()
		}
	}
	return pathParams, queryParams
}

// buildBody constructs the JSON request body from plan/state attributes.
func (r *GenericResource) buildBody(ctx context.Context, op *OperationDef, source stateAccessor, diags *diag.Diagnostics) string {
	if !op.HasBody || len(op.BodyFields) == 0 {
		return ""
	}

	bodyObj := map[string]any{}
	for _, bf := range op.BodyFields {
		attrName := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
		var val types.String
		d := source.GetAttribute(ctx, attrPath(attrName), &val)
		diags.Append(d...)
		if d.HasError() || val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
			continue
		}
		handlers.SetNested(bodyObj, bf.Path, val.ValueString())
	}

	if len(bodyObj) == 0 {
		return ""
	}
	b, err := json.Marshal(bodyObj)
	if err != nil {
		diags.AddError("Body Error", fmt.Sprintf("Failed to marshal request body: %v", err))
		return ""
	}
	return string(b)
}

// copyAttributes copies all param and body field values from source to target state.
// It copies from ALL CRUD operations so that attributes defined by one operation
// (e.g., Read's path params) are preserved when another operation (e.g., Create) runs.
func (r *GenericResource) copyAttributes(ctx context.Context, _ *OperationDef, source, target stateAccessor, diags *diag.Diagnostics) {
	seen := map[string]bool{}
	for _, op := range r.crudOps() {
		for _, p := range op.Params {
			attrName := toSnakeCase(p.Name)
			if seen[attrName] {
				continue
			}
			seen[attrName] = true
			var val types.String
			d := source.GetAttribute(ctx, attrPath(attrName), &val)
			diags.Append(d...)
			if !d.HasError() && !val.IsNull() && !val.IsUnknown() {
				diags.Append(target.SetAttribute(ctx, attrPath(attrName), val)...)
			}
		}
		for _, bf := range op.BodyFields {
			attrName := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
			if seen[attrName] {
				continue
			}
			seen[attrName] = true
			var val types.String
			d := source.GetAttribute(ctx, attrPath(attrName), &val)
			diags.Append(d...)
			if !d.HasError() && !val.IsNull() && !val.IsUnknown() {
				diags.Append(target.SetAttribute(ctx, attrPath(attrName), val)...)
			}
		}
	}
}

// crudOps returns all non-nil CRUD operations.
func (r *GenericResource) crudOps() []*OperationDef {
	ops := []*OperationDef{r.group.Ops.Create, r.group.Ops.Read, r.group.Ops.Update, r.group.Ops.Delete, r.group.Ops.List}
	var result []*OperationDef
	for _, op := range ops {
		if op != nil {
			result = append(result, op)
		}
	}
	return result
}

// extractID tries to extract an identifier from an API response map.
func extractID(m map[string]any) string {
	// Try common ID fields.
	for _, key := range []string{"uuid", "id", "slug", "name"} {
		if v, ok := m[key]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
