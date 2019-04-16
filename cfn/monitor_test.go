package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateMonitor_MonitorConnectivity(t *testing.T) {
	m := &monitor{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createMonitor: func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
					want := &mackerel.MonitorConnectivity{
						Name:                 "foo-bar",
						Memo:                 "monitor",
						NotificationInterval: 60,
						Scopes:               []string{"my-service"},
						ExcludeScopes:        []string{"my-service:my-role"},
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("monitor differs: (-got +want)\n%s", diff)
					}
					want.ID = "3yAYEDLXKL5"
					return want, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Monitor",
			LogicalResourceID: "Monitor",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type":                 "connectivity",
				"Name":                 "foo-bar",
				"Memo":                 "monitor",
				"Scopes":               []interface{}{"mkr:test-org:service:my-service"},
				"ExcludeScopes":        []interface{}{"mkr:test-org:role:my-service:my-role"},
				"NotificationInterval": 60,
			},
		},
	}
	id, param, err := m.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:monitor:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["MonitorId"].(string) != "3yAYEDLXKL5" {
		t.Errorf("unexpected monitor id, want %s, got %s", "3yAYEDLXKL5", param["MonitorId"].(string))
	}
	if param["Name"].(string) != "foo-bar" {
		t.Errorf("unexpected name, want %s, got %s", "foo-bar", param["Name"].(string))
	}
	if param["Type"].(string) != "connectivity" {
		t.Errorf("unexpected type, want %s, got %s", "connectivity", param["Type"].(string))
	}
}
