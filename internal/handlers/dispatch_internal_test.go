package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func newInternalTestClient(serverURL string) *client.BBClient {
	r := resty.New().SetBaseURL(serverURL).SetBasicAuth("u", "p")
	return &client.BBClient{Client: r}
}

func TestDispatchRaw_EmptyAndNonJSONResponses(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantNil     bool
	}{
		{name: "empty body", contentType: "application/json", body: "", wantNil: true},
		{name: "non json body", contentType: "text/plain", body: "plain text", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			got, err := DispatchRaw(context.Background(), newInternalTestClient(srv.URL), Request{
				Method:      http.MethodGet,
				URLTemplate: "/test",
			})
			if err != nil {
				t.Fatalf("DispatchRaw returned error: %v", err)
			}
			if tt.wantNil && got != nil {
				t.Fatalf("expected nil result, got %#v", got)
			}
		})
	}
}

func TestDispatchRaw_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"broken":`))
	}))
	defer srv.Close()

	_, err := DispatchRaw(context.Background(), newInternalTestClient(srv.URL), Request{
		Method:      http.MethodGet,
		URLTemplate: "/test",
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing response") {
		t.Fatalf("expected parsing response error, got %v", err)
	}
}

func TestFetchResultAndPaginationHelpers(t *testing.T) {
	if !shouldFetchNextPage(true, &pageResult{nextURL: "https://example.test/next"}) {
		t.Fatal("expected pagination helper to request the next page")
	}
	if shouldFetchNextPage(false, &pageResult{nextURL: "https://example.test/next"}) {
		t.Fatal("expected pagination helper to stop when all=false")
	}

	_, _, err := fetchResult(context.Background(), newInternalTestClient("http://127.0.0.1:1"), Request{
		Method:      http.MethodGet,
		URLTemplate: "/test",
	}, "http://127.0.0.1:1/test", "http://127.0.0.1:1/test")
	if err == nil {
		t.Fatal("expected fetchResult request error")
	}
}

func TestExtractPage(t *testing.T) {
	t.Run("non page payloads", func(t *testing.T) {
		for _, payload := range []any{
			"not-a-map",
			map[string]any{"values": "wrong-shape"},
			map[string]any{"id": 1},
		} {
			if page := extractPage(payload); page != nil {
				t.Fatalf("expected nil page for %#v, got %#v", payload, page)
			}
		}
	})

	t.Run("page with next", func(t *testing.T) {
		page := extractPage(map[string]any{
			"values": []any{map[string]any{"id": 1}},
			"next":   "https://example.com/page2",
		})
		if page == nil {
			t.Fatal("expected page result")
		}
		if len(page.values) != 1 || page.nextURL != "https://example.com/page2" {
			t.Fatalf("unexpected page result: %#v", page)
		}
	})
}

func TestSetNestedAndGetNested(t *testing.T) {
	body := map[string]any{}

	SetNested(body, "content.raw", "hello")
	SetNested(body, "reviewers", `[{"uuid":"1"}]`)
	SetNested(body, "metadata", `{"approved":true}`)
	SetNested(body, "raw.invalid", "{not-json}")

	if got, ok := GetNested(body, "content.raw"); !ok || got != "hello" {
		t.Fatalf("expected nested raw content, got %#v ok=%v", got, ok)
	}

	reviewers, ok := GetNested(body, "reviewers")
	if !ok {
		t.Fatal("expected reviewers to be present")
	}
	reviewerSlice, ok := reviewers.([]any)
	if !ok || len(reviewerSlice) != 1 {
		t.Fatalf("expected parsed reviewers array, got %#v", reviewers)
	}

	metadata, ok := GetNested(body, "metadata")
	if !ok {
		t.Fatal("expected metadata to be present")
	}
	metadataMap, ok := metadata.(map[string]any)
	if !ok || metadataMap["approved"] != true {
		t.Fatalf("expected parsed metadata object, got %#v", metadata)
	}

	if got, ok := GetNested(body, "raw.invalid"); !ok || got != "{not-json}" {
		t.Fatalf("expected invalid JSON string to stay raw, got %#v ok=%v", got, ok)
	}

	if got, ok := GetNested(body, "content.missing"); ok || got != nil {
		t.Fatalf("expected missing nested key to return nil,false, got %#v ok=%v", got, ok)
	}

	body["notMap"] = "value"
	if got, ok := GetNested(body, "notMap.child"); ok || got != nil {
		t.Fatalf("expected non-map traversal to fail, got %#v ok=%v", got, ok)
	}
}

func TestSetNested_MergesIntoExistingMap(t *testing.T) {
	body := map[string]any{
		"content": map[string]any{"raw": "existing"},
	}

	SetNested(body, "content.html", "<p>value</p>")

	content, ok := body["content"].(map[string]any)
	if !ok {
		t.Fatalf("expected content map, got %#v", body["content"])
	}
	if content["raw"] != "existing" || content["html"] != "<p>value</p>" {
		b, _ := json.Marshal(body)
		t.Fatalf("expected merged nested content, got %s", b)
	}
}
