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
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Resource identifier.",
			Computed:    true,
		},
		"api_response": schema.StringAttribute{
			Description: "The raw JSON response from the Bitbucket API.",
			Computed:    true,
		},
	}

	// Determine which operation to use for the data source schema.
	op := d.readOp()
	if op == nil {
		resp.Schema = schema.Schema{
			Description: d.group.Description,
			Attributes:  attrs,
		}
		return
	}

	// Add parameters as attributes.
	paramSeen := map[string]bool{}
	for _, p := range op.Params {
		if paramSeen[p.Name] {
			continue
		}
		paramSeen[p.Name] = true
		attrName := ParamAttrName(p.Name)
		// Skip if already defined.
		if _, exists := attrs[attrName]; exists {
			continue
		}
		attrs[attrName] = schema.StringAttribute{
			Description: fmt.Sprintf("%s parameter (%s)", p.Name, p.In),
			Required:    p.Required && p.In == "path",
			Optional:    !p.Required || p.In != "path",
		}
	}

	// Also add query params from list operation.
	if d.group.Ops.List != nil {
		for _, p := range d.group.Ops.List.Params {
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

	// Add computed attributes from response fields.
	if op != nil {
		for _, rf := range op.ResponseFields {
			key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
			// Skip reserved attributes and already-defined params.
			if key == "id" || key == "api_response" {
				continue
			}
			if _, exists := attrs[key]; exists {
				continue
			}
			desc := rf.Desc
			if desc == "" {
				desc = rf.Path
			}
			attrs[key] = schema.StringAttribute{
				Description: desc,
				Computed:    true,
			}
		}
	}

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
	op := d.readOp()
	if op == nil {
		resp.Diagnostics.AddError("Read not supported",
			fmt.Sprintf("Data source %s has no read operation", d.group.TypeName))
		return
	}

	pathParams := map[string]string{}
	queryParams := map[string]string{}

	for _, p := range op.Params {
		attrName := ParamAttrName(p.Name)
		var val types.String
		dd := req.Config.GetAttribute(ctx, attrPath(attrName), &val)
		resp.Diagnostics.Append(dd...)
		if dd.HasError() {
			continue
		}
		if val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
			if p.Required {
				resp.Diagnostics.AddError("Missing Required Parameter",
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

	if result != nil {
		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			resp.Diagnostics.AddError("Response Error", fmt.Sprintf("Failed to marshal response: %v", err))
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("api_response"), types.StringValue(string(jsonBytes)))...)

		idSet := false
		if m, ok := result.(map[string]any); ok {
			if id := extractID(m); id != "" {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("id"), types.StringValue(id))...)
				idSet = true
			}

			// Extract response fields from the API response.
			for _, rf := range op.ResponseFields {
				key := toSnakeCase(strings.ReplaceAll(rf.Path, ".", "_"))
				if key == "id" || key == "api_response" {
					continue
				}
				val, ok := handlers.GetNested(m, rf.Path)
				if !ok || val == nil {
					continue
				}
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath(key), types.StringValue(fmt.Sprintf("%v", val)))...)
			}
		}
		if !idSet {
			// For list results or missing ID fields, use a composite ID.
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath("id"), types.StringValue(buildListID(op.OperationID, pathParams)))...)
		}
	}

	// Copy config params to state.
	for _, p := range op.Params {
		attrName := ParamAttrName(p.Name)
		var val types.String
		dd := req.Config.GetAttribute(ctx, attrPath(attrName), &val)
		resp.Diagnostics.Append(dd...)
		if !dd.HasError() && !val.IsNull() && !val.IsUnknown() {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath(attrName), val)...)
		}
	}
}

// readOp returns the best operation for reading data (Read or List).
func (d *GenericDataSource) readOp() *OperationDef {
	if d.group.Ops.Read != nil {
		return d.group.Ops.Read
	}
	return d.group.Ops.List
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
