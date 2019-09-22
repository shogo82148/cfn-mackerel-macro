package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateDashboard(t *testing.T) {
	ptrString := func(v string) *string { return &v }
	d := &dashboard{
		Function: &Function{
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
							&mackerel.WidgetAlertStatus{
								Title:        "alert status",
								RoleFullname: ptrString("awesome-service:role-hogehoge"),
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
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Dashboard",
			LogicalResourceID: "Dashboard",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Title":   "dashboard-foobar",
				"Memo":    "memo",
				"UrlPath": "my-dashboard",
				"Widgets": []interface{}{
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
				"Roles": []interface{}{},
			},
		},
	}
	id, _, err := d.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:dashboard:dashboard-id" {
		t.Errorf("unexpected dashboard id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
}
