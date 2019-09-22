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

func TestFindDashboard(t *testing.T) {
	ptrString := func(v string) *string { return &v }
	tests := []struct {
		resp map[string]interface{} // the response of the mackerel api
		want *Dashboard
	}{
		/////////// Alert Status Widgets
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":         "alertStatus",
						"title":        "status title",
						"roleFullname": "service:role",
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetAlertStatus{
						Type:         WidgetTypeAlertStatus,
						Title:        "status title",
						RoleFullname: ptrString("service:role"),
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":         "alertStatus",
						"title":        "status title",
						"roleFullname": nil, // roleFullname may be nil
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetAlertStatus{
						Type:         WidgetTypeAlertStatus,
						Title:        "status title",
						RoleFullname: nil,
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		/////////// Graph Widgets
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]string{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphHost{
							Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]interface{}{
							"type":         "role",
							"roleFullname": "service-foo:role-bar",
							"name":         "host-graph",
							"isStacked":    true,
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphRole{
							Type:         GraphTypeRole,
							RoleFullname: "service-foo:role-bar",
							Name:         "host-graph",
							IsStacked:    true,
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]string{
							"type":        "service",
							"serviceName": "service-foo",
							"name":        "metric-name",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphService{
							Type:        GraphTypeService,
							ServiceName: "service-foo",
							Name:        "metric-name",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]string{
							"type":       "expression",
							"expression": "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphExpression{
							Type:       GraphTypeExpression,
							Expression: "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]string{
							"type": "unknown",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphUnknown{
							Type: GraphTypeUnknown,
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},

		/////////// Metric Widgets
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]string{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "hogehoge",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricHost{
							Type:   MetricTypeHost,
							HostID: "host-foobar",
							Name:   "hogehoge",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]string{
							"type":        "service",
							"serviceName": "service-foobar",
							"name":        "hogehoge",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricService{
							Type:        MetricTypeService,
							ServiceName: "service-foobar",
							Name:        "hogehoge",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]string{
							"type":       "expression",
							"expression": "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricExpression{
							Type:       MetricTypeExpression,
							Expression: "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]string{
							"type": "unknown",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricUnknown{
							Type: MetricTypeUnknown,
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":     "markdown",
						"title":    "markdown title",
						"markdown": "*FOOBAR*",
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetMarkdown{
						Type:     WidgetTypeMarkdown,
						Title:    "markdown title",
						Markdown: "*FOOBAR*",
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},

		////////// Graph ranges
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"range": map[string]interface{}{
							"type":   "relative",
							"period": 3600,
							"offset": -3600,
						},
						"graph": map[string]string{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Range: &GraphRangeRelative{
							Type:   GraphRangeTypeRelative,
							Period: 3600,
							Offset: -3600,
						},
						Graph: &GraphHost{
							Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"range": map[string]interface{}{
							"type":  "absolute",
							"start": -1234567890,
							"end":   1234567890,
						},
						"graph": map[string]string{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Range: &GraphRangeAbsolute{
							Type:  GraphRangeTypeAbsolute,
							Start: -1234567890,
							End:   1234567890,
						},
						Graph: &GraphHost{
							Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("FindDashboards-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards" {
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
			got, err := c.FindDashboards(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, []*Dashboard{tc.want}); diff != "" {
				t.Errorf("FindDashboards differs: (-got +want)\n%s", diff)
			}
		})

		t.Run(fmt.Sprintf("FindDashboard-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards/foobar" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards/foobar", r.URL.Path)
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
			got, err := c.FindDashboard(context.Background(), "foobar")
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("FindDashboard differs: (-got +want)\n%s", diff)
			}
		})

		t.Run(fmt.Sprintf("DeleteDashboard-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("unexpected method: want %s, got %s", http.MethodDelete, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards/foobar" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards/foobar", r.URL.Path)
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
			got, err := c.DeleteDashboard(context.Background(), "foobar")
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("DeleteDashboard differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCreateDashboard(t *testing.T) {
	ptrString := func(v string) *string { return &v }
	tests := []struct {
		in   *Dashboard
		want map[string]interface{}
	}{
		/////////// Alert Status Widgets
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetAlertStatus{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeAlertStatus,
						Title:        "status title",
						RoleFullname: ptrString("service:role"),
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":         "alertStatus",
						"title":        "status title",
						"roleFullname": "service:role",
					},
				},
			},
		},
		/////////// Graph Widgets
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphHost{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]interface{}{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphRole{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeRole,
							RoleFullname: "service-foo:role-bar",
							Name:         "host-graph",
							IsStacked:    true,
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]interface{}{
							"type":         "role",
							"roleFullname": "service-foo:role-bar",
							"name":         "host-graph",
							"isStacked":    true,
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphService{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeService,
							ServiceName: "service-name",
							Name:        "host-graph",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]interface{}{
							"type":        "service",
							"serviceName": "service-name",
							"name":        "host-graph",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphExpression{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeExpression,
							Expression: "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]interface{}{
							"type":       "expression",
							"expression": "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricHost{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-metric",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]interface{}{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-metric",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricService{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeService,
							ServiceName: "service-foobar",
							Name:        "service-metric",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]interface{}{
							"type":        "service",
							"serviceName": "service-foobar",
							"name":        "service-metric",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetValue{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeValue,
						Title: "metric title",
						Metric: &MetricExpression{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeExptression,
							Expression: "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "value",
						"title": "metric title",
						"metric": map[string]interface{}{
							"type":       "expression",
							"expression": "host(22CXRB3pZmu, memory.*)",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetMarkdown{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeMarkdown,
						Title:    "markdown title",
						Markdown: "*FOOBAR*",
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":     "markdown",
						"title":    "markdown title",
						"markdown": "*FOOBAR*",
					},
				},
			},
		},

		////////// Graph ranges
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Range: &GraphRangeRelative{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphRangeTypeRelative,
							Period: 3600,
							Offset: -3600,
						},
						Graph: &GraphHost{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"range": map[string]interface{}{
							"type":   "relative",
							"period": 3600.0,
							"offset": -3600.0,
						},
						"graph": map[string]interface{}{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
			},
		},
		{
			in: &Dashboard{
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						// the type field will be autocomplete from the Golang's type.
						// Type:  WidgetTypeGraph,
						Title: "graph title",
						Range: &GraphRangeAbsolute{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphRangeTypeAbsolute,
							Start: -1234567890,
							End:   1234567890,
						},
						Graph: &GraphHost{
							// the type field will be autocomplete from the Golang's type.
							// Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
			},
			want: map[string]interface{}{
				"title":   "title",
				"urlPath": "url path",
				"widgets": []interface{}{
					map[string]interface{}{
						"type":  "graph",
						"title": "graph title",
						"range": map[string]interface{}{
							"type":  "absolute",
							"start": -1234567890.0,
							"end":   1234567890.0,
						},
						"graph": map[string]interface{}{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("CreateDashboard-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var data map[string]interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&data); err != nil {
					t.Error(err)
					return
				}
				if diff := cmp.Diff(data, tc.want); diff != "" {
					t.Errorf("DeleteDashboard differs: (-got +want)\n%s", diff)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "{}")
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

			_, err = c.CreateDashboard(context.Background(), tc.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
