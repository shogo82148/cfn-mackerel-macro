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

func TestCreateRole(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: want %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v0/services/awesome-service/roles" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/services/awesome-service/roles", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"application","memo":"the application of awesome service"}`)
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

	got, err := c.CreateRole(context.Background(), "awesome-service", &CreateRoleParam{
		Name: "application",
		Memo: "the application of awesome service",
	})
	if err != nil {
		t.Error(err)
	}
	want := &Role{
		Name: "application",
		Memo: "the application of awesome service",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("role differs: (-got +want)\n%s", diff)
	}
}

func TestDeleteRole(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: want %s, got %s", http.MethodDelete, r.Method)
		}
		if r.URL.Path != "/api/v0/services/awesome-service/roles/application" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/services/awesome-service/roles", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"application","memo":"the application of awesome service"}`)
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

	got, err := c.DeleteRole(context.Background(), "awesome-service", "application")
	if err != nil {
		t.Error(err)
	}
	want := &Role{
		Name: "application",
		Memo: "the application of awesome service",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("role differs: (-got +want)\n%s", diff)
	}
}
