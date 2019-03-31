package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateNotificationChannel(t *testing.T) {
	ch := &notificationChannel{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createNotificationChannel: func(ctx context.Context, ch mackerel.NotificationChannel) (mackerel.NotificationChannel, error) {
					slack := ch.(*mackerel.NotificationChannelSlack)
					if ch.NotificationChannelName() != "channel-foobar" {
						t.Errorf("unexpected name, want %s, got %s", "channel-foobar", ch.NotificationChannelName())
					}
					if slack.URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
						t.Errorf("unexpected url: want %s, got %s", "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX", slack.URL)
					}
					slack.ID = "3yAYEDLXKL5"
					return slack, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "xxxx",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom::NotificationChannel",
			LogicalResourceID: "Channel",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type": "slack",
				"Name": "channel-foobar",
				"Url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
				"Mentions": map[string]interface{}{
					"Ok":      "ok message",
					"Warning": "warning message",
				},
				"Events": []interface{}{"alert"},
			},
		},
	}
	id, param, err := ch.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:notification-channel:3yAYEDLXKL5" {
		t.Errorf("unexpected notification channel id: want %s, got %s", "mkr:test-org:notification-channel:3yAYEDLXKL5", id)
	}
	if param["Name"].(string) != "channel-foobar" {
		t.Errorf("unexpected name, want %s, got %s", "channel-foobar", param["Name"].(string))
	}
}
