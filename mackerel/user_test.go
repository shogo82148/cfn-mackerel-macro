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

func TestFindUsers(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v0/users" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/users", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"users": [
			{
				"id": "2cdkEV8JB5d",
				"screenName": "shogo82148@gmail.com",
				"email": "shogo82148@gmail.com",
				"authority": "owner",
				"isInRegistrationProcess": false,
				"isMFAEnabled": false,
				"authenticationMethods": [
				  "password",
				  "github",
				  "google"
				],
				"joinedAt": 1411403412
			}
			]
		}`)
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

	got, err := c.FindUsers(context.Background())
	if err != nil {
		t.Error(err)
	}
	want := []*User{
		{
			ID:         "2cdkEV8JB5d",
			ScreenName: "shogo82148@gmail.com",
			Email:      "shogo82148@gmail.com",
			Authority:  UserAuthorityOwner,
			AuthenticationMethods: []UserAuthenticationMethod{
				UserAuthenticationMethodPassword,
				UserAuthenticationMethodGitHub,
				UserAuthenticationMethodGoogle,
			},
			JoinedAt: 1411403412,
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("users differs: (-got +want)\n%s", diff)
	}
}

func TestDeleteUser(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v0/users/2cdkEV8JB5d" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/users", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
				"id": "2cdkEV8JB5d",
				"screenName": "shogo82148@gmail.com",
				"email": "shogo82148@gmail.com",
				"authority": "owner",
				"isInRegistrationProcess": false,
				"isMFAEnabled": false,
				"authenticationMethods": [
				  "password",
				  "github",
				  "google"
				],
				"joinedAt": 1411403412
			}`)
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

	got, err := c.DeleteUser(context.Background(), "2cdkEV8JB5d")
	if err != nil {
		t.Error(err)
	}
	want := &User{
		ID:         "2cdkEV8JB5d",
		ScreenName: "shogo82148@gmail.com",
		Email:      "shogo82148@gmail.com",
		Authority:  UserAuthorityOwner,
		AuthenticationMethods: []UserAuthenticationMethod{
			UserAuthenticationMethodPassword,
			UserAuthenticationMethodGitHub,
			UserAuthenticationMethodGoogle,
		},
		JoinedAt: 1411403412,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("users differs: (-got +want)\n%s", diff)
	}
}
