package tfprovider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestReadObjectAttrsPreservesNestedArrays(t *testing.T) {
	reviewerFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	ownerFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	fields := []BodyFieldDef{
		{Path: "tags", Type: "string", IsArray: true},
		{Path: "reviewers", Type: "string", IsArray: true, ItemFields: reviewerFields},
		{
			Path:     "meta",
			Type:     "string",
			IsObject: true,
			ItemFields: []BodyFieldDef{
				{Path: "labels", Type: "string", IsArray: true},
				{Path: "owners", Type: "string", IsArray: true, ItemFields: ownerFields},
			},
		},
	}

	reviewerObjType := types.ObjectType{AttrTypes: itemAttrTypes(reviewerFields)}
	ownerObjType := types.ObjectType{AttrTypes: itemAttrTypes(ownerFields)}
	metaAttrTypes := itemAttrTypes(fields[2].ItemFields)
	root := types.ObjectValueMust(itemAttrTypes(fields), map[string]attr.Value{
		"tags": types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("one"),
			types.StringValue("two"),
		}),
		"reviewers": types.ListValueMust(reviewerObjType, []attr.Value{
			types.ObjectValueMust(reviewerObjType.AttrTypes, map[string]attr.Value{
				"name": types.StringValue("alice"),
			}),
		}),
		"meta": types.ObjectValueMust(metaAttrTypes, map[string]attr.Value{
			"labels": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("x"),
			}),
			"owners": types.ListValueMust(ownerObjType, []attr.Value{
				types.ObjectValueMust(ownerObjType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("bob"),
				}),
			}),
		}),
	})

	got := readObjectAttrs(root, fields)

	if !reflect.DeepEqual(got["tags"], []string{"one", "two"}) {
		t.Fatalf("tags mismatch: %#v", got["tags"])
	}
	if !reflect.DeepEqual(got["reviewers"], []map[string]any{{"name": "alice"}}) {
		t.Fatalf("reviewers mismatch: %#v", got["reviewers"])
	}

	meta, ok := got["meta"].(map[string]any)
	if !ok {
		t.Fatalf("meta type = %T, want map[string]any", got["meta"])
	}
	if !reflect.DeepEqual(meta["labels"], []string{"x"}) {
		t.Fatalf("meta.labels mismatch: %#v", meta["labels"])
	}
	if !reflect.DeepEqual(meta["owners"], []map[string]any{{"name": "bob"}}) {
		t.Fatalf("meta.owners mismatch: %#v", meta["owners"])
	}
}

func TestBuildListFromResponsePreservesNestedArrays(t *testing.T) {
	reviewerFields := []BodyFieldDef{{Path: "name", Type: "string"}}
	fields := []BodyFieldDef{
		{Path: "tags", Type: "string", IsArray: true},
		{Path: "reviewers", Type: "string", IsArray: true, ItemFields: reviewerFields},
	}

	list := buildListFromResponse([]any{
		map[string]any{
			"tags": []any{"one", "two"},
			"reviewers": []any{
				map[string]any{"name": "alice"},
			},
		},
	}, fields, types.ListNull(types.ObjectType{AttrTypes: itemAttrTypes(fields)}))

	elements := list.Elements()
	if len(elements) != 1 {
		t.Fatalf("element count = %d, want 1", len(elements))
	}
	obj, ok := elements[0].(types.Object)
	if !ok {
		t.Fatalf("element type = %T, want types.Object", elements[0])
	}

	tagsVal, ok := obj.Attributes()["tags"].(types.List)
	if !ok || len(tagsVal.Elements()) != 2 {
		t.Fatalf("tags value = %#v, want list with 2 elements", obj.Attributes()["tags"])
	}

	reviewersVal, ok := obj.Attributes()["reviewers"].(types.List)
	if !ok || len(reviewersVal.Elements()) != 1 {
		t.Fatalf("reviewers value = %#v, want list with 1 element", obj.Attributes()["reviewers"])
	}
	reviewerObj, ok := reviewersVal.Elements()[0].(types.Object)
	if !ok {
		t.Fatalf("reviewer element type = %T, want types.Object", reviewersVal.Elements()[0])
	}
	nameVal, ok := reviewerObj.Attributes()["name"].(types.String)
	if !ok || nameVal.ValueString() != "alice" {
		t.Fatalf("reviewer name = %#v, want alice", reviewerObj.Attributes()["name"])
	}
}
