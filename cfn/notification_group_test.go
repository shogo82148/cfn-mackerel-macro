package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateNotificationGroup(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			createNotificationGroup: func(ctx context.Context, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
				want := &mackerel.NotificationGroup{
					Name:                      "NotificationGroup",
					NotificationLevel:         mackerel.NotificationLevelAll,
					ChildNotificationGroupIDs: []string{"child-group"},
					ChildChannelIDs:           []string{"child-channel"},
					Monitors: []mackerel.NotificationGroupMonitor{
						{
							ID:          "monitor",
							SkipDefault: false,
						},
					},
					Services: []mackerel.NotificationGroupService{
						{
							Name: "service",
						},
					},
				}
				if diff := cmp.Diff(group, want); diff != "" {
					t.Errorf("group differs: (-got +want)\n%s", diff)
				}
				return &mackerel.NotificationGroup{
					ID:   "group-id",
					Name: "NotificationGroup",
				}, nil
			},
		},
	}

	event := cfn.Event{
		RequestType:       cfn.RequestCreate,
		RequestID:         "request-id123",
		ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:      "Custom::NotificationGroup",
		LogicalResourceID: "NotificationGroup",
		StackID:           "arn:aws:cloudformation:ap-northeast-1:123456789012:stack/foobar/12345678-1234-1234-1234-123456789abc",
		ResourceProperties: map[string]interface{}{
			"Name": "NotificationGroup",
			"ChildNotificationGroupIds": []interface{}{
				"mkr:test-org:notification-group:child-group",
			},
			"ChildChannelIds": []interface{}{
				"mkr:test-org:notification-channel:child-channel",
			},
			"Services": []interface{}{
				map[string]interface{}{
					"Id": "mkr:test-org:service:service",
				},
			},
			"Monitors": []interface{}{
				map[string]interface{}{
					"Id": "mkr:test-org:monitor:monitor",
				},
			},
		},
	}
	id, data, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:notification-group:group-id" {
		t.Errorf("unexpected aws integration id: want %s, got %s", "mkr:test-org:notification-group:group-id", id)
	}
	name, _ := data["Name"].(string)
	if name != "NotificationGroup" {
		t.Errorf("unexpected name: want %s, got %s", "NotificationGroup", name)
	}
}

func TestUpdateNotificationGroup(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			updateNotificationGroup: func(ctx context.Context, id string, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
				if id != "group-id" {
					t.Errorf("unexpected id: want %s, got %s", "group-id", id)
				}
				want := &mackerel.NotificationGroup{
					Name:                      "NotificationGroup",
					NotificationLevel:         mackerel.NotificationLevelAll,
					ChildNotificationGroupIDs: []string{"child-group"},
					ChildChannelIDs:           []string{"child-channel"},
					Monitors: []mackerel.NotificationGroupMonitor{
						{
							ID:          "monitor",
							SkipDefault: false,
						},
					},
					Services: []mackerel.NotificationGroupService{
						{
							Name: "service",
						},
					},
				}
				if diff := cmp.Diff(group, want); diff != "" {
					t.Errorf("group differs: (-got +want)\n%s", diff)
				}
				return &mackerel.NotificationGroup{
					ID:   "group-id",
					Name: "NotificationGroup",
				}, nil
			},
		},
	}

	event := cfn.Event{
		RequestType:        cfn.RequestUpdate,
		RequestID:          "request-id123",
		ResponseURL:        "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:       "Custom::NotificationGroup",
		PhysicalResourceID: "mkr:test-org:notification-group:group-id",
		LogicalResourceID:  "NotificationGroup",
		StackID:            "arn:aws:cloudformation:ap-northeast-1:123456789012:stack/foobar/12345678-1234-1234-1234-123456789abc",
		OldResourceProperties: map[string]interface{}{
			"Name": "NotificationGroup",
		},
		ResourceProperties: map[string]interface{}{
			"Name": "NotificationGroup",
			"ChildNotificationGroupIds": []interface{}{
				"mkr:test-org:notification-group:child-group",
			},
			"ChildChannelIds": []interface{}{
				"mkr:test-org:notification-channel:child-channel",
			},
			"Services": []interface{}{
				map[string]interface{}{
					"Id": "mkr:test-org:service:service",
				},
			},
			"Monitors": []interface{}{
				map[string]interface{}{
					"Id": "mkr:test-org:monitor:monitor",
				},
			},
		},
	}
	id, data, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:notification-group:group-id" {
		t.Errorf("unexpected aws integration id: want %s, got %s", "mkr:test-org:notification-group:group-id", id)
	}
	name, _ := data["Name"].(string)
	if name != "NotificationGroup" {
		t.Errorf("unexpected name: want %s, got %s", "NotificationGroup", name)
	}
}

func TestDeleteNotificationGroup(t *testing.T) {
	f := &Function{
		org: &mackerel.Org{
			Name: "test-org",
		},
		client: &fakeMackerelClient{
			deleteNotificationGroup: func(ctx context.Context, id string) (*mackerel.NotificationGroup, error) {
				if id != "group-id" {
					t.Errorf("unexpected id: want %s, got %s", "group-id", id)
				}
				return &mackerel.NotificationGroup{
					ID:   "group-id",
					Name: "NotificationGroup",
				}, nil
			},
		},
	}

	event := cfn.Event{
		RequestType:        cfn.RequestDelete,
		RequestID:          "request-id123",
		ResponseURL:        "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
		ResourceType:       "Custom::NotificationGroup",
		PhysicalResourceID: "mkr:test-org:notification-group:group-id",
		LogicalResourceID:  "NotificationGroup",
		StackID:            "arn:aws:cloudformation:ap-northeast-1:123456789012:stack/foobar/12345678-1234-1234-1234-123456789abc",
		OldResourceProperties: map[string]interface{}{
			"Name": "NotificationGroup",
		},
		ResourceProperties: map[string]interface{}{
			"Name": "NotificationGroup",
		},
	}
	id, _, err := f.Handle(context.Background(), event)
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:notification-group:group-id" {
		t.Errorf("unexpected aws integration id: want %s, got %s", "mkr:test-org:notification-group:group-id", id)
	}
}
