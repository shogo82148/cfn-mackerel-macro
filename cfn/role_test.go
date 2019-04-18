package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateRole(t *testing.T) {
	r := &role{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createRole: func(ctx context.Context, serviceName string, param *mackerel.CreateRoleParam) (*mackerel.Role, error) {
					if serviceName != "awesome-service" {
						t.Errorf("unexpected service name, want %s, got %s", "awesome-service", serviceName)
					}
					if param.Name != "role-app" {
						t.Errorf("unexpected name, want %s, got %s", "role-app", param.Name)
					}
					return &mackerel.Role{
						Name: param.Name,
					}, nil
				},
				putRoleMetaData: func(ctx context.Context, serviceName, roleName, namespace string, v interface{}) error {
					if namespace != "cloudformation" {
						t.Errorf("unexpected namespace: want cloudformation, got %s", namespace)
					}
					if serviceName != "awesome-service" {
						t.Errorf("unexpected service name, want %s, got %s", "awesome-service", serviceName)
					}
					if roleName != "role-app" {
						t.Errorf("unexpected host id, want %sn got %s", "role-app", roleName)
					}
					meta := v.(metadata)
					want := metadata{
						StackName: "foobar",
						StackID:   "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
						LogicalID: "Role",
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
				"Service": "mkr:test-org:service:awesome-service",
				"Name":    "role-app",
			},
		},
	}
	id, param, err := r.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:role:awesome-service:role-app" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:role:awesome-service:role-app", id)
	}
	if param["Name"].(string) != "role-app" {
		t.Errorf("unexpected name, want %s, got %s", "role-app", param["Name"].(string))
	}
	if param["FullName"].(string) != "awesome-service:role-app" {
		t.Errorf("unexpected name, want %s, got %s", "awesome-service:role-app", param["FullName"].(string))
	}
}
