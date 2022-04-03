package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
	"github.com/shogo82148/pointer"
)

func TestCreateAWSIntegration(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			createAWSIntegration: func(ctx context.Context, param *mackerel.AWSIntegration) (*mackerel.AWSIntegration, error) {
				want := &mackerel.AWSIntegration{
					Name:         "AWSIntegration",
					Region:       "ap-northeast-1",
					Key:          pointer.String("AKIAIOSFODNN7EXAMPLE"),
					SecretKey:    pointer.String("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
					IncludedTags: "TagKey:TagValue,\"Tag:Key\":\"Tag,Value\",'Tag\"Key':\"Tag' Value\"",
					Services: map[string]*mackerel.AWSIntegrationService{
						"Billing": {
							Enable:          false,
							ExcludedMetrics: []string{},
						},
						"S3": {
							Enable:          true,
							ExcludedMetrics: []string{"s3.some-metric"},
						},
					},
				}
				if diff := cmp.Diff(want, param); diff != "" {
					t.Errorf("creation parameter missmatch: (-want/+got):\n%s", diff)
				}
				return &mackerel.AWSIntegration{
					ID: "integration-id",
				}, nil
			},
		},
	}

	event := cfn.Event{
		RequestType:       cfn.RequestCreate,
		RequestID:         "request-id123",
		ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:      "Custom::AWSIntegration",
		LogicalResourceID: "AWSIntegration",
		StackID:           "arn:aws:cloudformation:ap-northeast-1:123456789012:stack/foobar/12345678-1234-1234-1234-123456789abc",
		ResourceProperties: map[string]interface{}{
			"Name":      "AWSIntegration",
			"Region":    "ap-northeast-1",
			"Key":       "AKIAIOSFODNN7EXAMPLE",
			"SecretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"IncludedTags": []interface{}{
				map[string]interface{}{
					"Key":   "TagKey",
					"Value": "TagValue",
				},
				map[string]interface{}{
					"Key":   "Tag:Key",
					"Value": "Tag,Value",
				},
				map[string]interface{}{
					"Key":   "Tag\"Key",
					"Value": "Tag' Value",
				},
			},
			"Services": []interface{}{
				map[string]interface{}{
					"ServiceId": "S3",
					"ExcludedMetrics": []interface{}{
						"s3.some-metric",
					},
				},
				map[string]interface{}{
					"ServiceId": "Billing",
					"Enable":    "false",
				},
			},
		},
	}
	id, _, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:aws-integration:integration-id" {
		t.Errorf("unexpected aws integration id: want %s, got %s", "mkr:test-org:aws-integration:integration-id", id)
	}
}

func TestCreateAWSIntegrationExternalID(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			createAWSIntegrationExternalID: func(ctx context.Context) (string, error) {
				return "external-id", nil
			},
		},
	}

	event := cfn.Event{
		RequestType:        cfn.RequestCreate,
		RequestID:          "request-id123",
		ResponseURL:        "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:       "Custom::AWSIntegrationExternalId",
		LogicalResourceID:  "AWSIntegrationExternalId",
		StackID:            "arn:aws:cloudformation:ap-northeast-1:123456789012:stack/foobar/12345678-1234-1234-1234-123456789abc",
		ResourceProperties: map[string]interface{}{},
	}
	id, _, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:aws-integration-external-id:external-id" {
		t.Errorf("unexpected aws integration external id: want %s, got %s", "mkr:test-org:aws-integration-external-id:external-id", id)
	}
}
