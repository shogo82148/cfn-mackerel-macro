package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
	"github.com/shogo82148/pointer"
)

func TestCreateDashboard(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			createDashboard: func(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
				want := &mackerel.Dashboard{
					Title:   "dashboard-foobar",
					Memo:    "memo",
					URLPath: "my-dashboard",
					Widgets: []mackerel.Widget{

						&mackerel.WidgetGraph{
							Title: "Host Graph",
							Graph: &mackerel.GraphHost{
								HostID: "host-id",
								Name:   "some.metric",
							},
							Range: &mackerel.GraphRangeRelative{
								Period: 3600,
								Offset: 0,
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetGraph{
							Title: "Role Graph",
							Graph: &mackerel.GraphRole{
								RoleFullname: "awesome-service:role-hogehoge",
								Name:         "some.metric",
								IsStacked:    true,
							},
							Range: &mackerel.GraphRangeAbsolute{
								Start: 1234567890,
								End:   1234567890,
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetGraph{
							Title: "Service Graph",
							Graph: &mackerel.GraphService{
								ServiceName: "awesome-service",
								Name:        "some.metric",
							},
							Range: &mackerel.GraphRangeRelative{
								Period: 3600,
								Offset: 0,
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetGraph{
							Title: "Expression Graph",
							Graph: &mackerel.GraphExpression{
								Expression: `avg(roleSlots("server:role","loadavg5"))`,
							},
							Range: &mackerel.GraphRangeRelative{
								Period: 3600,
								Offset: 0,
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetValue{
							Title: "Host Value",
							Metric: &mackerel.MetricHost{
								HostID: "host-id",
								Name:   "some.metric",
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetValue{
							Title: "Service Value",
							Metric: &mackerel.MetricService{
								ServiceName: "awesome-service",
								Name:        "some.metric",
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetValue{
							Title: "Expression Value",
							Metric: &mackerel.MetricExpression{
								Expression: `avg(roleSlots("server:role","loadavg5"))`,
							},
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetMarkdown{
							Title:    "Markdown",
							Markdown: "# Some Awesome Service\n- Markdown Text Here",
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 32,
							},
						},

						&mackerel.WidgetAlertStatus{
							Title:        "alert status",
							RoleFullname: pointer.String("awesome-service:role-hogehoge"),
							Layout: &mackerel.Layout{
								X:      0,
								Y:      0,
								Width:  24,
								Height: 6,
							},
						},
					},
				}
				if diff := cmp.Diff(param, want); diff != "" {
					t.Errorf("param differs: (-got +want)\n%s", diff)
				}
				ret := *param
				ret.ID = "dashboard-id"
				return &ret, nil
			},
		},
	}
	event := cfn.Event{
		RequestType:       cfn.RequestCreate,
		RequestID:         "",
		ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:      "Custom::Dashboard",
		LogicalResourceID: "Dashboard",
		StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
		ResourceProperties: map[string]interface{}{
			"Title":   "dashboard-foobar",
			"Memo":    "memo",
			"UrlPath": "my-dashboard",
			"Widgets": []interface{}{
				// Host Graph
				map[string]interface{}{
					"Type":  "graph",
					"Title": "Host Graph",
					"Graph": map[string]interface{}{
						"Type": "host",
						"Name": "some.metric",
						"Host": "mkr:test-org:host:host-id",
					},
					"Range": map[string]interface{}{
						"Type":   "relative",
						"Period": "3600",
						"Offset": "0",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Role Graph
				map[string]interface{}{
					"Type":  "graph",
					"Title": "Role Graph",
					"Graph": map[string]interface{}{
						"Type":      "role",
						"Role":      "mkr:test-org:role:awesome-service:role-hogehoge",
						"Name":      "some.metric",
						"IsStacked": "true",
					},
					"Range": map[string]interface{}{
						"Type":  "absolute",
						"Start": "1234567890",
						"End":   "1234567890",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Service Graph
				map[string]interface{}{
					"Type":  "graph",
					"Title": "Service Graph",
					"Graph": map[string]interface{}{
						"Type":    "service",
						"Service": "mkr:test-org:service:awesome-service",
						"Name":    "some.metric",
					},
					"Range": map[string]interface{}{
						"Type":   "relative",
						"Period": "3600",
						"Offset": "0",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Expression Graph
				map[string]interface{}{
					"Type":  "graph",
					"Title": "Expression Graph",
					"Graph": map[string]interface{}{
						"Type":       "expression",
						"Expression": `avg(roleSlots("server:role","loadavg5"))`,
					},
					"Range": map[string]interface{}{
						"Type":   "relative",
						"Period": "3600",
						"Offset": "0",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Host Value
				map[string]interface{}{
					"Type":  "value",
					"Title": "Host Value",
					"Metric": map[string]interface{}{
						"Type": "host",
						"Host": "mkr:test-org:host:host-id",
						"Name": "some.metric",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Service Value
				map[string]interface{}{
					"Type":  "value",
					"Title": "Service Value",
					"Metric": map[string]interface{}{
						"Type":    "service",
						"Service": "mkr:test-org:service:awesome-service",
						"Name":    "some.metric",
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Expression Value
				map[string]interface{}{
					"Type":  "value",
					"Title": "Expression Value",
					"Metric": map[string]interface{}{
						"Type":       "expression",
						"Expression": `avg(roleSlots("server:role","loadavg5"))`,
					},
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Markdown
				map[string]interface{}{
					"Type":     "markdown",
					"Title":    "Markdown",
					"Markdown": "# Some Awesome Service\n- Markdown Text Here",
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "32",
					},
				},

				// Alert Status Widget
				map[string]interface{}{
					"Type":  "alertStatus",
					"Title": "alert status",
					"Role":  "mkr:test-org:role:awesome-service:role-hogehoge",
					"Layout": map[string]interface{}{
						"X":      "0",
						"Y":      "0",
						"Width":  "24",
						"Height": "6",
					},
				},
			},
		},
	}
	id, _, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Fatal(err)
	}
	if id != "mkr:test-org:dashboard:dashboard-id" {
		t.Errorf("unexpected dashboard id: want %s, got %s", "mkr:test-org:dashboard:dashboard-id", id)
	}
}

func TestUpdateDashboard(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			updateDashboard: func(ctx context.Context, id string, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
				if id != "dashboard-id" {
					t.Errorf("unexpected dashboard id: want %s, got %s", "dashboard-id", id)
				}
				want := &mackerel.Dashboard{
					Title:   "dashboard-foobar",
					Memo:    "memo",
					URLPath: "my-dashboard",
					Widgets: []mackerel.Widget{},
				}
				if diff := cmp.Diff(param, want); diff != "" {
					t.Errorf("param differs: (-got +want)\n%s", diff)
				}
				ret := *param
				ret.ID = "dashboard-id"
				return &ret, nil
			},
		},
	}
	event := cfn.Event{
		RequestType:        cfn.RequestUpdate,
		RequestID:          "",
		ResponseURL:        "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:       "Custom::Dashboard",
		LogicalResourceID:  "Dashboard",
		PhysicalResourceID: "mkr:test-org:dashboard:dashboard-id",
		StackID:            "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
		OldResourceProperties: map[string]interface{}{
			"Title":   "dashboard-foobar",
			"Memo":    "memo",
			"UrlPath": "my-dashboard",
			"Widgets": []interface{}{},
			"Roles":   []interface{}{},
		},
		ResourceProperties: map[string]interface{}{
			"Title":   "dashboard-foobar",
			"Memo":    "memo",
			"UrlPath": "my-dashboard",
			"Widgets": []interface{}{},
			"Roles":   []interface{}{},
		},
	}
	id, _, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Fatal(err)
	}
	if id != "mkr:test-org:dashboard:dashboard-id" {
		t.Errorf("unexpected dashboard id: want %s, got %s", "mkr:test-org:dashboard:dashboard-id", id)
	}
}
