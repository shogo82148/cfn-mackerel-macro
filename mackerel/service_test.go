package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindServices(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"services":[{"name": "awesome-service", "memo": "some memo", "roles": ["role1", "role2"]}]}`)
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

	got, err := c.FindServices(context.Background())
	if err != nil {
		t.Error(err)
	}
	want := []*Service{
		{
			Name: "awesome-service",
			Memo: "some memo",
			Roles: []string{
				"role1", "role2",
			},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("metadata differs: (-got +want)\n%s", diff)
	}
}
