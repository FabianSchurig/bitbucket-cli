package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// TestDispatch_InternalAPI_UsesCSRFAndSessionCookies verifies that requests
// targeting Bitbucket's internal API (URLs containing "/!api/internal/") are
// authenticated with the csrftoken + cloud.session.token cookies and the
// X-CSRFToken header — the only auth combination that endpoint accepts.
//
// It also verifies that HTTP Basic Auth is *not* sent for internal URLs,
// because Bitbucket's internal API rejects requests that present an
// Authorization header (it expects browser-style cookie auth).
func TestDispatch_InternalAPI_UsesCSRFAndSessionCookies(t *testing.T) {
	output.Format = "json"

	var capturedAuth, capturedXCSRF, capturedXReqWith, capturedAccept string
	var capturedCookies []*http.Cookie
	var capturedBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedXCSRF = r.Header.Get("X-CSRFToken")
		capturedXReqWith = r.Header.Get("X-Requested-With")
		capturedAccept = r.Header.Get("Accept")
		capturedCookies = r.Cookies()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
		}
		capturedBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"values": []any{}})
	}))
	defer srv.Close()

	// Build a BBClient by hand that holds the internal-API cookies. We
	// deliberately do NOT call resty.SetBasicAuth so that the dispatcher is
	// the only place auth gets applied, exactly mirroring production wiring
	// from client.NewClientWithConfig.
	bb := &client.BBClient{
		Client:            resty.New().SetBaseURL(srv.URL),
		CSRFToken:         "csrf-abc",
		CloudSessionToken: "session-xyz",
	}

	// Use the test server URL but keep the "/!api/internal/" marker so the
	// dispatcher recognises this as an internal-API request. The dispatcher
	// passes absolute URL templates through unchanged, so we can swap host.
	url := srv.URL + "/!api/internal/workspaces/myorg/projects/PROJ/branch-restrictions/by-pattern/master"
	err := handlers.Dispatch(context.Background(), bb, handlers.Request{
		Method:      http.MethodPut,
		URLTemplate: url,
		Body:        `{"values":[]}`,
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}

	if capturedAuth != "" {
		t.Errorf("Authorization header should be empty for internal API, got %q", capturedAuth)
	}
	if capturedXCSRF != "csrf-abc" {
		t.Errorf("X-CSRFToken = %q, want %q", capturedXCSRF, "csrf-abc")
	}
	if capturedXReqWith != "XMLHttpRequest" {
		t.Errorf("X-Requested-With = %q, want XMLHttpRequest", capturedXReqWith)
	}
	if !strings.Contains(capturedAccept, "application/json") {
		t.Errorf("Accept = %q, want to contain application/json", capturedAccept)
	}

	cookieMap := map[string]string{}
	for _, c := range capturedCookies {
		cookieMap[c.Name] = c.Value
	}
	if cookieMap["csrftoken"] != "csrf-abc" {
		t.Errorf("csrftoken cookie = %q, want csrf-abc", cookieMap["csrftoken"])
	}
	if cookieMap["cloud.session.token"] != "session-xyz" {
		t.Errorf("cloud.session.token cookie = %q, want session-xyz", cookieMap["cloud.session.token"])
	}

	if capturedBody != `{"values":[]}` {
		t.Errorf("body forwarded incorrectly: %q", capturedBody)
	}
}

// TestDispatch_InternalAPI_MissingCookiesFails verifies a clear error when an
// internal API request is attempted without the required cookies (rather than
// silently sending an unauthenticated request that Bitbucket would 401).
func TestDispatch_InternalAPI_MissingCookiesFails(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("server should not be called when internal-API auth is missing, got %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	bb := &client.BBClient{Client: resty.New().SetBaseURL(srv.URL)}

	url := srv.URL + "/!api/internal/workspaces/myorg/projects/PROJ/branch-restrictions/by-pattern/master"
	err := handlers.Dispatch(context.Background(), bb, handlers.Request{
		Method:      http.MethodGet,
		URLTemplate: url,
	})
	if err == nil {
		t.Fatal("expected error for internal API request without csrf/session tokens")
	}
	if !strings.Contains(err.Error(), "BITBUCKET_CSRF_TOKEN") ||
		!strings.Contains(err.Error(), "BITBUCKET_CLOUD_SESSION_TOKEN") {
		t.Errorf("error should mention required env vars, got: %v", err)
	}
}

// TestDispatch_PublicAPI_UnaffectedByInternalAuth verifies that requests to
// the public REST API still rely on whatever Basic Auth the client was
// configured with and do NOT get cookie auth applied.
func TestDispatch_PublicAPI_UnaffectedByInternalAuth(t *testing.T) {
	output.Format = "json"

	var capturedAuth string
	var capturedCookies []*http.Cookie

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedCookies = r.Cookies()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1})
	}))
	defer srv.Close()

	bb := &client.BBClient{
		Client:            resty.New().SetBaseURL(srv.URL).SetBasicAuth("u", "p"),
		Username:          "u",
		Token:             "p",
		CSRFToken:         "csrf-abc",
		CloudSessionToken: "session-xyz",
	}

	err := handlers.Dispatch(context.Background(), bb, handlers.Request{
		Method:      http.MethodGet,
		URLTemplate: "/repositories/myorg/myrepo",
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if capturedAuth == "" {
		t.Error("expected Authorization header on public-API request")
	}
	for _, c := range capturedCookies {
		if c.Name == "csrftoken" || c.Name == "cloud.session.token" {
			t.Errorf("internal-API cookie %q leaked onto public-API request", c.Name)
		}
	}
}

// TestDispatch_InternalAPI_SuppressesPreconfiguredBasicAuth verifies that
// when the underlying resty client has Basic Auth configured, requests to
// /!api/internal/ still go out with NO Authorization header. The internal
// endpoint returns 401 if both cookies and Basic Auth are present, so
// suppression must work even when Basic Auth was set up at client
// construction time.
//
// This is the common production case: users set BITBUCKET_TOKEN (which would
// historically trigger client-level SetBasicAuth) AND the internal-API
// cookies, and call internal endpoints. The dispatcher applies Basic Auth
// per-request and the client installs a defensive PreRequestHook, so neither
// path can leak Authorization onto an internal request.
func TestDispatch_InternalAPI_SuppressesPreconfiguredBasicAuth(t *testing.T) {
	output.Format = "json"

	var capturedAuth string
	var capturedCookies []*http.Cookie

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedCookies = r.Cookies()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"values": []any{}})
	}))
	defer srv.Close()

	// Build the client through the production constructor so the defensive
	// PreRequestHook is installed, then ALSO call SetBasicAuth directly on
	// the underlying resty client to simulate a worst-case where Basic Auth
	// has been wired in at every level. The internal-API request must still
	// go out without an Authorization header.
	bb, err := client.NewClientWithConfig("u", "p", srv.URL, "csrf-abc", "session-xyz")
	if err != nil {
		t.Fatalf("NewClientWithConfig: %v", err)
	}
	bb.SetBasicAuth("u", "p")

	url := srv.URL + "/!api/internal/workspaces/myorg/projects/PROJ/branch-restrictions/group-by-branch/"
	if err := handlers.Dispatch(context.Background(), bb, handlers.Request{
		Method:      http.MethodGet,
		URLTemplate: url,
	}); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}

	if capturedAuth != "" {
		t.Errorf("Authorization header should be suppressed for internal API even with Basic Auth pre-configured, got %q", capturedAuth)
	}

	cookieMap := map[string]string{}
	for _, c := range capturedCookies {
		cookieMap[c.Name] = c.Value
	}
	if cookieMap["csrftoken"] != "csrf-abc" || cookieMap["cloud.session.token"] != "session-xyz" {
		t.Errorf("expected internal-API cookies to be sent, got %v", cookieMap)
	}
}
