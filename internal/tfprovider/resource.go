package tfprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
var _ resource.ResourceWithImportState = &GenericResource{}

// ─── Resource group metadata (shared with generators) ─────────────────────────

// OperationDef holds metadata for a single Bitbucket API operation.
type OperationDef struct {
	OperationID    string
	Method         string
	Path           string
	Summary        string
	Description    string
	Params         []ParamDef
	BodyFields     []BodyFieldDef
	ResponseFields []BodyFieldDef // Flattened fields from the response schema (computed)
	HasBody        bool
	Paginated      bool
	Scopes         []string // OAuth2 scopes from x-atlassian-oauth2-scopes
	DocURL         string   // Atlassian REST API documentation URL
}

// ParamDef describes a single API parameter.
type ParamDef struct {
	Name     string
	In       string // "path" or "query"
	Type     string // "string", "integer", "boolean"
	Required bool
}

// BodyFieldDef describes a request body field, potentially nested.
type BodyFieldDef struct {
	Path       string // relative field name (e.g., "hash" inside a "target" object)
	Type       string // "string", "integer", "boolean"
	Desc       string
	Required   bool           // true when the API schema lists this field as required
	IsArray    bool           // true when the field is an array
	IsObject   bool           // true when the field is a nested object
	ItemFields []BodyFieldDef // nested fields for array items or object properties
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
	Category    string // human-readable API group (e.g., "Pull Requests"); used as subcategory in docs
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
	resp.TypeName = req.ProviderTypeName + "_" + toSnakeCase(r.group.TypeName)
}

// Schema builds the resource schema dynamically from the operation definitions.
func (r *GenericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := resourceBaseAttrs()
	ops := r.crudOps()
	paramSeen := map[string]bool{}
	addResourceParams(attrs, ops, requiredPrimaryPathParams(r.primaryOp()), paramSeen)
	hasBody := addResourceBodyFields(attrs, ops)
	addResourceResponseFields(attrs, r.responseOp(), paramSeen)
	addRequestBodyAttr(attrs, hasBody)

	resp.Schema = schema.Schema{
		Description: r.group.Description,
		Attributes:  attrs,
	}
}

// buildNestedItemAttrs creates schema attributes for nested fields (array items
// or object properties). Recursively handles nested objects and arrays.
func buildNestedItemAttrs(itemFields []BodyFieldDef) map[string]schema.Attribute {
	nested := map[string]schema.Attribute{}
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		desc := f.Desc
		if desc == "" {
			desc = f.Path
		}
		if f.IsObject && len(f.ItemFields) > 0 {
			nested[key] = schema.SingleNestedAttribute{
				Description: desc,
				Optional:    true,
				Computed:    true,
				Attributes:  buildNestedItemAttrs(f.ItemFields),
			}
		} else if f.IsArray && len(f.ItemFields) > 0 {
			nested[key] = schema.ListNestedAttribute{
				Description: desc,
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: buildNestedItemAttrs(f.ItemFields),
				},
			}
		} else if f.IsArray {
			nested[key] = schema.ListAttribute{
				Description: desc,
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			}
		} else {
			nested[key] = schema.StringAttribute{
				Description: desc,
				Optional:    true,
				Computed:    true,
			}
		}
	}
	return nested
}

// itemAttrTypes returns the attr.Type map for nested attribute items.
// Recursively handles nested objects and arrays.
func itemAttrTypes(itemFields []BodyFieldDef) map[string]attr.Type {
	attrTypes := map[string]attr.Type{}
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		if f.IsObject && len(f.ItemFields) > 0 {
			attrTypes[key] = types.ObjectType{AttrTypes: itemAttrTypes(f.ItemFields)}
		} else if f.IsArray && len(f.ItemFields) > 0 {
			attrTypes[key] = types.ListType{ElemType: types.ObjectType{AttrTypes: itemAttrTypes(f.ItemFields)}}
		} else if f.IsArray {
			attrTypes[key] = types.ListType{ElemType: types.StringType}
		} else {
			attrTypes[key] = types.StringType
		}
	}
	return attrTypes
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

// Create calls the create API operation and stores the result in state. When
// the resource also defines a Read operation that differs from Create, a
// follow-up Read is performed against the freshly-written state so that
// computed attributes are populated from the canonical Read response (and any
// Read-side response-shape transformer, e.g. the project branch-restrictions
// `group-by-branch` reshape, is applied). Without this follow-up, Create-only
// responses whose shape diverges from the Read schema leave Computed
// attributes unset, producing "Provider produced inconsistent result after
// apply" errors on every fresh `terraform apply`.
func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	op := r.group.Ops.Create
	if op == nil {
		resp.Diagnostics.AddError("Create not supported", fmt.Sprintf("Resource %s does not support create", r.group.TypeName))
		return
	}
	r.dispatch(ctx, op, &req.Plan, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.refreshAfterWrite(ctx, op, &resp.State, nil, &resp.Diagnostics)
}

// refreshAfterWrite performs a post-write Read against the freshly-written
// state when the Read operation differs from the write operation that just
// ran. This is the generic mechanism by which any resource whose write
// response shape diverges from the Read schema (e.g. the project
// branch-restrictions endpoints, where the write PUT response is not the
// shape the Read transformer expects) gets its state populated from the
// canonical Read response. Without this follow-up, Terraform sees mismatches
// between the planned state and the post-apply state and aborts with
// "Provider produced inconsistent result after apply".
//
// Both Create and Update funnel through this helper so the refresh
// behaviour stays symmetric: a write op that needs a follow-up Read after
// Create needs the same follow-up after Update for the same reasons.
//
// paramFallback supplies an additional state to consult when a path/query
// param required by the Read op is null/unknown/empty in the freshly-
// written state. Update passes the prior state here so that Computed-only
// required params (e.g. a numeric "id" that is "(known after apply)" in
// the plan but present in prior state) can still satisfy the post-write
// Read. Create passes nil — its written state always carries every Read
// param either from the user-supplied plan or from populateComputedParams.
func (r *GenericResource) refreshAfterWrite(ctx context.Context, writeOp *OperationDef, state, paramFallback stateAccessor, diags *diag.Diagnostics) {
	readOp := r.group.Ops.Read
	if readOp == nil || readOp.OperationID == writeOp.OperationID {
		return
	}
	r.refreshState(ctx, readOp, state, paramFallback, diags)
}

// Read calls the read API operation and refreshes state. The resource `id` is
// preserved from the prior state across Read: the id is established once at
// Create time (from the Create operation and its path params) and must not
// change across refreshes. Allowing Read to overwrite it would route
// subsequent Update/Delete calls to the wrong operation/path — for example,
// project branch-restrictions whose Read GET path lacks the `pattern`
// segment, which after a refresh would produce a Delete against the
// `group-by-branch` endpoint and a `400: values array must be specified`
// from Bitbucket.
func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	op := r.group.Ops.Read
	if op == nil {
		// If no read operation, preserve existing state.
		return
	}
	priorID := readPriorID(ctx, &req.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.dispatch(ctx, op, &req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	restorePriorID(ctx, &resp.State, priorID, &resp.Diagnostics)
}

// refreshState performs a Read-style dispatch using the current state as
// both source and target, after preserving the resource id. It is invoked
// after a successful write so the post-write state matches the canonical
// Read response.
//
// paramFallback, when non-nil, is consulted by the dispatcher when a
// path/query param required by readOp is null/unknown/empty in state.
// Update supplies the prior state here so that Computed-only required
// params can still satisfy the Read; pass nil when no fallback applies.
func (r *GenericResource) refreshState(ctx context.Context, readOp *OperationDef, state, paramFallback stateAccessor, diags *diag.Diagnostics) {
	priorID := readPriorID(ctx, state, diags)
	if diags.HasError() {
		return
	}
	r.dispatchWithParamFallback(ctx, readOp, state, paramFallback, state, diags)
	if diags.HasError() {
		return
	}
	restorePriorID(ctx, state, priorID, diags)
}

// readPriorID returns the existing `id` attribute from state. Diagnostics
// returned by the framework's GetAttribute call are appended to diags so
// real state-handling problems (e.g. schema mismatches) surface to the user
// instead of being silently swallowed; callers should bail out via
// diags.HasError() when that happens. An empty string is returned when the
// attribute is null or unknown — i.e. no id has been written yet (e.g. the
// Create dispatch failed before storeDispatchResult ran). Callers treat the
// empty case as "no id to preserve" and let the dispatch-written id stand.
func readPriorID(ctx context.Context, state stateAccessor, diags *diag.Diagnostics) string {
	var id types.String
	diags.Append(state.GetAttribute(ctx, attrPath("id"), &id)...)
	if diags.HasError() || id.IsNull() || id.IsUnknown() {
		return ""
	}
	return id.ValueString()
}

// restorePriorID writes priorID back to state when non-empty. This undoes the
// id (re)write performed by storeDispatchResult during a Read-style dispatch.
func restorePriorID(ctx context.Context, state stateAccessor, priorID string, diags *diag.Diagnostics) {
	if priorID == "" {
		return
	}
	diags.Append(state.SetAttribute(ctx, attrPath("id"), types.StringValue(priorID))...)
}

// Update calls the update API operation and updates state.
//
// Path and query parameters are read from the plan first and fall back to the
// prior state. The fallback is required for resources whose identifying path
// parameter is Computed-only (e.g., the numeric "id" of a branch restriction
// rule returned by Create): such params appear as "(known after apply)" in the
// plan even for in-place updates, so without the fallback the dispatch would
// fail with "Missing Required Parameter".
//
// As with Create, when the Read operation differs from Update a follow-up
// Read is performed against the freshly-written state so the canonical Read
// response (and any Read-side response-shape transformer) populates state.
// Without this, write responses whose shape diverges from the Read schema
// (e.g. the project branch-restrictions PUT endpoints, whose response is not
// in the flat `{"values": [...]}` form the schema declares) leave nested
// Computed attributes mis-mapped, producing "Provider produced inconsistent
// result after apply" errors on every modify.
func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	op := r.group.Ops.Update
	if op == nil {
		resp.Diagnostics.AddError("Update not supported", fmt.Sprintf("Resource %s does not support update", r.group.TypeName))
		return
	}
	r.dispatchWithParamFallback(ctx, op, &req.Plan, &req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.refreshAfterWrite(ctx, op, &resp.State, &req.State, &resp.Diagnostics)
}

// ImportState implements resource import. The import ID must be the slash-separated
// required path parameter values in the order they appear in the URL path template
// (e.g. "my-workspace/my-repo-slug" for /repositories/{workspace}/{repo_slug}).
// After import, a Read is performed automatically to populate all computed attributes.
func (r *GenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	readOp := r.group.Ops.Read
	if readOp == nil {
		resp.Diagnostics.AddError("Import not supported", fmt.Sprintf("Resource %s has no Read operation", r.group.TypeName))
		return
	}

	// Collect required path params in URL template order (e.g. {workspace} before {repo_slug}).
	requiredSet := map[string]bool{}
	for _, p := range readOp.Params {
		if p.In == "path" && p.Required {
			requiredSet[p.Name] = true
		}
	}
	var paramNames []string
	// Extract {param} placeholders from the path in order.
	pathTemplate := readOp.Path
	for len(pathTemplate) > 0 {
		start := strings.Index(pathTemplate, "{")
		if start == -1 {
			break
		}
		end := strings.Index(pathTemplate[start:], "}")
		if end == -1 {
			break
		}
		name := pathTemplate[start+1 : start+end]
		if requiredSet[name] {
			paramNames = append(paramNames, ParamAttrName(name))
		}
		pathTemplate = pathTemplate[start+end+1:]
	}

	if len(paramNames) == 0 {
		// Fallback: treat ID as the resource id attribute.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("id"), types.StringValue(req.ID))...)
		return
	}

	parts := strings.Split(req.ID, "/")
	if len(parts) < len(paramNames) {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected %d slash-separated values (%s), got %q",
				len(paramNames), strings.Join(paramNames, "/"), req.ID))
		return
	}

	// When the last param may itself contain slashes (e.g. file paths), join
	// any extra parts back into the last segment.
	values := make([]string, len(paramNames))
	copy(values, parts[:len(paramNames)-1])
	values[len(paramNames)-1] = strings.Join(parts[len(paramNames)-1:], "/")

	for i, name := range paramNames {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath(name), types.StringValue(values[i]))...)
	}
}

// Delete calls the delete API operation.
func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	op := r.group.Ops.Delete
	if op == nil {
		// If no delete, just remove from state.
		resp.State.RemoveResource(ctx)
		return
	}

	plan := &req.State
	pathParams, queryParams := r.extractParams(ctx, op, plan, nil, &resp.Diagnostics)
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
		return
	}

	resp.State.RemoveResource(ctx)
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// dispatch executes an API operation, reading params from source and writing results to target.
func (r *GenericResource) dispatch(ctx context.Context, op *OperationDef, source stateAccessor, target stateAccessor, diags *diag.Diagnostics) {
	r.dispatchWithParamFallback(ctx, op, source, nil, target, diags)
}

// dispatchWithParamFallback is like dispatch but consults paramFallback when a
// path/query parameter is null/unknown/empty in the primary source. The body is
// always built from the primary source.
func (r *GenericResource) dispatchWithParamFallback(ctx context.Context, op *OperationDef, source, paramFallback stateAccessor, target stateAccessor, diags *diag.Diagnostics) {
	pathParams, queryParams := r.extractParams(ctx, op, source, paramFallback, diags)
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

	// Reshape responses whose on-the-wire shape differs from the schema the
	// generic extractor expects (e.g. the project branch-restrictions
	// `group-by-branch` GET).
	result = transformProjectBranchRestrictionsRead(ctx, op, r.group.TypeName, source, result, diags)

	resultMap := r.storeDispatchResult(ctx, op, pathParams, source, target, diags, result)
	if resultMap != nil {
		r.populateComputedParams(ctx, resultMap, source, target, diags)
	}
	r.copyAttributes(ctx, op, source, target, diags)
}

// extractParams reads path and query params from the source state/plan. When a
// value is null/unknown/empty in source and fallback is non-nil, the value is
// looked up in fallback instead. This is used during Update where Computed-only
// path params (e.g., a resource id returned by Create) are unknown in the plan
// and must be read from prior state.
func (r *GenericResource) extractParams(ctx context.Context, op *OperationDef, source, fallback stateAccessor, diags *diag.Diagnostics) (map[string]string, map[string]string) {
	pathParams := map[string]string{}
	queryParams := map[string]string{}

	for _, p := range op.Params {
		attrName := ParamAttrName(p.Name)
		var val types.String
		d := source.GetAttribute(ctx, attrPath(attrName), &val)
		diags.Append(d...)
		if d.HasError() {
			continue
		}
		if val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
			if fallback != nil {
				var fb types.String
				fd := fallback.GetAttribute(ctx, attrPath(attrName), &fb)
				if !fd.HasError() && !fb.IsNull() && !fb.IsUnknown() && fb.ValueString() != "" {
					val = fb
				}
			}
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
	if !op.HasBody {
		return ""
	}

	if len(op.BodyFields) > 0 {
		return marshalBodyObject(diags, buildExplicitBody(ctx, source, op.BodyFields, diags))
	}

	return rawRequestBody(ctx, source, diags)
}

// readListNested reads a ListNestedAttribute from state and returns it as a
// []map[string]any suitable for JSON marshaling. Returns nil if the list is
// null, unknown, or empty.
func readListNested(ctx context.Context, source stateAccessor, attrName string, itemFields []BodyFieldDef, diags *diag.Diagnostics) []map[string]any {
	var list types.List
	d := source.GetAttribute(ctx, attrPath(attrName), &list)
	diags.Append(d...)
	if d.HasError() || list.IsNull() || list.IsUnknown() {
		return nil
	}
	return readListNestedValue(list, itemFields)
}

func readListNestedValue(list types.List, itemFields []BodyFieldDef) []map[string]any {
	elements := list.Elements()
	if len(elements) == 0 {
		return nil
	}
	var result []map[string]any
	for _, elem := range elements {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		item := map[string]any{}
		objAttrs := obj.Attributes()
		for _, f := range itemFields {
			key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
			if v, exists := objAttrs[key]; exists {
				if value, ok := readAttrValue(v, f); ok {
					item[f.Path] = value
				}
			}
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	return result
}

// readSingleNested reads a SingleNestedAttribute from state and returns it as a
// map[string]any suitable for JSON marshaling. Returns nil if the object is
// null, unknown, or all fields are empty.
func readSingleNested(ctx context.Context, source stateAccessor, attrName string, itemFields []BodyFieldDef, diags *diag.Diagnostics) map[string]any {
	var obj types.Object
	d := source.GetAttribute(ctx, attrPath(attrName), &obj)
	diags.Append(d...)
	if d.HasError() || obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	result := map[string]any{}
	objAttrs := obj.Attributes()
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		v, exists := objAttrs[key]
		if !exists {
			continue
		}
		if value, ok := readAttrValue(v, f); ok {
			result[f.Path] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// readObjectAttrs reads attribute values from a types.Object, used by
// readSingleNested for recursively reading nested objects.
func readObjectAttrs(obj types.Object, itemFields []BodyFieldDef) map[string]any {
	result := map[string]any{}
	objAttrs := obj.Attributes()
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		v, exists := objAttrs[key]
		if !exists {
			continue
		}
		if value, ok := readAttrValue(v, f); ok {
			result[f.Path] = value
		}
	}
	return result
}

func readAttrValue(v attr.Value, f BodyFieldDef) (any, bool) {
	switch {
	case f.IsObject && len(f.ItemFields) > 0:
		return readObjectAttrValue(v, f.ItemFields)
	case f.IsArray && len(f.ItemFields) > 0:
		return readNestedListAttrValue(v, f.ItemFields)
	case f.IsArray:
		return readSimpleListAttrValue(v)
	default:
		return readStringAttrValue(v)
	}
}

// readSimpleList reads a ListAttribute (list of strings) from state and returns
// it as a []string suitable for JSON marshaling. Returns nil if null/unknown/empty.
func readSimpleList(ctx context.Context, source stateAccessor, attrName string, diags *diag.Diagnostics) []string {
	var list types.List
	d := source.GetAttribute(ctx, attrPath(attrName), &list)
	diags.Append(d...)
	if d.HasError() || list.IsNull() || list.IsUnknown() {
		return nil
	}
	return readSimpleListValue(list)
}

func readSimpleListValue(list types.List) []string {
	elements := list.Elements()
	if len(elements) == 0 {
		return nil
	}
	var result []string
	for _, elem := range elements {
		if sv, ok := elem.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
			result = append(result, sv.ValueString())
		}
	}
	return result
}

// copyAttributes copies all param and body field values from source to target state.
// It copies from ALL CRUD operations so that attributes defined by one operation
// (e.g., Read's path params) are preserved when another operation (e.g., Create) runs.
func (r *GenericResource) copyAttributes(ctx context.Context, _ *OperationDef, source, target stateAccessor, diags *diag.Diagnostics) {
	seen := map[string]bool{}
	hasBody := false
	for _, op := range r.crudOps() {
		if op.HasBody {
			hasBody = true
		}
		copyParamAttributes(ctx, op.Params, seen, source, target, diags)
		copyBodyAttributes(ctx, op.BodyFields, seen, source, target, diags)
	}
	if hasBody {
		copyStringAttribute(ctx, "request_body", source, target, diags)
	}
}

// populateComputedParams sets path param attributes from the API response for
// params that the user did not provide (Optional+Computed). For example, after
// a Create on branch-restrictions, the API returns the new restriction's "id"
// which is needed for subsequent Read/Update/Delete operations.
func (r *GenericResource) populateComputedParams(ctx context.Context, m map[string]any, source, target stateAccessor, diags *diag.Diagnostics) {
	seen := map[string]bool{}
	for _, crudOp := range r.crudOps() {
		for _, p := range crudOp.Params {
			attrName := ParamAttrName(p.Name)
			if shouldSkipComputedParam(p, attrName, seen, ctx, source, diags) {
				continue
			}
			setComputedParam(ctx, m, p.Name, attrName, target, diags)
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

// responseFields returns the response fields from the Read operation (preferred)
// or Create operation as fallback.
func (r *GenericResource) responseFields() []BodyFieldDef {
	if r.group.Ops.Read != nil && len(r.group.Ops.Read.ResponseFields) > 0 {
		return r.group.Ops.Read.ResponseFields
	}
	if r.group.Ops.Create != nil && len(r.group.Ops.Create.ResponseFields) > 0 {
		return r.group.Ops.Create.ResponseFields
	}
	return nil
}

// extractResponseFields extracts individual field values from the API response
// and sets them as computed attributes in the target state. The optional
// `source` accessor (typically the plan or prior state used to dispatch the
// API call) is consulted to align response-side nested-object lists with
// the configured/planned order, preserving order-insensitivity without
// mutating the planned value.
func (r *GenericResource) extractResponseFields(ctx context.Context, m map[string]any, source, target stateAccessor, diags *diag.Diagnostics) {
	for _, rf := range r.responseFields() {
		setResponseField(ctx, rf, m, source, target, diags)
	}
}

func resourceBaseAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Resource identifier.",
			Computed:    true,
		},
		"api_response": schema.StringAttribute{
			Description: "The raw JSON response from the Bitbucket API.",
			Computed:    true,
		},
	}
}

func (r *GenericResource) primaryOp() *OperationDef {
	if r.group.Ops.Create != nil {
		return r.group.Ops.Create
	}
	return r.group.Ops.Read
}

func (r *GenericResource) responseOp() *OperationDef {
	if r.group.Ops.Read != nil {
		return r.group.Ops.Read
	}
	return r.group.Ops.Create
}

func requiredPrimaryPathParams(op *OperationDef) map[string]bool {
	required := map[string]bool{}
	if op == nil {
		return required
	}
	for _, p := range op.Params {
		if p.Required && p.In == "path" {
			required[p.Name] = true
		}
	}
	return required
}

func addResourceParams(attrs map[string]schema.Attribute, ops []*OperationDef, primaryRequired map[string]bool, paramSeen map[string]bool) {
	for _, op := range ops {
		for _, p := range op.Params {
			if paramSeen[p.Name] {
				continue
			}
			paramSeen[p.Name] = true
			attrName := ParamAttrName(p.Name)
			if _, exists := attrs[attrName]; exists {
				continue
			}
			attrs[attrName] = resourceParamAttr(p, primaryRequired[p.Name])
		}
	}
}

func resourceParamAttr(p ParamDef, isRequired bool) schema.StringAttribute {
	attr := schema.StringAttribute{
		Description: fmt.Sprintf("%s parameter (%s)", p.Name, p.In),
	}
	if isRequired {
		attr.Required = true
		return attr
	}
	attr.Optional = true
	attr.Computed = p.In == "path"
	return attr
}

func addResourceBodyFields(attrs map[string]schema.Attribute, ops []*OperationDef) bool {
	bodyFieldSeen := map[string]bool{}
	hasBody := false
	for _, op := range ops {
		if op.HasBody {
			hasBody = true
		}
		for _, bf := range op.BodyFields {
			key := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
			if bodyFieldSeen[key] {
				continue
			}
			bodyFieldSeen[key] = true
			attrs[key] = bodyFieldAttr(bf)
		}
	}
	return hasBody
}

func bodyFieldAttr(bf BodyFieldDef) schema.Attribute {
	desc := bf.Desc
	if desc == "" {
		desc = bf.Path
	}
	if bf.IsObject && len(bf.ItemFields) > 0 {
		return schema.SingleNestedAttribute{
			Description: desc,
			Optional:    true,
			Attributes:  buildNestedItemAttrs(bf.ItemFields),
		}
	}
	if bf.IsArray && len(bf.ItemFields) > 0 {
		return schema.ListNestedAttribute{
			Description: desc,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: buildNestedItemAttrs(bf.ItemFields),
			},
		}
	}
	if bf.IsArray {
		return schema.ListAttribute{
			Description: desc,
			Optional:    true,
			ElementType: types.StringType,
		}
	}
	if bf.Type == "int" {
		if bf.Required {
			return schema.Int64Attribute{
				Description: desc,
				Required:    true,
			}
		}
		return schema.Int64Attribute{
			Description: desc,
			Optional:    true,
		}
	}
	if bf.Required {
		return schema.StringAttribute{
			Description: desc,
			Required:    true,
		}
	}
	return schema.StringAttribute{
		Description: desc,
		Optional:    true,
	}
}

func addResourceResponseFields(attrs map[string]schema.Attribute, responseOp *OperationDef, paramSeen map[string]bool) {
	if responseOp == nil {
		return
	}
	for _, rf := range responseOp.ResponseFields {
		key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
		if isReservedResourceAttr(key) {
			continue
		}
		if existing, exists := attrs[key]; exists {
			attrs[key] = mergeResponseAttr(existing, rf)
			continue
		}
		if !paramSeen[key] {
			attrs[key] = responseFieldAttr(rf)
		}
	}
}

func isReservedResourceAttr(key string) bool {
	return key == "id" || key == "api_response" || key == "request_body"
}

func responseFieldAttr(rf BodyFieldDef) schema.Attribute {
	desc := rf.Desc
	if desc == "" {
		desc = rf.Path
	}
	if rf.IsObject && len(rf.ItemFields) > 0 {
		return schema.SingleNestedAttribute{
			Description: desc,
			Computed:    true,
			Attributes:  buildNestedItemAttrs(rf.ItemFields),
		}
	}
	if rf.IsArray && len(rf.ItemFields) > 0 {
		return schema.ListNestedAttribute{
			Description: desc,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: buildNestedItemAttrs(rf.ItemFields),
			},
		}
	}
	if rf.IsArray {
		return schema.ListAttribute{
			Description: desc,
			Computed:    true,
			ElementType: types.StringType,
		}
	}
	if rf.Type == "int" {
		return schema.Int64Attribute{
			Description: desc,
			Computed:    true,
		}
	}
	return schema.StringAttribute{
		Description: desc,
		Computed:    true,
	}
}

func mergeResponseAttr(existing schema.Attribute, rf BodyFieldDef) schema.Attribute {
	switch sa := existing.(type) {
	case schema.StringAttribute:
		return mergeStringResponseAttr(sa, rf)
	case schema.Int64Attribute:
		return mergeInt64ResponseAttr(sa, rf)
	case schema.SingleNestedAttribute:
		return mergeSingleNestedResponseAttr(sa, rf)
	case schema.ListNestedAttribute:
		return mergeListNestedResponseAttr(sa, rf)
	case schema.ListAttribute:
		return mergeListResponseAttr(sa, rf)
	default:
		return existing
	}
}

func mergeStringResponseAttr(attr schema.StringAttribute, rf BodyFieldDef) schema.Attribute {
	if !canMergeComputedAttr(attr.Computed, attr.Required) {
		return attr
	}
	attr.Computed = true
	attr.Description = fieldDescription(rf)
	return attr
}

func mergeInt64ResponseAttr(attr schema.Int64Attribute, rf BodyFieldDef) schema.Attribute {
	if !canMergeComputedAttr(attr.Computed, attr.Required) {
		return attr
	}
	attr.Computed = true
	attr.Description = fieldDescription(rf)
	return attr
}

func mergeSingleNestedResponseAttr(attr schema.SingleNestedAttribute, rf BodyFieldDef) schema.Attribute {
	if !canMergeComputedAttr(attr.Computed, attr.Required) {
		return attr
	}
	attr.Computed = true
	attr.Description = fieldDescription(rf)
	if rf.IsObject && len(rf.ItemFields) > 0 {
		attr.Attributes = buildNestedItemAttrs(rf.ItemFields)
	}
	return attr
}

func mergeListNestedResponseAttr(attr schema.ListNestedAttribute, rf BodyFieldDef) schema.Attribute {
	if !canMergeComputedAttr(attr.Computed, attr.Required) {
		return attr
	}
	attr.Computed = true
	attr.Description = fieldDescription(rf)
	if rf.IsArray && len(rf.ItemFields) > 0 {
		attr.NestedObject = schema.NestedAttributeObject{
			Attributes: buildNestedItemAttrs(rf.ItemFields),
		}
	}
	return attr
}

func mergeListResponseAttr(attr schema.ListAttribute, rf BodyFieldDef) schema.Attribute {
	if !canMergeComputedAttr(attr.Computed, attr.Required) {
		return attr
	}
	attr.Computed = true
	attr.Description = fieldDescription(rf)
	return attr
}

func canMergeComputedAttr(computed, required bool) bool {
	return !computed && !required
}

func fieldDescription(field BodyFieldDef) string {
	if field.Desc != "" {
		return field.Desc
	}
	return field.Path
}

func addRequestBodyAttr(attrs map[string]schema.Attribute, hasBody bool) {
	if !hasBody {
		return
	}
	attrs["request_body"] = schema.StringAttribute{
		Description: "Raw JSON request body for create/update operations. Use this to pass fields not exposed as individual attributes.",
		Optional:    true,
	}
}

func (r *GenericResource) storeDispatchResult(ctx context.Context, op *OperationDef, pathParams map[string]string, source, target stateAccessor, diags *diag.Diagnostics, result any) map[string]any {
	if result == nil {
		diags.Append(target.SetAttribute(ctx, attrPath("api_response"), types.StringValue(""))...)
		diags.Append(target.SetAttribute(ctx, attrPath("id"), types.StringValue(op.OperationID))...)
		return nil
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		diags.AddError("Response Error", fmt.Sprintf("Failed to marshal response: %v", err))
		return nil
	}
	diags.Append(target.SetAttribute(ctx, attrPath("api_response"), types.StringValue(string(jsonBytes)))...)

	resultMap, _ := result.(map[string]any)
	if resultMap != nil {
		r.extractResponseFields(ctx, resultMap, source, target, diags)
		if id := extractID(resultMap); id != "" {
			diags.Append(target.SetAttribute(ctx, attrPath("id"), types.StringValue(id))...)
			return resultMap
		}
	}
	diags.Append(target.SetAttribute(ctx, attrPath("id"), types.StringValue(fallbackResourceID(op, pathParams)))...)
	return resultMap
}

func fallbackResourceID(op *OperationDef, pathParams map[string]string) string {
	fallbackID := op.OperationID
	for _, p := range op.Params {
		if p.In == "path" {
			if v, ok := pathParams[p.Name]; ok {
				fallbackID += "/" + v
			}
		}
	}
	return fallbackID
}

func buildExplicitBody(ctx context.Context, source stateAccessor, fields []BodyFieldDef, diags *diag.Diagnostics) map[string]any {
	bodyObj := map[string]any{}
	for _, bf := range fields {
		if value, ok := bodyFieldValue(ctx, source, bf, diags); ok {
			handlers.SetNested(bodyObj, bf.Path, value)
		}
	}
	return bodyObj
}

func bodyFieldValue(ctx context.Context, source stateAccessor, bf BodyFieldDef, diags *diag.Diagnostics) (any, bool) {
	attrName := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
	switch {
	case bf.IsObject && len(bf.ItemFields) > 0:
		obj := readSingleNested(ctx, source, attrName, bf.ItemFields, diags)
		return obj, obj != nil
	case bf.IsArray && len(bf.ItemFields) > 0:
		arr := readListNested(ctx, source, attrName, bf.ItemFields, diags)
		return arr, arr != nil
	case bf.IsArray:
		arr := readSimpleList(ctx, source, attrName, diags)
		return arr, arr != nil
	case bf.Type == "int":
		return readBodyInt64Value(ctx, source, attrName, diags)
	default:
		return readBodyStringValue(ctx, source, attrName, diags)
	}
}

func readBodyStringValue(ctx context.Context, source stateAccessor, attrName string, diags *diag.Diagnostics) (string, bool) {
	var val types.String
	d := source.GetAttribute(ctx, attrPath(attrName), &val)
	diags.Append(d...)
	if d.HasError() || val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
		return "", false
	}
	return val.ValueString(), true
}

func readBodyInt64Value(ctx context.Context, source stateAccessor, attrName string, diags *diag.Diagnostics) (int64, bool) {
	var val types.Int64
	d := source.GetAttribute(ctx, attrPath(attrName), &val)
	diags.Append(d...)
	if d.HasError() || val.IsNull() || val.IsUnknown() {
		return 0, false
	}
	return val.ValueInt64(), true
}

func marshalBodyObject(diags *diag.Diagnostics, bodyObj map[string]any) string {
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

func rawRequestBody(ctx context.Context, source stateAccessor, diags *diag.Diagnostics) string {
	body, ok := readBodyStringValue(ctx, source, "request_body", diags)
	if !ok {
		return ""
	}
	return body
}

func readObjectAttrValue(v attr.Value, itemFields []BodyFieldDef) (any, bool) {
	innerObj, ok := v.(types.Object)
	if !ok || innerObj.IsNull() || innerObj.IsUnknown() {
		return nil, false
	}
	inner := readObjectAttrs(innerObj, itemFields)
	if len(inner) == 0 {
		return nil, false
	}
	return inner, true
}

func readNestedListAttrValue(v attr.Value, itemFields []BodyFieldDef) (any, bool) {
	list, ok := v.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return nil, false
	}
	items := readListNestedValue(list, itemFields)
	if len(items) == 0 {
		return nil, false
	}
	return items, true
}

func readSimpleListAttrValue(v attr.Value) (any, bool) {
	list, ok := v.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return nil, false
	}
	items := readSimpleListValue(list)
	if len(items) == 0 {
		return nil, false
	}
	return items, true
}

func readStringAttrValue(v attr.Value) (any, bool) {
	sv, ok := v.(types.String)
	if !ok || sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
		return nil, false
	}
	return sv.ValueString(), true
}

func copyParamAttributes(ctx context.Context, params []ParamDef, seen map[string]bool, source, target stateAccessor, diags *diag.Diagnostics) {
	for _, p := range params {
		attrName := ParamAttrName(p.Name)
		if seen[attrName] {
			continue
		}
		seen[attrName] = true
		copyStringAttribute(ctx, attrName, source, target, diags)
	}
}

func copyBodyAttributes(ctx context.Context, fields []BodyFieldDef, seen map[string]bool, source, target stateAccessor, diags *diag.Diagnostics) {
	for _, bf := range fields {
		attrName := toSnakeCase(strings.ReplaceAll(bf.Path, ".", "_"))
		if seen[attrName] {
			continue
		}
		seen[attrName] = true
		if bf.IsArray || bf.IsObject {
			continue
		}
		if bf.Type == "int" {
			copyInt64Attribute(ctx, attrName, source, target, diags)
		} else {
			copyStringAttribute(ctx, attrName, source, target, diags)
		}
	}
}

func copyStringAttribute(ctx context.Context, attrName string, source, target stateAccessor, diags *diag.Diagnostics) {
	var val types.String
	d := source.GetAttribute(ctx, attrPath(attrName), &val)
	diags.Append(d...)
	if !d.HasError() && !val.IsNull() && !val.IsUnknown() {
		diags.Append(target.SetAttribute(ctx, attrPath(attrName), val)...)
	}
}

func copyInt64Attribute(ctx context.Context, attrName string, source, target stateAccessor, diags *diag.Diagnostics) {
	var val types.Int64
	d := source.GetAttribute(ctx, attrPath(attrName), &val)
	diags.Append(d...)
	if !d.HasError() && !val.IsNull() && !val.IsUnknown() {
		diags.Append(target.SetAttribute(ctx, attrPath(attrName), val)...)
	}
}

func shouldSkipComputedParam(p ParamDef, attrName string, seen map[string]bool, ctx context.Context, source stateAccessor, diags *diag.Diagnostics) bool {
	if p.In != "path" || seen[attrName] {
		return true
	}
	seen[attrName] = true
	var existing types.String
	d := source.GetAttribute(ctx, attrPath(attrName), &existing)
	diags.Append(d...)
	return d.HasError() || (!existing.IsNull() && !existing.IsUnknown() && existing.ValueString() != "")
}

func setComputedParam(ctx context.Context, m map[string]any, paramName, attrName string, target stateAccessor, diags *diag.Diagnostics) {
	if val, ok := responseParamValue(m, paramName); ok && val != "" {
		diags.Append(target.SetAttribute(ctx, attrPath(attrName), types.StringValue(fmt.Sprintf("%v", val)))...)
	}
}

func setResponseField(ctx context.Context, rf BodyFieldDef, m map[string]any, source, target stateAccessor, diags *diag.Diagnostics) {
	key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
	if isReservedResourceAttr(key) {
		return
	}
	val, ok := handlers.GetNested(m, rf.Path)
	if !ok || val == nil {
		return
	}
	// For nested-object lists, read the prior/planned value from `source`
	// so the API response can be reordered to match it (preserves the
	// operator's order without mutating the plan). Errors reading the prior
	// value are non-fatal: we fall back to deterministic sorting.
	var priorList types.List
	if source != nil && rf.IsArray && len(rf.ItemFields) > 0 {
		var probe types.List
		if d := source.GetAttribute(ctx, attrPath(key), &probe); !d.HasError() {
			priorList = probe
		}
	}
	attrValue, ok := responseFieldValue(val, rf, priorList)
	if !ok {
		return
	}
	diags.Append(target.SetAttribute(ctx, attrPath(key), attrValue)...)
}

func responseFieldValue(val any, rf BodyFieldDef, priorList types.List) (any, bool) {
	if rf.IsObject && len(rf.ItemFields) > 0 {
		obj, ok := val.(map[string]any)
		if !ok {
			return nil, false
		}
		return buildObjectFromResponse(obj, rf.ItemFields), true
	}
	if rf.IsArray && len(rf.ItemFields) > 0 {
		arr, ok := val.([]any)
		if !ok {
			return nil, false
		}
		return buildListFromResponse(arr, rf.ItemFields, priorList), true
	}
	if rf.IsArray {
		arr, ok := val.([]any)
		if !ok {
			return nil, false
		}
		return buildSimpleListFromResponse(arr), true
	}
	if rf.Type == "int" {
		switch v := val.(type) {
		case float64:
			return types.Int64Value(int64(v)), true
		case int64:
			return types.Int64Value(v), true
		case int:
			return types.Int64Value(int64(v)), true
		}
		return nil, false
	}
	return types.StringValue(stringifyComplexValue(val)), true
}

func stringifyComplexValue(val any) string {
	switch val.(type) {
	case []any, map[string]any:
		if b, err := json.Marshal(val); err == nil {
			return string(b)
		}
	}
	return fmt.Sprintf("%v", val)
}

// buildListFromResponse converts a JSON array from the API response into a
// types.List value suitable for a ListNestedAttribute. When `priorList` is
// a non-null/non-unknown reference list (e.g. the planned/prior state
// value), response items are reordered to match that reference order via
// stable identity keys — this is what makes nested-object lists
// order-insensitive without ever mutating the plan. When no reference
// exists (initial Read with no prior state, or data source reads), items
// are sorted by a stable identity key with a canonical JSON tiebreaker so
// two equivalent API responses still produce byte-identical state.
func buildListFromResponse(arr []any, itemFields []BodyFieldDef, priorList types.List) types.List {
	attrTypes := itemAttrTypes(itemFields)
	objType := types.ObjectType{AttrTypes: attrTypes}

	if len(arr) == 0 {
		return types.ListValueMust(objType, []attr.Value{})
	}

	// Work on a copy so we never mutate the caller's slice.
	ordered := make([]any, len(arr))
	copy(ordered, arr)

	aligned := false
	if !priorList.IsNull() && !priorList.IsUnknown() {
		refs := make([]types.Object, 0, len(priorList.Elements()))
		for _, e := range priorList.Elements() {
			if obj, ok := e.(types.Object); ok && !obj.IsNull() && !obj.IsUnknown() {
				refs = append(refs, obj)
			}
		}
		if len(refs) > 0 {
			aligned = alignResponseItemsToReference(ordered, refs, itemFields)
		}
	}
	if !aligned {
		// Fall back (or stabilise leftovers after partial alignment) with
		// the deterministic identity-key sort so equivalent API responses
		// still produce byte-identical state.
		if priorList.IsNull() || priorList.IsUnknown() || len(priorList.Elements()) == 0 {
			sortResponseItems(ordered, itemFields)
		}
	}

	elements := make([]attr.Value, 0, len(ordered))
	for _, item := range ordered {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		objAttrs := map[string]attr.Value{}
		for _, f := range itemFields {
			key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
			v, exists := m[f.Path]
			if !exists || v == nil {
				objAttrs[key] = attrNullValue(f)
				continue
			}
			objAttrs[key] = buildAttrValueFromResponse(v, f)
		}
		elements = append(elements, types.ObjectValueMust(attrTypes, objAttrs))
	}
	return types.ListValueMust(objType, elements)
}

// buildObjectFromResponse converts a JSON object from the API response into a
// types.Object value suitable for a SingleNestedAttribute.
func buildObjectFromResponse(m map[string]any, itemFields []BodyFieldDef) types.Object {
	attrTypes := itemAttrTypes(itemFields)
	objAttrs := map[string]attr.Value{}
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		v, exists := m[f.Path]
		if !exists || v == nil {
			objAttrs[key] = attrNullValue(f)
			continue
		}
		objAttrs[key] = buildAttrValueFromResponse(v, f)
	}
	return types.ObjectValueMust(attrTypes, objAttrs)
}

func buildAttrValueFromResponse(v any, f BodyFieldDef) attr.Value {
	if f.IsObject && len(f.ItemFields) > 0 {
		sub, ok := v.(map[string]any)
		if !ok {
			return attrNullValue(f)
		}
		return buildObjectFromResponse(sub, f.ItemFields)
	}
	if f.IsArray && len(f.ItemFields) > 0 {
		arr, ok := v.([]any)
		if !ok {
			return attrNullValue(f)
		}
		// Nested-inside-object arrays don't have a per-item prior list
		// available here; fall back to deterministic sort. The Required
		// attrs constraint that motivated alignment only applies to
		// top-level lists which are reached via setResponseField.
		return buildListFromResponse(arr, f.ItemFields, types.ListNull(types.ObjectType{AttrTypes: itemAttrTypes(f.ItemFields)}))
	}
	if f.IsArray {
		arr, ok := v.([]any)
		if !ok {
			return attrNullValue(f)
		}
		return buildSimpleListFromResponse(arr)
	}
	return types.StringValue(fmt.Sprintf("%v", v))
}

// attrNullValue returns the appropriate null value for a field's type.
func attrNullValue(f BodyFieldDef) attr.Value {
	if f.IsObject && len(f.ItemFields) > 0 {
		return types.ObjectNull(itemAttrTypes(f.ItemFields))
	}
	if f.IsArray && len(f.ItemFields) > 0 {
		return types.ListNull(types.ObjectType{AttrTypes: itemAttrTypes(f.ItemFields)})
	}
	if f.IsArray {
		return types.ListNull(types.StringType)
	}
	return types.StringNull()
}

// buildSimpleListFromResponse converts a JSON array of simple values into a
// types.List of strings, suitable for a schema.ListAttribute{ElementType: types.StringType}.
func buildSimpleListFromResponse(arr []any) types.List {
	if len(arr) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
	elements := make([]attr.Value, 0, len(arr))
	for _, item := range arr {
		elements = append(elements, types.StringValue(fmt.Sprintf("%v", item)))
	}
	return types.ListValueMust(types.StringType, elements)
}

// extractID tries to extract an identifier from an API response map.
func extractID(m map[string]any) string {
	// Try common ID fields.
	for _, key := range []string{"uuid", "id", "slug", "name"} {
		if v, ok := m[key]; ok {
			return formatResponseValue(v)
		}
	}
	return ""
}

// formatResponseValue converts an API response value to a string,
// formatting float64 integer values (JSON numbers) as plain integers
// to avoid scientific notation (e.g. 7.4332764e+07 instead of 74332764).
func formatResponseValue(v any) string {
	if f, ok := v.(float64); ok && f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%v", v)
}

func responseParamValue(m map[string]any, paramName string) (string, bool) {
	tryKeys := []string{paramName}
	if strings.HasSuffix(paramName, "_uuid") {
		base := strings.TrimSuffix(paramName, "_uuid")
		tryKeys = append(tryKeys, base, "uuid")
	}
	if strings.HasSuffix(paramName, "_id") {
		base := strings.TrimSuffix(paramName, "_id")
		tryKeys = append(tryKeys, base, "id")
	}
	for _, key := range tryKeys {
		if key == "" {
			continue
		}
		if v, ok := m[key]; ok && v != nil {
			return formatResponseValue(v), true
		}
	}
	if id := extractID(m); id != "" {
		return id, true
	}
	return "", false
}
