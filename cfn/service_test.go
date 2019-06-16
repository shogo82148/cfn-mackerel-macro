package cfn

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateService(t *testing.T) {
	s := &service{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createService: func(ctx context.Context, param *mackerel.CreateServiceParam) (*mackerel.Service, error) {
					if param.Name != "awesome-service" {
						t.Errorf("unexpected name, want %s, got %s", "awesome-service", param.Name)
					}
					return &mackerel.Service{
						Name:  param.Name,
						Memo:  "",
						Roles: []string{},
					}, nil
				},
				putServiceMetaData: func(ctx context.Context, serviceName, namespace string, v interface{}) error {
					if namespace != "cloudformation" {
						t.Errorf("unexpected namespace: want cloudformation, got %s", namespace)
					}
					if serviceName != "awesome-service" {
						t.Errorf("unexpected host id, want %sn got %s", "awesome-service", serviceName)
					}
					meta := v.(metadata)
					want := metadata{
						StackName: "foobar",
						StackID:   "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
						LogicalID: "Service",
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
			ResourceType:      "Custom:Service",
			LogicalResourceID: "Service",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Name": "awesome-service",
			},
		},
	}
	id, param, err := s.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:service:awesome-service" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:service:awesome-service", id)
	}
	if param["Name"].(string) != "awesome-service" {
		t.Errorf("unexpected name, want %s, got %s", "awesome-service", param["Name"].(string))
	}
}

func TestCreateService_AlreadyExists(t *testing.T) {
	s := &service{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createService: func(ctx context.Context, param *mackerel.CreateServiceParam) (*mackerel.Service, error) {
					return nil, mkrError{
						statusCode: http.StatusBadRequest,
					}
				},
				putServiceMetaData: func(ctx context.Context, serviceName, namespace string, v interface{}) error {
					if namespace != "cloudformation" {
						t.Errorf("unexpected namespace: want cloudformation, got %s", namespace)
					}
					if serviceName != "awesome-service" {
						t.Errorf("unexpected host id, want %sn got %s", "awesome-service", serviceName)
					}
					meta := v.(metadata)
					want := metadata{
						StackName: "foobar",
						StackID:   "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
						LogicalID: "Service",
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
			ResourceType:      "Custom:Service",
			LogicalResourceID: "Service",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Name": "awesome-service",
			},
		},
	}
	id, param, err := s.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:service:awesome-service" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:service:awesome-service", id)
	}
	if param["Name"].(string) != "awesome-service" {
		t.Errorf("unexpected name, want %s, got %s", "awesome-service", param["Name"].(string))
	}
}

func TestDeleteService(t *testing.T) {
	var deleted bool
	s := &service{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				deleteService: func(ctx context.Context, serviceName string) (*mackerel.Service, error) {
					deleted = true
					if serviceName != "awesome-service" {
						t.Errorf("unexpected service name, want awesome-service, got %s", serviceName)
					}
					return &mackerel.Service{
						Name: serviceName,
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestDelete,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Service",
			LogicalResourceID: "Service",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Name": "awesome-service",
			},
			PhysicalResourceID: "mkr:test-org:service:awesome-service",
		},
	}
	_, _, err := s.delete(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !deleted {
		t.Error("the service is not deleted")
	}
}
