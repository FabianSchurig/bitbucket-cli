package spec

import "testing"

func TestResolveRefPropertyAddsIDFallback(t *testing.T) {
	schemas := map[string]any{
		"root": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"parent": map[string]any{
					"$ref": "#/components/schemas/comment_ref",
				},
			},
		},
		"comment_ref": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{
					"type":        "integer",
					"description": "Referenced comment ID",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "Referenced comment content",
				},
			},
		},
	}

	fields := ResolveBodyFields(schemas, "#/components/schemas/root", "", map[string]bool{})
	if len(fields) != 1 {
		t.Fatalf("field count = %d, want 1", len(fields))
	}
	parent := fields[0]
	if !parent.IsObject {
		t.Fatalf("parent.IsObject = false, want true")
	}
	if len(parent.ItemFields) != 2 {
		t.Fatalf("nested field count = %d, want 2", len(parent.ItemFields))
	}

	hasContent := false
	hasID := false
	for _, field := range parent.ItemFields {
		switch field.Path {
		case "content":
			hasContent = true
		case "id":
			hasID = true
			if field.GoType != "int" {
				t.Fatalf("id field type = %s, want int", field.GoType)
			}
		}
	}

	if !hasContent {
		t.Fatal("expected content field to be present")
	}
	if !hasID {
		t.Fatal("expected id fallback field to be present")
	}
}
