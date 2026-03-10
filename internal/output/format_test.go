package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

func TestRenderJSON(t *testing.T) {
	output.Format = "json"

	id := 42
	title := "Test PR"
	state := generated.PullrequestStateOPEN
	pr := generated.Pullrequest{
		Id:    &id,
		Title: &title,
		State: &state,
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, pr); err != nil {
		t.Fatalf("RenderTo: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"id": 42`) {
		t.Errorf("expected JSON to contain id=42, got: %s", got)
	}
	if !strings.Contains(got, `"title": "Test PR"`) {
		t.Errorf("expected JSON to contain title, got: %s", got)
	}
}

func TestRenderTable_EmptySlice(t *testing.T) {
	output.Format = "table"

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, []generated.Pullrequest{}); err != nil {
		t.Fatalf("RenderTo: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "(no results)") {
		t.Errorf("expected '(no results)', got: %s", got)
	}
}

func TestRenderTable_Slice(t *testing.T) {
	output.Format = "table"

	id1, id2 := 1, 2
	t1, t2 := "PR One", "PR Two"
	s1 := generated.PullrequestStateOPEN
	s2 := generated.PullrequestStateMERGED

	prs := []generated.Pullrequest{
		{Id: &id1, Title: &t1, State: &s1},
		{Id: &id2, Title: &t2, State: &s2},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, prs); err != nil {
		t.Fatalf("RenderTo: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "PR One") {
		t.Errorf("expected 'PR One' in table output, got: %s", got)
	}
	if !strings.Contains(got, "PR Two") {
		t.Errorf("expected 'PR Two' in table output, got: %s", got)
	}
}

func TestRenderID(t *testing.T) {
	output.Format = "id"

	id1, id2 := 10, 20
	t1, t2 := "A", "B"
	prs := []generated.Pullrequest{
		{Id: &id1, Title: &t1},
		{Id: &id2, Title: &t2},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, prs); err != nil {
		t.Fatalf("RenderTo: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "10") || !strings.Contains(got, "20") {
		t.Errorf("expected IDs 10 and 20, got: %s", got)
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	output.Format = "xml"

	var buf bytes.Buffer
	err := output.RenderTo(&buf, "anything")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
