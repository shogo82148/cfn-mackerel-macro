package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateHost(t *testing.T) {
	h := &host{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createHost: func(ctx context.Context, param *mackerel.CreateHostParam) (string, error) {
					if param.Name != "host-foobar" {
						t.Errorf("unexpected name, want %s, got %s", "host-foobar", param.Name)
					}
					return "3yAYEDLXKL5", nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "http://example.com/",
			ResourceType:      "Custom:Host",
			LogicalResourceID: "Host",
			StackID:           "arn:hogehoge",
			ResourceProperties: map[string]interface{}{
				"Name":  "host-foobar",
				"Roles": []interface{}{"mkr:test-org:role:awesome-service:role-hogehoge"},
			},
		},
	}
	id, param, err := h.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:host:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["Name"].(string) != "host-foobar" {
		t.Errorf("unexpected name, want %s, got %s", "host-foobar", param["Name"].(string))
	}
}
