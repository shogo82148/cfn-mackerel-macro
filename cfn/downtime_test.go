package cfn

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateDowntime(t *testing.T) {
	r := &downtime{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createDowntime: func(ctx context.Context, param *mackerel.Downtime) (*mackerel.Downtime, error) {
					want := &mackerel.Downtime{
						Name:     "Maintenance #1",
						Memo:     "Memo #1",
						Start:    1563700000,
						Duration: 60,
						Recurrence: &mackerel.DowntimeRecurrence{
							Type:     mackerel.DowntimeRecurrenceTypeWeekly,
							Interval: 3,
							Weekdays: []mackerel.DowntimeWeekday{
								mackerel.DowntimeWeekday(time.Monday),
								mackerel.DowntimeWeekday(time.Thursday),
								mackerel.DowntimeWeekday(time.Saturday),
							},
						},
						ServiceScopes: []string{
							"service1",
						},
						ServiceExcludeScopes: []string{
							"service2",
						},
						RoleScopes: []string{
							"service3:role1",
						},
						RoleExcludeScopes: []string{
							"service1:role1",
						},
						MonitorScopes: []string{
							"monitor0",
						},
						MonitorExcludeScopes: []string{
							"monitor1",
						},
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("param differs: (-got +want)\n%s", diff)
					}
					return &mackerel.Downtime{
						ID:       "3yAYEDLXKL5",
						Name:     "Maintenance #1",
						Memo:     "Memo #1",
						Start:    1563700000,
						Duration: 60,
						Recurrence: &mackerel.DowntimeRecurrence{
							Type:     mackerel.DowntimeRecurrenceTypeWeekly,
							Interval: 3,
							Weekdays: []mackerel.DowntimeWeekday{
								mackerel.DowntimeWeekday(time.Monday),
								mackerel.DowntimeWeekday(time.Thursday),
								mackerel.DowntimeWeekday(time.Saturday),
							},
						},
						ServiceScopes: []string{
							"service1",
						},
						ServiceExcludeScopes: []string{
							"service2",
						},
						RoleScopes: []string{
							"service3:role1",
						},
						RoleExcludeScopes: []string{
							"service1:role1",
						},
						MonitorScopes: []string{
							"monitor0",
						},
						MonitorExcludeScopes: []string{
							"monitor1",
						},
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Downtime",
			LogicalResourceID: "Downtime",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Name":     "Maintenance #1",
				"Memo":     "Memo #1",
				"Start":    1563700000.0,
				"Duration": 60,
				"Recurrence": map[string]interface{}{
					"Type":     "weekly",
					"Interval": 3,
					"Weekdays": []interface{}{
						"Monday", "Thursday", "Saturday",
					},
				},
				"ServiceScopes": []interface{}{
					"mkr:test-org:service:service1",
				},
				"ServiceExcludeScopes": []interface{}{
					"mkr:test-org:service:service2",
				},
				"RoleScopes": []interface{}{
					"mkr:test-org:role:service3:role1",
				},
				"RoleExcludeScopes": []interface{}{
					"mkr:test-org:role:service1:role1",
				},
				"MonitorScopes": []interface{}{
					"mkr:test-org:monitor:monitor0",
				},
				"MonitorExcludeScopes": []interface{}{
					"mkr:test-org:monitor:monitor1",
				},
			},
		},
	}
	id, _, err := r.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:downtime:3yAYEDLXKL5" {
		t.Errorf("unexpected downtime id: want %s, got %s", "mkr:test-org:downtime:3yAYEDLXKL5", id)
	}
}

func TestDeleteDowntime(t *testing.T) {
}
