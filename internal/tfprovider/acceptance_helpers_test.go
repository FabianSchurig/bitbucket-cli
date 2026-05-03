package tfprovider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func TestAccEnsureRepoUserWritePermissionRestoresPriorState(t *testing.T) {
	for _, tc := range []struct {
		name          string
		oldPermission string
		wantPuts      []string
		wantDeletes   int
	}{
		{
			name:          "none is restored by delete",
			oldPermission: "none",
			wantPuts:      []string{"write"},
			wantDeletes:   1,
		},
		{
			name:          "read is restored by put",
			oldPermission: "read",
			wantPuts:      []string{"write", "read"},
			wantDeletes:   0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var putPermissions []string
			deleteCount := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					_ = json.NewEncoder(w).Encode(map[string]string{"permission": tc.oldPermission})
				case http.MethodPut:
					var body map[string]string
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode PUT body: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					putPermissions = append(putPermissions, body["permission"])
					_ = json.NewEncoder(w).Encode(map[string]string{"permission": body["permission"]})
				case http.MethodDelete:
					deleteCount++
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			defer srv.Close()

			c, err := client.NewClientWithConfig("user", "token", srv.URL, "", "")
			if err != nil {
				t.Fatalf("NewClientWithConfig: %v", err)
			}
			restore, err := testAccEnsureRepoUserWritePermission(context.Background(), c, "ws", "repo", "{user}")
			if err != nil {
				t.Fatalf("testAccEnsureRepoUserWritePermission: %v", err)
			}
			if len(putPermissions) != 1 || putPermissions[0] != "write" {
				t.Fatalf("initial PUT permissions = %#v, want [write]", putPermissions)
			}
			if err := restore(context.Background()); err != nil {
				t.Fatalf("restore: %v", err)
			}

			if len(putPermissions) != len(tc.wantPuts) {
				t.Fatalf("PUT permissions = %#v, want %#v", putPermissions, tc.wantPuts)
			}
			for i, want := range tc.wantPuts {
				if putPermissions[i] != want {
					t.Fatalf("PUT permissions = %#v, want %#v", putPermissions, tc.wantPuts)
				}
			}
			if deleteCount != tc.wantDeletes {
				t.Fatalf("DELETE count = %d, want %d", deleteCount, tc.wantDeletes)
			}
		})
	}
}
