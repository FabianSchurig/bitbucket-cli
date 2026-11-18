package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplyBodyEnvelope(t *testing.T) {
	tests := []struct {
		name, body, envelope, want string
	}{
		{"no envelope", `{"name":"x"}`, "", `{"name":"x"}`},
		{"empty body", "", "change", ""},
		{"wraps object", `{"name":"x"}`, "change", `{"change":{"name":"x"}}`},
		{"already wrapped", `{"change":{"name":"x"}}`, "change", `{"change":{"name":"x"}}`},
		{"wraps when envelope key has siblings", `{"change":1,"name":"x"}`, "change", `{"change":{"change":1,"name":"x"}}`},
		{"invalid JSON is passed through", `not json`, "change", `not json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyBodyEnvelope(tt.body, tt.envelope); got != tt.want {
				t.Errorf("applyBodyEnvelope(%q, %q) = %q, want %q", tt.body, tt.envelope, got, tt.want)
			}
		})
	}
}

func TestDispatchRawSendsEnvelopedBody(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = string(b)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	_, err := DispatchRaw(context.Background(), newInternalTestClient(srv.URL), Request{
		Method:       http.MethodPost,
		URLTemplate:  srv.URL + "/environments/{uuid}/changes",
		PathParams:   map[string]string{"uuid": "{env}"},
		Body:         `{"name":"renamed"}`,
		BodyEnvelope: "change",
	})
	if err != nil {
		t.Fatalf("DispatchRaw: %v", err)
	}
	if want := `{"change":{"name":"renamed"}}`; got != want {
		t.Errorf("sent body %s, want %s", got, want)
	}
}

func TestSetNestedReplacesScalarParent(t *testing.T) {
	m := map[string]any{}
	SetNested(m, "target", "main")
	SetNested(m, "target.hash", "abc")
	target, ok := m["target"].(map[string]any)
	if !ok || target["hash"] != "abc" {
		t.Fatalf("SetNested produced %#v", m)
	}
}
