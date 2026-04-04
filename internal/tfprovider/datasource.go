package tfprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
)

// Ensure the implementation satisfies the datasource interface.
var _ datasource.DataSource = &GenericDataSource{}
var _ datasource.DataSourceWithConfigure = &GenericDataSource{}

// GenericDataSource implements a Terraform data source backed by Bitbucket API operations.
type GenericDataSource struct {
	group  ResourceGroup
	client *client.BBClient
}

// Metadata returns the data source type name.
func (d *GenericDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + toSnakeCase(d.group.TypeName)
}

// Schema builds the data source schema from the Read or List operation.
func (d *GenericDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := dataSourceBaseAttrs()

	op := d.readOp()
	if op == nil {
		resp.Schema = schema.Schema{
			Description: d.group.Description,
			Attributes:  attrs,
		}
		return
	}

	paramSeen := map[string]bool{}
	addDataSourceParams(attrs, op.Params, listRequiredPathParams(d.group.Ops.List), paramSeen)
	addDataSourceListParams(attrs, d.group.Ops.List, paramSeen)
	addDataSourceResponseFields(attrs, op.ResponseFields)

	resp.Schema = schema.Schema{
		Description: d.group.Description + " (data source - read-only)",
		Attributes:  attrs,
	}
}

// Configure receives the provider-configured client.
func (d *GenericDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.BBClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.BBClient, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

// Read fetches data from the Bitbucket API.
func (d *GenericDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Choose operation: if Read op exists and all its path params are provided,
	// use Read. Otherwise fall back to List (if available).
	op := d.selectReadOp(ctx, req)
	if op == nil {
		resp.Diagnostics.AddError("Read not supported",
			fmt.Sprintf("Data source %s has no read operation", d.group.TypeName))
		return
	}

	pathParams, queryParams := readDataSourceParams(ctx, req, resp, op)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := handlers.DispatchRaw(ctx, d.client, handlers.Request{
		Method:      op.Method,
		URLTemplate: op.Path,
		PathParams:  pathParams,
		QueryParams: queryParams,
		All:         op.Paginated,
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Read failed: %v", err))
		return
	}

	setDataSourceResult(ctx, resp, op, pathParams, result)
	copyDataSourceConfigToState(ctx, req, resp, op.Params)
}

// readOp returns the best operation for building the data source schema (Read or List).
func (d *GenericDataSource) readOp() *OperationDef {
	if d.group.Ops.Read != nil {
		return d.group.Ops.Read
	}
	return d.group.Ops.List
}

// selectReadOp chooses which operation to use at runtime:
// If all of Read's required path params are provided, use Read.
// Otherwise fall back to List (if available).
func (d *GenericDataSource) selectReadOp(ctx context.Context, req datasource.ReadRequest) *OperationDef {
	readOp := d.group.Ops.Read
	if readOp != nil {
		allProvided := true
		for _, p := range readOp.Params {
			if p.In != "path" || !p.Required {
				continue
			}
			attrName := ParamAttrName(p.Name)
			var val types.String
			dd := req.Config.GetAttribute(ctx, attrPath(attrName), &val)
			if dd.HasError() || val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
				allProvided = false
				break
			}
		}
		if allProvided {
			return readOp
		}
	}
	// Fall back to List if Read params were incomplete.
	if d.group.Ops.List != nil {
		return d.group.Ops.List
	}
	return readOp // If no List, use Read anyway (errors will be reported for missing params).
}

// buildDSNestedItemAttrs creates data source schema attributes for array item fields.
// All nested fields are Computed in data sources.
func buildDSNestedItemAttrs(itemFields []BodyFieldDef) map[string]schema.Attribute {
	nested := map[string]schema.Attribute{}
	for _, f := range itemFields {
		key := toSnakeCase(strings.ReplaceAll(f.Path, ".", "_"))
		desc := f.Desc
		if desc == "" {
			desc = f.Path
		}
		nested[key] = schema.StringAttribute{
			Description: desc,
			Computed:    true,
		}
	}
	return nested
}

// buildListID creates a composite ID from operation ID and path parameters for list data sources.
// Keys are sorted for deterministic ordering.
func buildListID(operationID string, pathParams map[string]string) string {
	keys := make([]string, 0, len(pathParams))
	for k := range pathParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys)+1)
	if operationID != "" {
		parts = append(parts, operationID)
	}
	for _, k := range keys {
		parts = append(parts, pathParams[k])
	}
	return strings.Join(parts, "/")
}

func dataSourceBaseAttrs() map[string]schema.Attribute {
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

func listRequiredPathParams(op *OperationDef) map[string]bool {
	listPathParams := map[string]bool{}
	if op == nil {
		return listPathParams
	}
	for _, p := range op.Params {
		if p.In == "path" && p.Required {
			listPathParams[p.Name] = true
		}
	}
	return listPathParams
}

func addDataSourceParams(attrs map[string]schema.Attribute, params []ParamDef, listPathParams map[string]bool, paramSeen map[string]bool) {
	for _, p := range params {
		if paramSeen[p.Name] {
			continue
		}
		paramSeen[p.Name] = true
		attrName := ParamAttrName(p.Name)
		if _, exists := attrs[attrName]; exists {
			continue
		}
		attrs[attrName] = schema.StringAttribute{
			Description: fmt.Sprintf("%s parameter (%s)", p.Name, p.In),
			Required:    isDataSourceParamRequired(p, listPathParams),
			Optional:    !isDataSourceParamRequired(p, listPathParams),
		}
	}
}

func isDataSourceParamRequired(p ParamDef, listPathParams map[string]bool) bool {
	if !p.Required || p.In != "path" {
		return false
	}
	if len(listPathParams) == 0 {
		return true
	}
	return listPathParams[p.Name]
}

func addDataSourceListParams(attrs map[string]schema.Attribute, listOp *OperationDef, paramSeen map[string]bool) {
	if listOp == nil {
		return
	}
	for _, p := range listOp.Params {
		if paramSeen[p.Name] {
			continue
		}
		paramSeen[p.Name] = true
		attrName := ParamAttrName(p.Name)
		if _, exists := attrs[attrName]; exists {
			continue
		}
		attrs[attrName] = schema.StringAttribute{
			Description: fmt.Sprintf("%s parameter (%s)", p.Name, p.In),
			Optional:    true,
		}
	}
}

func addDataSourceResponseFields(attrs map[string]schema.Attribute, fields []BodyFieldDef) {
	for _, rf := range fields {
		key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
		if key == "id" || key == "api_response" {
			continue
		}
		if _, exists := attrs[key]; exists {
			continue
		}
		attrs[key] = dataSourceResponseAttr(rf)
	}
}

func dataSourceResponseAttr(rf BodyFieldDef) schema.Attribute {
	desc := rf.Desc
	if desc == "" {
		desc = rf.Path
	}
	if rf.IsArray && len(rf.ItemFields) > 0 {
		return schema.ListNestedAttribute{
			Description: desc,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: buildDSNestedItemAttrs(rf.ItemFields),
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
	return schema.StringAttribute{
		Description: desc,
		Computed:    true,
	}
}

func readDataSourceParams(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, op *OperationDef) (map[string]string, map[string]string) {
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	for _, p := range op.Params {
		val, ok := getDataSourceParamValue(ctx, req, resp, op, p)
		if !ok {
			continue
		}
		assignDataSourceParam(pathParams, queryParams, p, val)
	}
	return pathParams, queryParams
}

func getDataSourceParamValue(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, op *OperationDef, p ParamDef) (string, bool) {
	attrName := ParamAttrName(p.Name)
	var val types.String
	dd := req.Config.GetAttribute(ctx, attrPath(attrName), &val)
	resp.Diagnostics.Append(dd...)
	if dd.HasError() {
		return "", false
	}
	if val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
		if p.Required {
			resp.Diagnostics.AddError("Missing Required Parameter",
				fmt.Sprintf("Parameter %q is required for operation %s", p.Name, op.OperationID))
		}
		return "", false
	}
	return val.ValueString(), true
}

func assignDataSourceParam(pathParams, queryParams map[string]string, p ParamDef, val string) {
	switch p.In {
	case "path":
		pathParams[p.Name] = val
	case "query":
		queryParams[p.Name] = val
	}
}

func setDataSourceResult(ctx context.Context, resp *datasource.ReadResponse, op *OperationDef, pathParams map[string]string, result any) {
	if result == nil {
		return
	}
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("Response Error", fmt.Sprintf("Failed to marshal response: %v", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("api_response"), types.StringValue(string(jsonBytes)))...)

	if m, ok := result.(map[string]any); ok {
		if setDataSourceMapResult(ctx, resp, op, m) {
			return
		}
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("id"), types.StringValue(buildListID(op.OperationID, pathParams)))...)
}

func setDataSourceMapResult(ctx context.Context, resp *datasource.ReadResponse, op *OperationDef, m map[string]any) bool {
	id := extractID(m)
	if id != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("id"), types.StringValue(id))...)
	}
	for _, rf := range op.ResponseFields {
		setDataSourceResponseField(ctx, resp, rf, m)
	}
	return id != ""
}

func setDataSourceResponseField(ctx context.Context, resp *datasource.ReadResponse, rf BodyFieldDef, m map[string]any) {
	key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
	if key == "id" || key == "api_response" {
		return
	}
	val, ok := handlers.GetNested(m, rf.Path)
	if !ok || val == nil {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath(key), dataSourceResponseValue(val, rf))...)
}

func dataSourceResponseValue(val any, rf BodyFieldDef) any {
	if rf.IsArray && len(rf.ItemFields) > 0 {
		if arr, ok := val.([]any); ok {
			return buildListFromResponse(arr, rf.ItemFields)
		}
		return nil
	}
	if rf.IsArray {
		if arr, ok := val.([]any); ok {
			return buildSimpleListFromResponse(arr)
		}
		return nil
	}
	return types.StringValue(stringifyResponseValue(val))
}

func stringifyResponseValue(val any) string {
	switch val.(type) {
	case []any, map[string]any:
		if b, err := json.Marshal(val); err == nil {
			return string(b)
		}
	}
	return fmt.Sprintf("%v", val)
}

func copyDataSourceConfigToState(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, params []ParamDef) {
	for _, p := range params {
		attrName := ParamAttrName(p.Name)
		var val types.String
		dd := req.Config.GetAttribute(ctx, attrPath(attrName), &val)
		resp.Diagnostics.Append(dd...)
		if !dd.HasError() && !val.IsNull() && !val.IsUnknown() {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath(attrName), val)...)
		}
	}
}
