package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
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
				if param.Name != "AWSIntegration" {
					t.Errorf("unexpected name: want %s, got %s", "AWSIntegration", param.Name)
				}
				if param.Region != "ap-northeast-1" {
					t.Errorf("unexpected region: want %s, got %s", "ap-northeast-1", param.Region)
				}
				if pointer.StringValue(param.Key) != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("unexpected key: want %s, got %s", "AKIAIOSFODNN7EXAMPLE", pointer.StringValue(param.Key))
				}
				if pointer.StringValue(param.SecretKey) != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
					t.Errorf("unexpected secret key: want %s, got %s", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", pointer.StringValue(param.SecretKey))
				}
				if param.IncludedTags != "TagName:TagValue,\"Tag:Name\":\"Tag,Value\",'Tag\"Name':\"Tag' Value\"" {
					t.Errorf("unexpected included tags: want %q, got %q", "TagName:TagValue,\"Tag:Name\":\"Tag,Value\",'Tag\"Name':\"Tag' Value\"", param.IncludedTags)
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
					"Name":  "TagName",
					"Value": "TagValue",
				},
				map[string]interface{}{
					"Name":  "Tag:Name",
					"Value": "Tag,Value",
				},
				map[string]interface{}{
					"Name":  "Tag\"Name",
					"Value": "Tag' Value",
				},
			},
			"Services": []interface{}{
				map[string]interface{}{
					"ServiceId": "S3",
					"Enable":    "true",
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
