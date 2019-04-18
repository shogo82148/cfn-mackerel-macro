package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateOrg(t *testing.T) {
	r := &org{
		Function: &Function{
			client: &fakeMackerelClient{
				getOrg: func(ctx context.Context) (*mackerel.Org, error) {
					return &mackerel.Org{
						Name: "test-org",
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:        cfn.RequestCreate,
			RequestID:          "",
			ResponseURL:        "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:       "Custom:Org",
			LogicalResourceID:  "Org",
			StackID:            "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{},
		},
	}
	id, param, err := r.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org", id)
	}
	if param["Name"].(string) != "test-org" {
		t.Errorf("unexpected name, want %s, got %s", "test-org", param["Name"].(string))
	}
}
