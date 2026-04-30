package tfprovider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

const testItemReadPath = "/items/{workspace}/{id}"

type mockState struct {
	values map[string]attr.Value
	set    map[string]any
	diags  map[string]diag.Diagnostics
}

func newMockState(values map[string]attr.Value) *mockState {
	return &mockState{
		values: values,
		set:    map[string]any{},
		diags:  map[string]diag.Diagnostics{},
	}
}

func (m *mockState) GetAttribute(_ context.Context, p path.Path, target interface{}) diag.Diagnostics {
	name := p.String()
	if d, ok := m.diags[name]; ok {
		return d
	}
	v, ok := m.values[name]
	if !ok {
		switch t := target.(type) {
		case *types.String:
			*t = types.StringNull()
		case *types.List:
			*t = types.ListNull(types.StringType)
		case *types.Object:
			*t = types.ObjectNull(map[string]attr.Type{})
		default:
			panic("unsupported target type")
		}
		return nil
	}

	switch t := target.(type) {
	case *types.String:
		*t = v.(types.String)
	case *types.List:
		*t = v.(types.List)
	case *types.Object:
		*t = v.(types.Object)
	default:
		panic("unsupported target type")
	}
	return nil
}

func (m *mockState) SetAttribute(_ context.Context, p path.Path, val interface{}) diag.Diagnostics {
	m.set[p.String()] = val
	return nil
}

func testBBClient(serverURL string) *client.BBClient {
	return &client.BBClient{Client: resty.New().SetBaseURL(serverURL).SetBasicAuth("u", "p")}
}

func nestedObjectValue(fields []BodyFieldDef, values map[string]attr.Value) types.Object {
	return types.ObjectValueMust(itemAttrTypes(fields), values)
}

func nestedListValue(fields []BodyFieldDef, values ...map[string]attr.Value) types.List {
	objType := types.ObjectType{AttrTypes: itemAttrTypes(fields)}
	items := make([]attr.Value, 0, len(values))
	for _, item := range values {
		items = append(items, types.ObjectValueMust(objType.AttrTypes, item))
	}
	return types.ListValueMust(objType, items)
}

func stringListValue(items ...string) types.List {
	vals := make([]attr.Value, 0, len(items))
	for _, item := range items {
		vals = append(vals, types.StringValue(item))
	}
	return types.ListValueMust(types.StringType, vals)
}

func testResourceGroup() ResourceGroup {
	itemFields := []BodyFieldDef{{Path: "name", Type: "string", Desc: "name"}}
	return ResourceGroup{
		TypeName:    "sample_group",
		Description: "Sample resource",
		Ops: CRUDOps{
			Create: &OperationDef{
				OperationID: "createSample",
				Method:      http.MethodPost,
				Path:        "/items/{workspace}",
				Params: []ParamDef{
					{Name: "workspace", In: "path", Required: true},
					{Name: "filter", In: "query"},
				},
				HasBody: true,
				BodyFields: []BodyFieldDef{
					{Path: "title", Type: "string", Desc: "title"},
					{Path: "settings", Type: "string", Desc: "settings", IsObject: true, ItemFields: itemFields},
					{Path: "reviewers", Type: "string", Desc: "reviewers", IsArray: true, ItemFields: itemFields},
					{Path: "tags", Type: "string", Desc: "tags", IsArray: true},
				},
			},
			Read: &OperationDef{
				OperationID: "getSample",
				Method:      http.MethodGet,
				Path:        testItemReadPath,
				Params: []ParamDef{
					{Name: "workspace", In: "path", Required: true},
					{Name: "id", In: "path", Required: true},
				},
				ResponseFields: []BodyFieldDef{
					{Path: "title", Type: "string", Desc: "title"},
					{Path: "settings", Type: "string", Desc: "settings", IsObject: true, ItemFields: itemFields},
					{Path: "reviewers", Type: "string", Desc: "reviewers", IsArray: true, ItemFields: itemFields},
					{Path: "tags", Type: "string", Desc: "tags", IsArray: true},
					{Path: "metadata", Type: "string", Desc: "metadata"},
				},
			},
			Update: &OperationDef{
				OperationID: "updateSample",
				Method:      http.MethodPut,
				Path:        testItemReadPath,
				Params: []ParamDef{
					{Name: "workspace", In: "path", Required: true},
					{Name: "id", In: "path", Required: true},
				},
				HasBody: true,
				BodyFields: []BodyFieldDef{
					{Path: "title", Type: "string", Desc: "title"},
				},
			},
			Delete: &OperationDef{
				OperationID: "deleteSample",
				Method:      http.MethodDelete,
				Path:        testItemReadPath,
				Params: []ParamDef{
					{Name: "workspace", In: "path", Required: true},
					{Name: "id", In: "path", Required: true},
				},
			},
			List: &OperationDef{
				OperationID: "listSamples",
				Method:      http.MethodGet,
				Path:        "/items/{workspace}",
				Params: []ParamDef{
					{Name: "workspace", In: "path", Required: true},
					{Name: "state", In: "query"},
				},
				Paginated: true,
			},
		},
	}
}

func TestGenericResourceSchema(t *testing.T) {
	r := &GenericResource{group: testResourceGroup()}

	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	attrs := resp.Schema.Attributes
	workspace := attrs["workspace"].(resourceschema.StringAttribute)
	if !workspace.Required {
		t.Fatal("workspace should be required")
	}

	paramID := attrs["param_id"].(resourceschema.StringAttribute)
	if !paramID.Optional || !paramID.Computed {
		t.Fatalf("param_id should be optional+computed: %#v", paramID)
	}

	title := attrs["title"].(resourceschema.StringAttribute)
	if !title.Optional || !title.Computed {
		t.Fatalf("title should be optional+computed: %#v", title)
	}

	settings := attrs["settings"].(resourceschema.SingleNestedAttribute)
	if !settings.Optional || !settings.Computed {
		t.Fatalf("settings should be optional+computed: %#v", settings)
	}

	reviewers := attrs["reviewers"].(resourceschema.ListNestedAttribute)
	if !reviewers.Optional || !reviewers.Computed {
		t.Fatalf("reviewers should be optional+computed: %#v", reviewers)
	}

	tags := attrs["tags"].(resourceschema.ListAttribute)
	if !tags.Optional || !tags.Computed {
		t.Fatalf("tags should be optional+computed: %#v", tags)
	}

	metadata := attrs["metadata"].(resourceschema.StringAttribute)
	if !metadata.Computed {
		t.Fatalf("metadata should be computed: %#v", metadata)
	}

	requestBody := attrs["request_body"].(resourceschema.StringAttribute)
	if !requestBody.Optional {
		t.Fatalf("request_body should be optional: %#v", requestBody)
	}
}

func TestGenericResourceConfigureAndWrappers(t *testing.T) {
	r := &GenericResource{group: ResourceGroup{TypeName: "sample"}}

	var cfgResp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{}, &cfgResp)
	if cfgResp.Diagnostics.HasError() {
		t.Fatal("expected nil provider data to be ignored")
	}

	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "wrong"}, &cfgResp)
	if !cfgResp.Diagnostics.HasError() {
		t.Fatal("expected wrong provider data type error")
	}

	var createResp resource.CreateResponse
	r.Create(context.Background(), resource.CreateRequest{}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected create unsupported error")
	}

	var updateResp resource.UpdateResponse
	r.Update(context.Background(), resource.UpdateRequest{}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected update unsupported error")
	}

	var readResp resource.ReadResponse
	r.Read(context.Background(), resource.ReadRequest{}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatal("expected read without op to be ignored")
	}
}

func TestResourceExtractParamsAndBody(t *testing.T) {
	group := testResourceGroup()
	r := &GenericResource{group: group}
	ctx := context.Background()

	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	source := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"filter":    types.StringValue("open"),
		"title":     types.StringValue("Hello"),
		"settings":  nestedObjectValue(itemFields, map[string]attr.Value{"name": types.StringValue("cfg")}),
		"reviewers": nestedListValue(itemFields, map[string]attr.Value{"name": types.StringValue("alice")}),
		"tags":      stringListValue("one", "two"),
	})
	var diags diag.Diagnostics

	pathParams, queryParams := r.extractParams(ctx, group.Ops.Create, source, nil, &diags)
	if !reflect.DeepEqual(pathParams, map[string]string{"workspace": "ws"}) {
		t.Fatalf("unexpected path params: %#v", pathParams)
	}
	if !reflect.DeepEqual(queryParams, map[string]string{"filter": "open"}) {
		t.Fatalf("unexpected query params: %#v", queryParams)
	}

	body := r.buildBody(ctx, group.Ops.Create, source, &diags)
	var got map[string]any
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("buildBody returned invalid JSON: %v", err)
	}
	if got["title"] != "Hello" {
		t.Fatalf("expected title in body, got %#v", got)
	}
	if !reflect.DeepEqual(got["tags"], []any{"one", "two"}) {
		t.Fatalf("expected tags in body, got %#v", got["tags"])
	}

	rawSource := newMockState(map[string]attr.Value{
		"request_body": types.StringValue(`{"custom":true}`),
	})
	rawOp := &OperationDef{OperationID: "rawBody", HasBody: true}
	if raw := r.buildBody(ctx, rawOp, rawSource, &diags); raw != `{"custom":true}` {
		t.Fatalf("expected raw request_body fallback, got %q", raw)
	}
}

// TestResourceExtractParamsFallbackToState reproduces the bug reported in
// "bitbucket_branch_restrictions Update fails with Missing Required Parameter id".
// When an in-place Update is planned, the Computed-only path param (e.g. the
// numeric "id"/"param_id") is "(known after apply)" in the plan. Without a
// fallback to prior state, extractParams reports it as missing. With the
// fallback, the value is sourced from state and the dispatch can proceed.
func TestResourceExtractParamsFallbackToState(t *testing.T) {
	group := testResourceGroup()
	r := &GenericResource{group: group}
	ctx := context.Background()

	// Plan: workspace is known, param_id is unknown (typical for in-place Update
	// of a resource whose id is Computed-only).
	plan := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringUnknown(),
	})
	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("75417602"),
	})

	// Without fallback, the missing/unknown id should produce a diagnostic.
	var diagsNoFallback diag.Diagnostics
	r.extractParams(ctx, group.Ops.Update, plan, nil, &diagsNoFallback)
	if !diagsNoFallback.HasError() {
		t.Fatal("expected Missing Required Parameter without state fallback")
	}

	// With state as fallback, the id is recovered and no error is reported.
	var diagsWithFallback diag.Diagnostics
	pathParams, _ := r.extractParams(ctx, group.Ops.Update, plan, state, &diagsWithFallback)
	if diagsWithFallback.HasError() {
		t.Fatalf("unexpected diagnostics with state fallback: %#v", diagsWithFallback)
	}
	if got := pathParams["id"]; got != "75417602" {
		t.Fatalf("expected id %q from state fallback, got %q", "75417602", got)
	}
	if got := pathParams["workspace"]; got != "ws" {
		t.Fatalf("expected workspace %q from plan, got %q", "ws", got)
	}

	// When both plan and state have a known value, the plan takes precedence
	// (the fallback is only consulted when the plan value is null/unknown/empty).
	planWithID := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("plan-id"),
	})
	stateWithID := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("state-id"),
	})
	var diagsPrecedence diag.Diagnostics
	pathParams, _ = r.extractParams(ctx, group.Ops.Update, planWithID, stateWithID, &diagsPrecedence)
	if diagsPrecedence.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diagsPrecedence)
	}
	if got := pathParams["id"]; got != "plan-id" {
		t.Fatalf("expected plan to take precedence over state fallback, got %q", got)
	}
}

func TestResourceHelperSelectionAndCopyAttributes(t *testing.T) {
	group := testResourceGroup()
	r := &GenericResource{group: group}
	ctx := context.Background()
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	source := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"title":     types.StringValue("Hello"),
		"settings":  nestedObjectValue(itemFields, map[string]attr.Value{"name": types.StringValue("cfg")}),
		"reviewers": nestedListValue(itemFields, map[string]attr.Value{"name": types.StringValue("alice")}),
		"tags":      stringListValue("one", "two"),
	})
	target := newMockState(nil)
	var diags diag.Diagnostics

	missingSource := newMockState(map[string]attr.Value{})
	diags = nil
	r.extractParams(ctx, group.Ops.Read, missingSource, nil, &diags)
	if !diags.HasError() {
		t.Fatal("expected missing required parameter diagnostics")
	}

	if len(r.crudOps()) != 5 {
		t.Fatalf("expected 5 CRUD ops, got %d", len(r.crudOps()))
	}
	if len(r.responseFields()) != len(group.Ops.Read.ResponseFields) {
		t.Fatalf("unexpected response fields: %#v", r.responseFields())
	}

	r.copyAttributes(ctx, group.Ops.Create, source, target, &diags)
	if _, ok := target.set["reviewers"]; ok {
		t.Fatal("copyAttributes should skip list-nested body fields")
	}
	if _, ok := target.set["settings"]; ok {
		t.Fatal("copyAttributes should skip object body fields")
	}
	if target.set["title"] != types.StringValue("Hello") {
		t.Fatalf("expected title copied to target, got %#v", target.set["title"])
	}
}

func TestResourcePopulateComputedParamsAndExtractResponseFields(t *testing.T) {
	group := testResourceGroup()
	r := &GenericResource{group: group}
	ctx := context.Background()
	var diags diag.Diagnostics
	computedSource := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringNull(),
	})
	computedTarget := newMockState(nil)
	r.populateComputedParams(ctx, map[string]any{"id": 42}, computedSource, computedTarget, &diags)
	if got := computedTarget.set["param_id"]; got != types.StringValue("42") {
		t.Fatalf("expected computed param_id, got %#v", got)
	}

	responseTarget := newMockState(nil)
	diags = nil
	r.extractResponseFields(ctx, map[string]any{
		"title":    "Hello",
		"settings": map[string]any{"name": "cfg"},
		"reviewers": []any{
			map[string]any{"name": "alice"},
		},
		"tags":     []any{"one", "two"},
		"metadata": map[string]any{"mode": "full"},
	}, responseTarget, &diags)
	if got := responseTarget.set["title"]; got != types.StringValue("Hello") {
		t.Fatalf("expected title response field, got %#v", got)
	}
	if _, ok := responseTarget.set["settings"].(types.Object); !ok {
		t.Fatalf("expected settings object value, got %T", responseTarget.set["settings"])
	}
	if _, ok := responseTarget.set["reviewers"].(setLikeListValue); !ok {
		t.Fatalf("expected reviewers setLikeListValue, got %T", responseTarget.set["reviewers"])
	}
	if _, ok := responseTarget.set["tags"].(types.List); !ok {
		t.Fatalf("expected tags list value, got %T", responseTarget.set["tags"])
	}
	if got := responseTarget.set["metadata"]; got != types.StringValue(`{"mode":"full"}`) {
		t.Fatalf("expected metadata JSON string, got %#v", got)
	}
}

func TestResourceReadValueHelpers(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}

	if got := readSimpleListValue(stringListValue("one", "two")); !reflect.DeepEqual(got, []string{"one", "two"}) {
		t.Fatalf("unexpected simple list: %#v", got)
	}
	if got := readSimpleListValue(types.ListValueMust(types.StringType, []attr.Value{})); got != nil {
		t.Fatalf("expected nil for empty list, got %#v", got)
	}

	listWithSparse := types.ListValueMust(types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}, []attr.Value{
		types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{"name": types.StringNull()}),
		types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{"name": types.StringValue("alice")}),
	})
	if got := readListNestedValue(listWithSparse, itemFields); !reflect.DeepEqual(got, []map[string]any{{"name": "alice"}}) {
		t.Fatalf("unexpected nested list: %#v", got)
	}

	if got := readSingleNested(context.Background(), newMockState(map[string]attr.Value{
		"settings": types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{"name": types.StringNull()}),
	}), "settings", itemFields, &diag.Diagnostics{}); got != nil {
		t.Fatalf("expected nil for empty single nested object, got %#v", got)
	}

	if got, ok := readAttrValue(types.StringUnknown(), BodyFieldDef{Path: "title", Type: "string"}); ok || got != nil {
		t.Fatalf("expected unknown string to be skipped, got %#v ok=%v", got, ok)
	}
	if got, ok := readAttrValue(types.ListNull(types.StringType), BodyFieldDef{Path: "tags", Type: "string", IsArray: true}); ok || got != nil {
		t.Fatalf("expected null list to be skipped, got %#v ok=%v", got, ok)
	}
}

func TestResourceResponseValueHelpers(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	if got := buildSimpleListFromResponse([]any{"one", 2}); len(got.Elements()) != 2 {
		t.Fatalf("expected two simple response elements, got %#v", got)
	}
	if got := buildSimpleListFromResponse([]any{}); len(got.Elements()) != 0 {
		t.Fatalf("expected empty simple response list, got %#v", got)
	}

	obj := buildObjectFromResponse(map[string]any{"name": "alice"}, itemFields)
	if name := obj.Attributes()["name"].(types.String).ValueString(); name != "alice" {
		t.Fatalf("expected object name, got %q", name)
	}

	if _, ok := buildAttrValueFromResponse("bad", BodyFieldDef{Path: "settings", IsObject: true, ItemFields: itemFields}).(types.Object); !ok {
		t.Fatal("expected null object attr value for invalid nested object")
	}
	if _, ok := buildAttrValueFromResponse("bad", BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: itemFields}).(setLikeListValue); !ok {
		t.Fatal("expected null setLikeList attr value for invalid nested list")
	}
	if _, ok := buildAttrValueFromResponse("bad", BodyFieldDef{Path: "tags", IsArray: true}).(types.List); !ok {
		t.Fatal("expected null simple list attr value for invalid list")
	}
	if got := buildAttrValueFromResponse(99, BodyFieldDef{Path: "id"}).(types.String).ValueString(); got != "99" {
		t.Fatalf("expected scalar attr string conversion, got %q", got)
	}

	if got := attrNullValue(BodyFieldDef{Path: "settings", IsObject: true, ItemFields: itemFields}); !got.IsNull() {
		t.Fatal("expected null object attr")
	}
	if got := attrNullValue(BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: itemFields}); !got.IsNull() {
		t.Fatal("expected null nested list attr")
	}
	if got := attrNullValue(BodyFieldDef{Path: "tags", IsArray: true}); !got.IsNull() {
		t.Fatal("expected null simple list attr")
	}
}

func TestResourceIDHelpers(t *testing.T) {
	if id := extractID(map[string]any{"uuid": "u-1"}); id != "u-1" {
		t.Fatalf("expected uuid id, got %q", id)
	}
	if id := extractID(map[string]any{"slug": "repo"}); id != "repo" {
		t.Fatalf("expected slug id, got %q", id)
	}
	if val, ok := responseParamValue(map[string]any{"issue": "7"}, "issue_id"); !ok || val != "7" {
		t.Fatalf("expected issue_id fallback, got %q ok=%v", val, ok)
	}
	if val, ok := responseParamValue(map[string]any{"name": "demo"}, "missing"); !ok || val != "demo" {
		t.Fatalf("expected extractID fallback, got %q ok=%v", val, ok)
	}
	// Large integer IDs from JSON (decoded as float64) must not use scientific notation.
	if id := extractID(map[string]any{"id": float64(74332764)}); id != "74332764" {
		t.Fatalf("expected integer id without scientific notation, got %q", id)
	}
	if val, ok := responseParamValue(map[string]any{"id": float64(74332764)}, "id"); !ok || val != "74332764" {
		t.Fatalf("expected integer responseParamValue without scientific notation, got %q ok=%v", val, ok)
	}
}

func TestGenericResourceDispatch(t *testing.T) {
	group := testResourceGroup()
	ctx := context.Background()
	source := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"title":     types.StringValue("Hello"),
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/items/ws":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 5, "title": "Hello"})
		case r.Method == http.MethodDelete && r.URL.Path == "/items/ws/5":
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	r := &GenericResource{group: group, client: testBBClient(srv.URL)}
	target := newMockState(nil)
	var diags diag.Diagnostics

	r.dispatch(ctx, group.Ops.Create, source, target, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected dispatch diagnostics: %#v", diags)
	}
	if got := target.set["id"]; got != types.StringValue("5") {
		t.Fatalf("expected response ID, got %#v", got)
	}
	if got := target.set["api_response"].(types.String).ValueString(); !strings.Contains(got, `"id": 5`) {
		t.Fatalf("expected api_response JSON, got %q", got)
	}

	deleteSource := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
	})
	deleteTarget := newMockState(nil)
	diags = nil
	r.dispatch(ctx, group.Ops.Delete, deleteSource, deleteTarget, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected delete dispatch diagnostics: %#v", diags)
	}
	if got := deleteTarget.set["id"]; got != types.StringValue(group.Ops.Delete.OperationID) {
		t.Fatalf("expected fallback delete ID, got %#v", got)
	}
	if got := deleteTarget.set["api_response"]; got != types.StringValue("") {
		t.Fatalf("expected empty api_response on delete, got %#v", got)
	}
}

func TestResourceSchemaHelpersAndFallbacks(t *testing.T) {
	attrs := resourceBaseAttrs()
	if _, ok := attrs["id"].(resourceschema.StringAttribute); !ok {
		t.Fatalf("expected id base attribute, got %T", attrs["id"])
	}

	readOnly := &GenericResource{group: ResourceGroup{Ops: CRUDOps{Read: testResourceGroup().Ops.Read}}}
	if readOnly.primaryOp() != readOnly.group.Ops.Read {
		t.Fatal("expected primaryOp to fall back to read")
	}
	if readOnly.responseOp() != readOnly.group.Ops.Read {
		t.Fatal("expected responseOp to prefer read")
	}
	if required := requiredPrimaryPathParams(nil); len(required) != 0 {
		t.Fatalf("expected nil primary op to produce no required params, got %#v", required)
	}

	queryAttr := resourceParamAttr(ParamDef{Name: "state", In: "query"}, false)
	if !queryAttr.Optional || queryAttr.Computed {
		t.Fatalf("expected optional non-computed query attr, got %#v", queryAttr)
	}
	requiredAttr := resourceParamAttr(ParamDef{Name: "workspace", In: "path"}, true)
	if !requiredAttr.Required {
		t.Fatalf("expected required attr, got %#v", requiredAttr)
	}

	if !canMergeComputedAttr(false, false) || canMergeComputedAttr(true, false) || canMergeComputedAttr(false, true) {
		t.Fatal("unexpected canMergeComputedAttr result")
	}
	if got := fieldDescription(BodyFieldDef{Path: "title"}); got != "title" {
		t.Fatalf("expected path fallback description, got %q", got)
	}

	requestBodyAttrs := map[string]resourceschema.Attribute{}
	addRequestBodyAttr(requestBodyAttrs, false)
	if len(requestBodyAttrs) != 0 {
		t.Fatalf("expected no request_body attr when no body exists, got %#v", requestBodyAttrs)
	}
}

func TestResourceNestedSchemaHelpers(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	nested := buildNestedItemAttrs([]BodyFieldDef{
		{Path: "settings", IsObject: true, ItemFields: itemFields},
		{Path: "reviewers", IsArray: true, ItemFields: itemFields},
		{Path: "tags", IsArray: true},
		{Path: "title"},
	})
	if _, ok := nested["settings"].(resourceschema.SingleNestedAttribute); !ok {
		t.Fatalf("expected nested object attr, got %T", nested["settings"])
	}
	if _, ok := nested["reviewers"].(resourceschema.ListNestedAttribute); !ok {
		t.Fatalf("expected nested list attr, got %T", nested["reviewers"])
	}
	if _, ok := nested["tags"].(resourceschema.ListAttribute); !ok {
		t.Fatalf("expected list attr, got %T", nested["tags"])
	}

	attrTypes := itemAttrTypes([]BodyFieldDef{
		{Path: "settings", IsObject: true, ItemFields: itemFields},
		{Path: "reviewers", IsArray: true, ItemFields: itemFields},
		{Path: "tags", IsArray: true},
		{Path: "title"},
	})
	if _, ok := attrTypes["settings"].(types.ObjectType); !ok {
		t.Fatalf("expected object attr type, got %T", attrTypes["settings"])
	}
	if _, ok := attrTypes["reviewers"].(setLikeListType); !ok {
		t.Fatalf("expected setLikeListType for nested-object array, got %T", attrTypes["reviewers"])
	}
}

func TestResourceResponseAttributeHelpers(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	if _, ok := responseFieldAttr(BodyFieldDef{Path: "settings", IsObject: true, ItemFields: itemFields}).(resourceschema.SingleNestedAttribute); !ok {
		t.Fatal("expected single nested response attribute")
	}
	if _, ok := responseFieldAttr(BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute); !ok {
		t.Fatal("expected list nested response attribute")
	}
	if _, ok := responseFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}).(resourceschema.ListAttribute); !ok {
		t.Fatal("expected list response attribute")
	}

	merged := mergeResponseAttr(resourceschema.StringAttribute{Optional: true}, BodyFieldDef{Path: "title"}).(resourceschema.StringAttribute)
	if !merged.Computed || merged.Description != "title" {
		t.Fatalf("expected merged string attr, got %#v", merged)
	}
	required := mergeResponseAttr(resourceschema.StringAttribute{Required: true}, BodyFieldDef{Path: "title"}).(resourceschema.StringAttribute)
	if required.Computed {
		t.Fatalf("expected required string attr to remain non-computed, got %#v", required)
	}
	objectAttr := mergeResponseAttr(resourceschema.SingleNestedAttribute{Optional: true}, BodyFieldDef{Path: "settings", IsObject: true, ItemFields: itemFields}).(resourceschema.SingleNestedAttribute)
	if !objectAttr.Computed || len(objectAttr.Attributes) == 0 {
		t.Fatalf("expected merged object attr, got %#v", objectAttr)
	}
	listAttr := mergeResponseAttr(resourceschema.ListNestedAttribute{Optional: true}, BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute)
	if !listAttr.Computed || len(listAttr.NestedObject.Attributes) == 0 {
		t.Fatalf("expected merged list nested attr, got %#v", listAttr)
	}
}

func TestResourceBodyAndComputedParamHelpers(t *testing.T) {
	ctx := context.Background()
	itemFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	source := newMockState(map[string]attr.Value{
		"settings": types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{"name": types.StringValue("cfg")}),
		"reviewers": types.ListValueMust(types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}, []attr.Value{
			types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{"name": types.StringValue("alice")}),
		}),
		"tags":  stringListValue("one"),
		"title": types.StringValue("Hello"),
	})
	var diags diag.Diagnostics

	if val, ok := bodyFieldValue(ctx, source, BodyFieldDef{Path: "settings", IsObject: true, ItemFields: itemFields}, &diags); !ok || val == nil {
		t.Fatalf("expected object body field value, got %#v ok=%v", val, ok)
	}
	if val, ok := bodyFieldValue(ctx, source, BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: itemFields}, &diags); !ok || val == nil {
		t.Fatalf("expected nested list body field value, got %#v ok=%v", val, ok)
	}
	if val, ok := bodyFieldValue(ctx, source, BodyFieldDef{Path: "tags", IsArray: true}, &diags); !ok || val == nil {
		t.Fatalf("expected list body field value, got %#v ok=%v", val, ok)
	}
	if val, ok := bodyFieldValue(ctx, source, BodyFieldDef{Path: "title"}, &diags); !ok || val != "Hello" {
		t.Fatalf("expected string body field value, got %#v ok=%v", val, ok)
	}

	if body := marshalBodyObject(&diags, nil); body != "" {
		t.Fatalf("expected empty marshaled body, got %q", body)
	}
	if body := rawRequestBody(ctx, newMockState(map[string]attr.Value{}), &diags); body != "" {
		t.Fatalf("expected missing request body to be empty, got %q", body)
	}

	seen := map[string]bool{}
	if !shouldSkipComputedParam(ParamDef{Name: "workspace", In: "query"}, "workspace", seen, ctx, source, &diags) {
		t.Fatal("expected non-path computed param to be skipped")
	}
	if !shouldSkipComputedParam(ParamDef{Name: "title", In: "path"}, "title", map[string]bool{}, ctx, source, &diags) {
		t.Fatal("expected existing path value to be skipped")
	}

	target := newMockState(nil)
	setComputedParam(ctx, map[string]any{"id": 7}, "missing_id", "param_id", target, &diags)
	if got := target.set["param_id"]; got != types.StringValue("7") {
		t.Fatalf("expected fallback computed param, got %#v", got)
	}
}

func TestResourceDispatchResultAndResponseHelpers(t *testing.T) {
	r := &GenericResource{group: testResourceGroup()}
	ctx := context.Background()
	target := newMockState(nil)
	var diags diag.Diagnostics

	if result := r.storeDispatchResult(ctx, &OperationDef{OperationID: "delete"}, nil, target, &diags, nil); result != nil {
		t.Fatalf("expected nil result map, got %#v", result)
	}
	if got := target.set["id"]; got != types.StringValue("delete") {
		t.Fatalf("expected fallback nil-result id, got %#v", got)
	}

	target = newMockState(nil)
	diags = nil
	if result := r.storeDispatchResult(ctx, &OperationDef{OperationID: "bad"}, nil, target, &diags, map[string]any{"bad": make(chan int)}); result != nil || !diags.HasError() {
		t.Fatalf("expected marshal failure diagnostics, got result=%#v diags=%#v", result, diags)
	}

	if got := fallbackResourceID(&OperationDef{OperationID: "get", Params: []ParamDef{{Name: "workspace", In: "path"}, {Name: "id", In: "path"}}}, map[string]string{"workspace": "ws", "id": "5"}); got != "get/ws/5" {
		t.Fatalf("unexpected fallback resource id %q", got)
	}

	reservedTarget := newMockState(nil)
	setResponseField(ctx, BodyFieldDef{Path: "id"}, map[string]any{"id": 1}, reservedTarget, &diags)
	if len(reservedTarget.set) != 0 {
		t.Fatalf("expected reserved field to be skipped, got %#v", reservedTarget.set)
	}

	if val, ok := responseFieldValue("bad", BodyFieldDef{Path: "settings", IsObject: true, ItemFields: []BodyFieldDef{{Path: "name"}}}); ok || val != nil {
		t.Fatalf("expected invalid object response field to be skipped, got %#v ok=%v", val, ok)
	}
	if val, ok := responseFieldValue("bad", BodyFieldDef{Path: "reviewers", IsArray: true, ItemFields: []BodyFieldDef{{Path: "name"}}}); ok || val != nil {
		t.Fatalf("expected invalid nested list response field to be skipped, got %#v ok=%v", val, ok)
	}
	if val, ok := responseFieldValue("bad", BodyFieldDef{Path: "tags", IsArray: true}); ok || val != nil {
		t.Fatalf("expected invalid list response field to be skipped, got %#v ok=%v", val, ok)
	}
}
