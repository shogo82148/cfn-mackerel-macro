package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindNotificationGroups(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v0/notification-groups" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/notification-groups", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"notificationGroups":[{
			"id": "2oWY1xPXrco",
			"name": "Example notification group",
			"notificationLevel": "all",
			"childNotificationGroupIds": [],
			"childChannelIds": [
			  "2vh7AZ21abc"
			],
			"monitors": [
			  {
				"id": "2qtozU21abc",
				"skipDefault": false
			  }
			],
			"services": [
			  {
				"name": "Example-Service-1"
			  },
			  {
				"name": "Example-Service-2"
			  }
			]
		  }]}`)
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

	got, err := c.FindNotificationGroups(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []*NotificationGroup{
		{
			ID:                        "2oWY1xPXrco",
			Name:                      "Example notification group",
			NotificationLevel:         NotificationLevelAll,
			ChildNotificationGroupIDs: []string{},
			ChildChannelIDs:           []string{"2vh7AZ21abc"},
			Monitors: []NotificationGroupMonitor{
				{
					ID: "2qtozU21abc",
				},
			},
			Services: []NotificationGroupService{
				{
					Name: "Example-Service-1",
				},
				{
					Name: "Example-Service-2",
				},
			},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("notification groups missmatch (-got +want):\n%s", diff)
	}
}

func TestCreateNotificationGroups(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: want %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v0/notification-groups" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/notification-groups", r.URL.Path)
		}
		var got interface{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&got); err != nil {
			t.Error(err)
		}
		want := map[string]interface{}{
			"name":                      "Example notification group",
			"notificationLevel":         "all",
			"childChannelIds":           []interface{}{},
			"childNotificationGroupIds": []interface{}{},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("notification group missmatch (-got +want):\n%s", diff)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "2oWY1xPXrco",
			"name": "Example notification group",
			"notificationLevel": "all",
			"childNotificationGroupIds": [],
			"childChannelIds": [],
			"monitors": [],
			"services": []
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

	got, err := c.CreateNotificationGroup(context.Background(), &NotificationGroup{
		Name:              "Example notification group",
		NotificationLevel: NotificationLevelAll,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := &NotificationGroup{
		ID:                        "2oWY1xPXrco",
		Name:                      "Example notification group",
		NotificationLevel:         NotificationLevelAll,
		ChildNotificationGroupIDs: []string{},
		ChildChannelIDs:           []string{},
		Monitors:                  []NotificationGroupMonitor{},
		Services:                  []NotificationGroupService{},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("notification group missmatch (-got +want):\n%s", diff)
	}
}

func TestUpdateNotificationGroups(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: want %s, got %s", http.MethodPut, r.Method)
		}
		if r.URL.Path != "/api/v0/notification-groups/2oWY1xPXrco" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/notification-groups/2oWY1xPXrco", r.URL.Path)
		}
		var got interface{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&got); err != nil {
			t.Error(err)
		}
		want := map[string]interface{}{
			"name":                      "Example notification group",
			"notificationLevel":         "all",
			"childChannelIds":           []interface{}{},
			"childNotificationGroupIds": []interface{}{},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("notification group missmatch (-got +want):\n%s", diff)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "2oWY1xPXrco",
			"name": "Example notification group",
			"notificationLevel": "all",
			"childNotificationGroupIds": [],
			"childChannelIds": [],
			"monitors": [],
			"services": []
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

	got, err := c.UpdateNotificationGroup(context.Background(), "2oWY1xPXrco", &NotificationGroup{
		Name:              "Example notification group",
		NotificationLevel: NotificationLevelAll,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := &NotificationGroup{
		ID:                        "2oWY1xPXrco",
		Name:                      "Example notification group",
		NotificationLevel:         NotificationLevelAll,
		ChildNotificationGroupIDs: []string{},
		ChildChannelIDs:           []string{},
		Monitors:                  []NotificationGroupMonitor{},
		Services:                  []NotificationGroupService{},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("notification group missmatch (-got +want):\n%s", diff)
	}
}

func TestDeleteNotificationGroups(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: want %s, got %s", http.MethodDelete, r.Method)
		}
		if r.URL.Path != "/api/v0/notification-groups/2oWY1xPXrco" {
			t.Errorf("unexpected path: want %s, got %s", "/api/v0/notification-groups/2oWY1xPXrco", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "2oWY1xPXrco",
			"name": "Example notification group",
			"notificationLevel": "all",
			"childNotificationGroupIds": [],
			"childChannelIds": [],
			"monitors": [],
			"services": []
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

	got, err := c.DeleteNotificationGroup(context.Background(), "2oWY1xPXrco")
	if err != nil {
		t.Fatal(err)
	}
	want := &NotificationGroup{
		ID:                        "2oWY1xPXrco",
		Name:                      "Example notification group",
		NotificationLevel:         NotificationLevelAll,
		ChildNotificationGroupIDs: []string{},
		ChildChannelIDs:           []string{},
		Monitors:                  []NotificationGroupMonitor{},
		Services:                  []NotificationGroupService{},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("notification group missmatch (-got +want):\n%s", diff)
	}
}
