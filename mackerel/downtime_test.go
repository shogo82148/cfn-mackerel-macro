package mackerel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestFindDowntimes(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/downtimes" {
			t.Errorf("unexpected request path: want %s, got %s", "/api/v0/downtimes", r.URL.Path)
		}
		ret := map[string]interface{}{
			"downtimes": []interface{}{
				map[string]interface{}{
					"id":       "abcde0",
					"name":     "Maintenance #0",
					"memo":     "Memo #0",
					"start":    1563600000,
					"duration": 120,
				},
				map[string]interface{}{
					"id":       "abcde1",
					"name":     "Maintenance #1",
					"memo":     "Memo #1",
					"start":    1563700000,
					"duration": 60,
					"recurrence": map[string]interface{}{
						"interval": 3,
						"type":     "weekly",
						"weekdays": []string{
							"Monday",
							"Thursday",
							"Saturday",
						},
					},
					"serviceScopes": []string{
						"service1",
					},
					"serviceExcludeScopes": []string{
						"service2",
					},
					"roleScopes": []string{
						"service3: role1",
					},
					"roleExcludeScopes": []string{
						"service1: role1",
					},
					"monitorScopes": []string{
						"monitor0",
					},
					"monitorExcludeScopes": []string{
						"monitor1",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
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

	got, err := c.FindDowntimes(context.Background())
	if err != nil {
		t.Error(err)
	}

	want := []*Downtime{
		{
			ID:       "abcde0",
			Name:     "Maintenance #0",
			Memo:     "Memo #0",
			Start:    1563600000,
			Duration: 120,
		},
		{
			ID:       "abcde1",
			Name:     "Maintenance #1",
			Memo:     "Memo #1",
			Start:    1563700000,
			Duration: 60,
			Recurrence: &DowntimeRecurrence{
				Type:     DowntimeRecurrenceTypeWeekly,
				Interval: 3,
				Weekdays: []DowntimeWeekday{
					DowntimeWeekday(time.Monday),
					DowntimeWeekday(time.Thursday),
					DowntimeWeekday(time.Saturday),
				},
			},
			ServiceScopes: []string{
				"service1",
			},
			ServiceExcludeScopes: []string{
				"service2",
			},
			RoleScopes: []string{
				"service3: role1",
			},
			RoleExcludeScopes: []string{
				"service1: role1",
			},
			MonitorScopes: []string{
				"monitor0",
			},
			MonitorExcludeScopes: []string{
				"monitor1",
			},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("FindDowntimes differs: (-got +want)\n%s", diff)
	}
}

func TestCreateDowntime(t *testing.T) {
	const (
		downtimeID = "9rxGOHfVF8F"
	)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		want := map[string]interface{}{
			"name":     "downtime name",
			"memo":     "memo",
			"start":    1234567890.0,
			"duration": 10.0,
		}
		if diff := cmp.Diff(data, want); diff != "" {
			t.Errorf("downtime differs: (-got +want)\n%s", diff)
		}
		w.WriteHeader(http.StatusOK)
		want["id"] = downtimeID
		json.NewEncoder(w).Encode(want)
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

	param := &Downtime{
		Name:     "downtime name",
		Memo:     "memo",
		Start:    1234567890,
		Duration: 10,
	}
	got, err := c.CreateDowntime(context.Background(), param)
	if err != nil {
		t.Error(err)
	}

	param.ID = downtimeID
	if diff := cmp.Diff(got, param); diff != "" {
		t.Errorf("downtime differs: (-got +want)\n%s", diff)
	}
}
