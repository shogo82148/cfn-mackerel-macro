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

// Ensure each monitor type conforms to the Monitor interface.
var (
	_ Monitor = (*MonitorConnectivity)(nil)
	_ Monitor = (*MonitorHostMetric)(nil)
	_ Monitor = (*MonitorServiceMetric)(nil)
	_ Monitor = (*MonitorExternalHTTP)(nil)
	_ Monitor = (*MonitorExpression)(nil)
)

func TestFindMonitors(t *testing.T) {
	ptrFloat64 := func(v float64) *float64 { return &v }
	tests := []struct {
		resp map[string]interface{} // the response of the mackerel api
		want Monitor
	}{
		{
			resp: map[string]interface{}{
				"id":                   "2cSZzK3XfmG",
				"type":                 "host",
				"name":                 "disk.aa-00.writes.delta",
				"memo":                 "This monitor is for Hatena Blog.",
				"duration":             3,
				"metric":               "disk.aa-00.writes.delta",
				"operator":             ">",
				"warning":              20000.0,
				"critical":             400000.0,
				"maxCheckAttempts":     3,
				"notificationInterval": 60,
				"scopes":               []interface{}{"Hatena-Blog"},
				"excludeScopes":        []interface{}{"Hatena-Bookmark:db-master"},
			},
			want: &MonitorHostMetric{
				ID:                   "2cSZzK3XfmG",
				Name:                 "disk.aa-00.writes.delta",
				Memo:                 "This monitor is for Hatena Blog.",
				Type:                 MonitorTypeHostMetric,
				NotificationInterval: 60,

				Metric:           "disk.aa-00.writes.delta",
				Operator:         ">",
				Warning:          ptrFloat64(20000.0),
				Critical:         ptrFloat64(400000.0),
				Duration:         3,
				MaxCheckAttempts: 3,

				Scopes:        []string{"Hatena-Blog"},
				ExcludeScopes: []string{"Hatena-Bookmark:db-master"},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("FindMonitors-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/monitors" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode([]interface{}{tc.resp})
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
			got, err := c.FindMonitors(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, []Monitor{tc.want}); diff != "" {
				t.Errorf("FindMonitors differs: (-got +want)\n%s", diff)
			}
		})

		t.Run(fmt.Sprintf("FindMonitor-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/monitors/2cSZzK3XfmG" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards/2cSZzK3XfmG", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode(tc.resp)
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
			got, err := c.FindMonitor(context.Background(), "2cSZzK3XfmG")
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("FindMonitors differs: (-got +want)\n%s", diff)
			}
		})
	}
}
