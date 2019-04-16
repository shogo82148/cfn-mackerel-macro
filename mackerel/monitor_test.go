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
	ptrUint64 := func(v uint64) *uint64 { return &v }
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
		{
			resp: map[string]interface{}{
				"id":            "2cSZzK3XfmG",
				"type":          "connectivity",
				"name":          "connectivity service1",
				"memo":          "A monitor that checks connectivity.",
				"scopes":        []interface{}{"service1"},
				"excludeScopes": []interface{}{"service1:role3"},
			},
			want: &MonitorConnectivity{
				ID:            "2cSZzK3XfmG",
				Name:          "connectivity service1",
				Memo:          "A monitor that checks connectivity.",
				Type:          MonitorTypeConnectivity,
				Scopes:        []string{"service1"},
				ExcludeScopes: []string{"service1:role3"},
			},
		},
		{
			resp: map[string]interface{}{
				"id":                      "2cSZzK3XfmG",
				"type":                    "service",
				"name":                    "Hatena-Blog - access_num.4xx_count",
				"memo":                    "A monitor that checks the number of 4xx for Hatena Blog",
				"service":                 "Hatena-Blog",
				"duration":                1,
				"metric":                  "access_num.4xx_count",
				"operator":                ">",
				"warning":                 50.0,
				"critical":                100.0,
				"maxCheckAttempts":        3,
				"missingDurationWarning":  360,
				"missingDurationCritical": 720,
				"notificationInterval":    60,
			},
			want: &MonitorServiceMetric{
				ID:                   "2cSZzK3XfmG",
				Name:                 "Hatena-Blog - access_num.4xx_count",
				Memo:                 "A monitor that checks the number of 4xx for Hatena Blog",
				Type:                 MonitorTypeServiceMetric,
				NotificationInterval: 60,

				Service:          "Hatena-Blog",
				Metric:           "access_num.4xx_count",
				Operator:         ">",
				Warning:          ptrFloat64(50.0),
				Critical:         ptrFloat64(100.0),
				Duration:         1,
				MaxCheckAttempts: 3,

				MissingDurationWarning:  ptrUint64(360),
				MissingDurationCritical: ptrUint64(720),
			},
		},
		{
			resp: map[string]interface{}{
				"id":                              "2cSZzK3XfmG",
				"type":                            "external",
				"name":                            "Example Domain",
				"memo":                            "Monitors example.com",
				"method":                          "GET",
				"url":                             "https://example.com",
				"service":                         "Hatena-Blog",
				"notificationInterval":            60,
				"responseTimeWarning":             5000,
				"responseTimeCritical":            10000,
				"responseTimeDuration":            3,
				"containsString":                  "Example",
				"maxCheckAttempts":                3,
				"certificationExpirationWarning":  90,
				"certificationExpirationCritical": 30,
				"isMute":                          false,
				"headers": []interface{}{
					map[string]interface{}{
						"name":  "Cache-Control",
						"value": "no-cache",
					},
				},
			},
			want: &MonitorExternalHTTP{
				ID:                   "2cSZzK3XfmG",
				Name:                 "Example Domain",
				Memo:                 "Monitors example.com",
				Type:                 MonitorTypeExternalHTTP,
				NotificationInterval: 60,

				Method:                          "GET",
				URL:                             "https://example.com",
				MaxCheckAttempts:                3,
				Service:                         "Hatena-Blog",
				ResponseTimeCritical:            ptrFloat64(10000),
				ResponseTimeWarning:             ptrFloat64(5000),
				ResponseTimeDuration:            ptrUint64(3),
				ContainsString:                  "Example",
				CertificationExpirationCritical: ptrUint64(30),
				CertificationExpirationWarning:  ptrUint64(90),
				Headers: []HeaderField{
					HeaderField{
						Name:  "Cache-Control",
						Value: "no-cache",
					},
				},
			},
		},
		{
			resp: map[string]interface{}{
				"id":                   "2cSZzK3XfmG",
				"type":                 "expression",
				"name":                 "role average",
				"memo":                 "Monitors the average of loadavg5",
				"expression":           "avg(roleSlots(\"server:role\",\"loadavg5\"))",
				"operator":             ">",
				"warning":              5.0,
				"critical":             10.0,
				"notificationInterval": 60,
			},
			want: &MonitorExpression{
				ID:                   "2cSZzK3XfmG",
				Name:                 "role average",
				Memo:                 "Monitors the average of loadavg5",
				Type:                 MonitorTypeExpression,
				NotificationInterval: 60,

				Expression: "avg(roleSlots(\"server:role\",\"loadavg5\"))",
				Operator:   ">",
				Warning:    ptrFloat64(5.0),
				Critical:   ptrFloat64(10.0),
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
				t.Errorf("FindMonitor differs: (-got +want)\n%s", diff)
			}
		})

		t.Run(fmt.Sprintf("DeleteMonitor-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("unexpected method: want %s, got %s", http.MethodDelete, r.Method)
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
			got, err := c.DeleteMonitor(context.Background(), "2cSZzK3XfmG")
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("DeleteMonitor differs: (-got +want)\n%s", diff)
			}
		})

	}
}

func TestCreateMonitor(t *testing.T) {
	ptrFloat64 := func(v float64) *float64 { return &v }
	ptrUint64 := func(v uint64) *uint64 { return &v }
	tests := []struct {
		in   Monitor
		want map[string]interface{}
	}{
		{
			in: &MonitorHostMetric{
				Name:                 "disk.aa-00.writes.delta",
				Memo:                 "This monitor is for Hatena Blog.",
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
			want: map[string]interface{}{
				"type":                 "host",
				"name":                 "disk.aa-00.writes.delta",
				"memo":                 "This monitor is for Hatena Blog.",
				"duration":             3.0,
				"metric":               "disk.aa-00.writes.delta",
				"operator":             ">",
				"warning":              20000.0,
				"critical":             400000.0,
				"maxCheckAttempts":     3.0,
				"notificationInterval": 60.0,
				"scopes":               []interface{}{"Hatena-Blog"},
				"excludeScopes":        []interface{}{"Hatena-Bookmark:db-master"},
			},
		},
		{
			in: &MonitorConnectivity{
				Name:          "connectivity service1",
				Memo:          "A monitor that checks connectivity.",
				Scopes:        []string{"service1"},
				ExcludeScopes: []string{"service1:role3"},
			},
			want: map[string]interface{}{
				"type":          "connectivity",
				"name":          "connectivity service1",
				"memo":          "A monitor that checks connectivity.",
				"scopes":        []interface{}{"service1"},
				"excludeScopes": []interface{}{"service1:role3"},
			},
		},
		{
			in: &MonitorServiceMetric{
				Name:                 "Hatena-Blog - access_num.4xx_count",
				Memo:                 "A monitor that checks the number of 4xx for Hatena Blog",
				NotificationInterval: 60,

				Service:          "Hatena-Blog",
				Metric:           "access_num.4xx_count",
				Operator:         ">",
				Warning:          ptrFloat64(50.0),
				Critical:         ptrFloat64(100.0),
				Duration:         1,
				MaxCheckAttempts: 3,

				MissingDurationWarning:  ptrUint64(360),
				MissingDurationCritical: ptrUint64(720),
			},
			want: map[string]interface{}{
				"type":                    "service",
				"name":                    "Hatena-Blog - access_num.4xx_count",
				"memo":                    "A monitor that checks the number of 4xx for Hatena Blog",
				"service":                 "Hatena-Blog",
				"duration":                1.0,
				"metric":                  "access_num.4xx_count",
				"operator":                ">",
				"warning":                 50.0,
				"critical":                100.0,
				"maxCheckAttempts":        3.0,
				"missingDurationWarning":  360.0,
				"missingDurationCritical": 720.0,
				"notificationInterval":    60.0,
			},
		},
		{
			in: &MonitorExternalHTTP{
				Name:                 "Example Domain",
				Memo:                 "Monitors example.com",
				NotificationInterval: 60,

				Method:                          "GET",
				URL:                             "https://example.com",
				MaxCheckAttempts:                3,
				Service:                         "Hatena-Blog",
				ResponseTimeCritical:            ptrFloat64(10000),
				ResponseTimeWarning:             ptrFloat64(5000),
				ResponseTimeDuration:            ptrUint64(3),
				ContainsString:                  "Example",
				CertificationExpirationCritical: ptrUint64(30),
				CertificationExpirationWarning:  ptrUint64(90),
				Headers: []HeaderField{
					HeaderField{
						Name:  "Cache-Control",
						Value: "no-cache",
					},
				},
			},
			want: map[string]interface{}{
				"type":                            "external",
				"name":                            "Example Domain",
				"memo":                            "Monitors example.com",
				"method":                          "GET",
				"url":                             "https://example.com",
				"service":                         "Hatena-Blog",
				"notificationInterval":            60.0,
				"responseTimeWarning":             5000.0,
				"responseTimeCritical":            10000.0,
				"responseTimeDuration":            3.0,
				"containsString":                  "Example",
				"maxCheckAttempts":                3.0,
				"certificationExpirationWarning":  90.0,
				"certificationExpirationCritical": 30.0,
				"headers": []interface{}{
					map[string]interface{}{
						"name":  "Cache-Control",
						"value": "no-cache",
					},
				},
			},
		},
		{
			in: &MonitorExpression{
				Name:                 "role average",
				Memo:                 "Monitors the average of loadavg5",
				NotificationInterval: 60,

				Expression: "avg(roleSlots(\"server:role\",\"loadavg5\"))",
				Operator:   ">",
				Warning:    ptrFloat64(5.0),
				Critical:   ptrFloat64(10.0),
			},
			want: map[string]interface{}{
				"type":                 "expression",
				"name":                 "role average",
				"memo":                 "Monitors the average of loadavg5",
				"expression":           "avg(roleSlots(\"server:role\",\"loadavg5\"))",
				"operator":             ">",
				"warning":              5.0,
				"critical":             10.0,
				"notificationInterval": 60.0,
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("CreateMonitor-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method: want %s, got %s", http.MethodPost, r.Method)
				}
				if r.URL.Path != "/api/v0/monitors" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/monitors", r.URL.Path)
				}

				var data map[string]interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&data); err != nil {
					t.Error(err)
					return
				}
				if diff := cmp.Diff(data, tc.want); diff != "" {
					t.Errorf("CreateMonitor differs: (-got +want)\n%s", diff)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"type":"host"}`) // DUMMY
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

			_, err = c.CreateMonitor(context.Background(), tc.in)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run(fmt.Sprintf("UpdateMonitor-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("unexpected method: want %s, got %s", http.MethodPut, r.Method)
				}
				if r.URL.Path != "/api/v0/monitors/2cSZzK3XfmG" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/monitors/2cSZzK3XfmG", r.URL.Path)
				}

				var data map[string]interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&data); err != nil {
					t.Error(err)
					return
				}
				if diff := cmp.Diff(data, tc.want); diff != "" {
					t.Errorf("CreateMonitor differs: (-got +want)\n%s", diff)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"type":"host"}`) // DUMMY
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

			_, err = c.UpdateMonitor(context.Background(), "2cSZzK3XfmG", tc.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
