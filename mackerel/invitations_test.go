package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRevokeInvitation(t *testing.T) {
	const (
		email = "macopy@example.com"
	)
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: got %s, want %s", http.MethodPost, r.Method)
		}
		if r.RequestURI != "/api/v0/invitations/revoke" {
			t.Errorf("want /api/v0/invitations/revoke, got %s", r.RequestURI)
		}
		var body struct {
			Email string `json:"email"`
		}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Email != email {
			t.Errorf("unexpected email address, want %s, got %s", email, body.Email)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"success":true}`)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		BaseURL:    u,
		APIKey:     "DUMMY-API-KEY",
		HTTPClient: ts.Client(),
	}

	err = c.RevokeInvitation(context.Background(), email)
	if err != nil {
		t.Error(err)
	}
}
