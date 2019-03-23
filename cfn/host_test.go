package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
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
				putHostMetaData: func(ctx context.Context, hostID, namespace string, v interface{}) error {
					if namespace != "cloudformation" {
						t.Errorf("unexpected namespace: want cloudformation, got %s", namespace)
					}
					if hostID != "3yAYEDLXKL5" {
						t.Errorf("unexpected host id, want %sn got %s", "3yAYEDLXKL5", hostID)
					}
					meta := v.(metadata)
					want := metadata{
						StackName: "foobar",
						StackID:   "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
						LogicalID: "Host",
					}
					if diff := cmp.Diff(meta, want); diff != "" {
						t.Errorf("metadata differs: (-got +want)\n%s", diff)
					}
					return nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Host",
			LogicalResourceID: "Host",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
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
