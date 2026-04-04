package tfprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func providerObjectType(attrs map[string]providerschema.Attribute) tftypes.Object {
	attrTypes := make(map[string]tftypes.Type, len(attrs))
	for name := range attrs {
		attrTypes[name] = tftypes.String
	}
	return tftypes.Object{AttributeTypes: attrTypes}
}

func datasourceObjectType(attrs map[string]datasourceschema.Attribute) tftypes.Object {
	attrTypes := make(map[string]tftypes.Type, len(attrs))
	for name, attr := range attrs {
		attrTypes[name] = datasourceAttrType(attr)
	}
	return tftypes.Object{AttributeTypes: attrTypes}
}

func datasourceAttrType(attr datasourceschema.Attribute) tftypes.Type {
	switch a := attr.(type) {
	case datasourceschema.StringAttribute:
		return tftypes.String
	case datasourceschema.ListAttribute:
		return tftypes.List{ElementType: tftypes.String}
	case datasourceschema.SingleNestedAttribute:
		return datasourceObjectType(a.Attributes)
	case datasourceschema.ListNestedAttribute:
		return tftypes.List{ElementType: datasourceObjectType(a.NestedObject.Attributes)}
	default:
		panic(fmt.Sprintf("unsupported datasource attribute type %T", attr))
	}
}

func objectValue(objectType tftypes.Object, values map[string]string) tftypes.Value {
	raw := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, typ := range objectType.AttributeTypes {
		val, ok := values[name]
		if !ok {
			raw[name] = tftypes.NewValue(typ, nil)
			continue
		}
		raw[name] = tftypes.NewValue(typ, val)
	}
	return tftypes.NewValue(objectType, raw)
}

func TestProviderRuntime(t *testing.T) {
	p := &BitbucketProvider{version: "v1.2.3"}

	var metaResp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &metaResp)
	if metaResp.TypeName != "bitbucket" || metaResp.Version != "v1.2.3" {
		t.Fatalf("unexpected metadata: %#v", metaResp)
	}

	var schemaResp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &schemaResp)
	if len(schemaResp.Schema.Attributes) != 3 {
		t.Fatalf("expected 3 provider attributes, got %d", len(schemaResp.Schema.Attributes))
	}

	t.Setenv("BITBUCKET_USERNAME", "env-user")
	t.Setenv("BITBUCKET_TOKEN", "env-token")
	t.Setenv("BITBUCKET_BASE_URL", "https://env.example")

	configType := providerObjectType(schemaResp.Schema.Attributes)
	req := provider.ConfigureRequest{
		Config: tfsdk.Config{
			Raw: objectValue(configType, map[string]string{
				"username": "cfg-user",
				"token":    "cfg-token",
				"base_url": "https://cfg.example",
			}),
			Schema: schemaResp.Schema,
		},
	}
	var cfgResp provider.ConfigureResponse
	p.Configure(context.Background(), req, &cfgResp)
	if cfgResp.Diagnostics.HasError() {
		t.Fatalf("unexpected provider configure diagnostics: %#v", cfgResp.Diagnostics)
	}
	if _, ok := cfgResp.ResourceData.(*client.BBClient); !ok {
		t.Fatalf("expected ResourceData client, got %T", cfgResp.ResourceData)
	}
	if _, ok := cfgResp.DataSourceData.(*client.BBClient); !ok {
		t.Fatalf("expected DataSourceData client, got %T", cfgResp.DataSourceData)
	}

	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_TOKEN", "")
	t.Setenv("BITBUCKET_BASE_URL", "")
	var errResp provider.ConfigureResponse
	p.Configure(context.Background(), provider.ConfigureRequest{
		Config: tfsdk.Config{
			Raw:    objectValue(configType, nil),
			Schema: schemaResp.Schema,
		},
	}, &errResp)
	if !errResp.Diagnostics.HasError() {
		t.Fatal("expected provider configure error when token is missing")
	}

	if len(p.Resources(context.Background())) == 0 {
		t.Fatal("expected registered resources")
	}
	if len(p.DataSources(context.Background())) == 0 {
		t.Fatal("expected registered data sources")
	}
	if len(RegisteredGroups()) == 0 {
		t.Fatal("expected registered groups")
	}
}

func TestDataSourceHelpers(t *testing.T) {
	group := testResourceGroup()
	d := &GenericDataSource{group: group}

	var metaResp datasource.MetadataResponse
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "bitbucket"}, &metaResp)
	if metaResp.TypeName != "bitbucket_sample_group" {
		t.Fatalf("unexpected data source type name %q", metaResp.TypeName)
	}

	var schemaResp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)
	attrs := schemaResp.Schema.Attributes
	if !attrs["workspace"].(datasourceschema.StringAttribute).Required {
		t.Fatal("workspace should be required in datasource schema")
	}
	if !attrs["param_id"].(datasourceschema.StringAttribute).Optional {
		t.Fatal("param_id should be optional in datasource schema")
	}
	if !attrs["tags"].(datasourceschema.ListAttribute).Computed {
		t.Fatal("tags should be computed list attribute")
	}
	if !attrs["reviewers"].(datasourceschema.ListNestedAttribute).Computed {
		t.Fatal("reviewers should be computed nested list attribute")
	}

	d.Configure(context.Background(), datasource.ConfigureRequest{}, &datasource.ConfigureResponse{})
	var wrongResp datasource.ConfigureResponse
	d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "wrong"}, &wrongResp)
	if !wrongResp.Diagnostics.HasError() {
		t.Fatal("expected wrong datasource configure type error")
	}

	if d.readOp().OperationID != group.Ops.Read.OperationID {
		t.Fatal("expected readOp to prefer read")
	}
	onlyList := &GenericDataSource{group: ResourceGroup{Ops: CRUDOps{List: group.Ops.List}}}
	if onlyList.readOp().OperationID != group.Ops.List.OperationID {
		t.Fatal("expected readOp fallback to list")
	}

	nested := buildDSNestedItemAttrs([]BodyFieldDef{{Path: "content.raw", Desc: "raw"}})
	if _, ok := nested["content_raw"].(datasourceschema.StringAttribute); !ok {
		t.Fatalf("expected nested datasource string attribute, got %T", nested["content_raw"])
	}

	if got := buildListID("listItems", map[string]string{"repo": "r", "workspace": "w"}); got != "listItems/r/w" {
		t.Fatalf("unexpected list ID %q", got)
	}
}

func TestDataSourceRead(t *testing.T) {
	group := testResourceGroup()
	d := &GenericDataSource{group: group}
	var schemaResp datasource.SchemaResponse
	d.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/items/ws/5":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":    5,
				"title": "one",
				"tags":  []any{"a", "b"},
				"reviewers": []any{
					map[string]any{"name": "alice"},
				},
			})
		case "/items/ws":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"values": []any{
					map[string]any{"id": 1},
					map[string]any{"id": 2},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	d.client = &client.BBClient{Client: resty.New().SetBaseURL(srv.URL).SetBasicAuth("u", "p")}

	dsType := datasourceObjectType(schemaResp.Schema.Attributes)
	makeConfig := func(values map[string]string) tfsdk.Config {
		return tfsdk.Config{
			Raw:    objectValue(dsType, values),
			Schema: schemaResp.Schema,
		}
	}
	makeState := func() tfsdk.State {
		return tfsdk.State{
			Raw:    tftypes.NewValue(dsType, nil),
			Schema: schemaResp.Schema,
		}
	}

	readReq := datasource.ReadRequest{Config: makeConfig(map[string]string{
		"workspace": "ws",
		"param_id":  "5",
	})}
	readResp := datasource.ReadResponse{State: makeState()}
	d.Read(context.Background(), readReq, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected datasource read diagnostics: %#v", readResp.Diagnostics)
	}

	var id types.String
	readResp.State.GetAttribute(context.Background(), attrPath("id"), &id)
	if id.ValueString() != "5" {
		t.Fatalf("expected datasource read id 5, got %q", id.ValueString())
	}

	var title types.String
	readResp.State.GetAttribute(context.Background(), attrPath("title"), &title)
	if title.ValueString() != "one" {
		t.Fatalf("expected datasource read title, got %q", title.ValueString())
	}

	var reviewers types.List
	readResp.State.GetAttribute(context.Background(), attrPath("reviewers"), &reviewers)
	if len(reviewers.Elements()) != 1 {
		t.Fatalf("expected nested reviewer list, got %#v", reviewers)
	}

	listReq := datasource.ReadRequest{Config: makeConfig(map[string]string{
		"workspace": "ws",
	})}
	listResp := datasource.ReadResponse{State: makeState()}
	d.Read(context.Background(), listReq, &listResp)
	if listResp.Diagnostics.HasError() {
		t.Fatalf("unexpected datasource list diagnostics: %#v", listResp.Diagnostics)
	}
	listResp.State.GetAttribute(context.Background(), attrPath("id"), &id)
	if id.ValueString() != "listSamples/ws" {
		t.Fatalf("expected list fallback id, got %q", id.ValueString())
	}
}

func TestDataSourceReadErrors(t *testing.T) {
	d := &GenericDataSource{group: ResourceGroup{TypeName: "empty"}}
	resp := datasource.ReadResponse{}
	d.Read(context.Background(), datasource.ReadRequest{}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected read not supported error")
	}
}

func TestProviderConfigUsesEnvironmentFallback(t *testing.T) {
	p := &BitbucketProvider{version: "test"}
	var schemaResp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &schemaResp)

	t.Setenv("BITBUCKET_USERNAME", "env-user")
	t.Setenv("BITBUCKET_TOKEN", "env-token")
	t.Setenv("BITBUCKET_BASE_URL", "https://env.example")

	configType := providerObjectType(schemaResp.Schema.Attributes)
	req := provider.ConfigureRequest{
		Config: tfsdk.Config{
			Raw:    objectValue(configType, nil),
			Schema: schemaResp.Schema,
		},
	}
	var resp provider.ConfigureResponse
	p.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected env-based provider configure to succeed: %#v", resp.Diagnostics)
	}
}
