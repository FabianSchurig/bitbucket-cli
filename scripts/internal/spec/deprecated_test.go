package spec

import "testing"

func TestBuildOperationPreservesDeprecation(t *testing.T) {
	schema := &Schema{
		Paths: map[string]PathItem{
			"/deprecated": {
				Get: &Op{
					OperationID: "deprecatedOperation",
					Deprecated:  true,
				},
			},
		},
	}

	operations := BuildOperations(schema)
	if len(operations) != 1 {
		t.Fatalf("expected one operation, got %d", len(operations))
	}
	if !operations[0].Deprecated {
		t.Fatal("expected deprecated metadata to be preserved")
	}
}
