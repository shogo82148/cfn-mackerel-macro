package cfn

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateMonitor_MonitorConnectivity(t *testing.T) {
	m := &monitor{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createMonitor: func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
					want := &mackerel.MonitorConnectivity{
						Name:                 "foo-bar",
						Memo:                 "monitor",
						NotificationInterval: 60,
						Scopes:               []string{"my-service"},
						ExcludeScopes:        []string{"my-service:my-role"},
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("monitor differs: (-got +want)\n%s", diff)
					}
					want.ID = "3yAYEDLXKL5"
					return want, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Monitor",
			LogicalResourceID: "Monitor",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type":                 "connectivity",
				"Name":                 "foo-bar",
				"Memo":                 "monitor",
				"Scopes":               []interface{}{"mkr:test-org:service:my-service"},
				"ExcludeScopes":        []interface{}{"mkr:test-org:role:my-service:my-role"},
				"NotificationInterval": 60,
			},
		},
	}
	id, param, err := m.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:monitor:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["MonitorId"].(string) != "3yAYEDLXKL5" {
		t.Errorf("unexpected monitor id, want %s, got %s", "3yAYEDLXKL5", param["MonitorId"].(string))
	}
	if param["Name"].(string) != "foo-bar" {
		t.Errorf("unexpected name, want %s, got %s", "foo-bar", param["Name"].(string))
	}
	if param["Type"].(string) != "connectivity" {
		t.Errorf("unexpected type, want %s, got %s", "connectivity", param["Type"].(string))
	}
}

func TestCreateMonitor_MonitorHostMetric(t *testing.T) {
	ptrFloat64 := func(v float64) *float64 { return &v }
	m := &monitor{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createMonitor: func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
					want := &mackerel.MonitorHostMetric{
						Name:                 "disk.aa-00.writes.delta",
						Memo:                 "This monitor is for Hatena Blog.",
						Duration:             3,
						Metric:               "disk.aa-00.writes.delta",
						Operator:             ">",
						Warning:              ptrFloat64(20000.0),
						Critical:             ptrFloat64(400000.0),
						MaxCheckAttempts:     3,
						Scopes:               []string{"Hatena-Blog"},
						ExcludeScopes:        []string{"Hatena-Bookmark:db-master"},
						NotificationInterval: 60,
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("monitor differs: (-got +want)\n%s", diff)
					}
					want.ID = "3yAYEDLXKL5"
					return want, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Monitor",
			LogicalResourceID: "Monitor",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type":                 "host",
				"Name":                 "disk.aa-00.writes.delta",
				"Memo":                 "This monitor is for Hatena Blog.",
				"Duration":             3,
				"Metric":               "disk.aa-00.writes.delta",
				"Operator":             ">",
				"Warning":              20000.0,
				"Critical":             400000.0,
				"MaxCheckAttempts":     3,
				"Scopes":               []interface{}{"mkr:test-org:service:Hatena-Blog"},
				"ExcludeScopes":        []interface{}{"mkr:test-org:role:Hatena-Bookmark:db-master"},
				"NotificationInterval": 60,
			},
		},
	}
	id, param, err := m.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:monitor:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["MonitorId"].(string) != "3yAYEDLXKL5" {
		t.Errorf("unexpected monitor id, want %s, got %s", "3yAYEDLXKL5", param["MonitorId"].(string))
	}
	if param["Name"].(string) != "disk.aa-00.writes.delta" {
		t.Errorf("unexpected name, want %s, got %s", "foo-bar", param["Name"].(string))
	}
	if param["Type"].(string) != "host" {
		t.Errorf("unexpected type, want %s, got %s", "host", param["Type"].(string))
	}
}

func TestCreateMonitor_MonitorServiceMetric(t *testing.T) {
	ptrFloat64 := func(v float64) *float64 { return &v }
	ptrUint64 := func(v uint64) *uint64 { return &v }
	m := &monitor{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createMonitor: func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
					want := &mackerel.MonitorServiceMetric{
						Name:                    "Hatena-Blog - access_num.4xx_count",
						Memo:                    "A monitor that checks the number of 4xx for Hatena Blog",
						Duration:                1,
						Service:                 "Hatena-Blog",
						Metric:                  "access_num.4xx_count",
						Operator:                ">",
						Warning:                 ptrFloat64(50.0),
						Critical:                ptrFloat64(100.0),
						MaxCheckAttempts:        3,
						NotificationInterval:    60,
						MissingDurationWarning:  ptrUint64(360),
						MissingDurationCritical: ptrUint64(720),
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("monitor differs: (-got +want)\n%s", diff)
					}
					want.ID = "3yAYEDLXKL5"
					return want, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Monitor",
			LogicalResourceID: "Monitor",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type":                    "service",
				"Name":                    "Hatena-Blog - access_num.4xx_count",
				"Memo":                    "A monitor that checks the number of 4xx for Hatena Blog",
				"Service":                 "mkr:test-org:service:Hatena-Blog",
				"Duration":                1,
				"Metric":                  "access_num.4xx_count",
				"Operator":                ">",
				"Warning":                 50.0,
				"Critical":                100.0,
				"MaxCheckAttempts":        3,
				"MissingDurationWarning":  360,
				"MissingDurationCritical": 720,
				"NotificationInterval":    60,
			},
		},
	}
	id, param, err := m.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:monitor:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["MonitorId"].(string) != "3yAYEDLXKL5" {
		t.Errorf("unexpected monitor id, want %s, got %s", "3yAYEDLXKL5", param["MonitorId"].(string))
	}
	if param["Name"].(string) != "Hatena-Blog - access_num.4xx_count" {
		t.Errorf("unexpected name, want %s, got %s", "Hatena-Blog - access_num.4xx_count", param["Name"].(string))
	}
	if param["Type"].(string) != "service" {
		t.Errorf("unexpected type, want %s, got %s", "service", param["Type"].(string))
	}
}

func TestCreateMonitor_MonitorExternalHTTP(t *testing.T) {
	ptrFloat64 := func(v float64) *float64 { return &v }
	ptrUint64 := func(v uint64) *uint64 { return &v }
	m := &monitor{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createMonitor: func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
					want := &mackerel.MonitorExternalHTTP{
						Name:                 "Example Domain",
						Memo:                 "Monitors example.com",
						NotificationInterval: 60,

						Method:                          "GET",
						URL:                             "https://example.com",
						MaxCheckAttempts:                3,
						Service:                         "Hatena-Blog",
						ResponseTimeCritical:            ptrFloat64(10000),
						ResponseTimeWarning:             ptrFloat64(5000),
						ResponseTimeDuration:            ptrUint64(3),
						ContainsString:                  "Example",
						CertificationExpirationCritical: ptrUint64(30),
						CertificationExpirationWarning:  ptrUint64(90),
						Headers: []mackerel.HeaderField{
							mackerel.HeaderField{
								Name:  "Cache-Control",
								Value: "no-cache",
							},
						},
					}
					if diff := cmp.Diff(param, want); diff != "" {
						t.Errorf("monitor differs: (-got +want)\n%s", diff)
					}
					want.ID = "3yAYEDLXKL5"
					return want, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:Monitor",
			LogicalResourceID: "Monitor",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Type":                            "external",
				"Name":                            "Example Domain",
				"Memo":                            "Monitors example.com",
				"Method":                          "GET",
				"Url":                             "https://example.com",
				"Service":                         "mkr:test-org:service:Hatena-Blog",
				"NotificationInterval":            60.0,
				"ResponseTimeWarning":             5000.0,
				"ResponseTimeCritical":            10000.0,
				"ResponseTimeDuration":            3.0,
				"ContainsString":                  "Example",
				"MaxCheckAttempts":                3.0,
				"CertificationExpirationWarning":  90.0,
				"CertificationExpirationCritical": 30.0,
				"Headers": []interface{}{
					map[string]interface{}{
						"Name":  "Cache-Control",
						"Value": "no-cache",
					},
				},
			},
		},
	}
	id, param, err := m.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if id != "mkr:test-org:monitor:3yAYEDLXKL5" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:host:3yAYEDLXKL5", id)
	}
	if param["MonitorId"].(string) != "3yAYEDLXKL5" {
		t.Errorf("unexpected monitor id, want %s, got %s", "3yAYEDLXKL5", param["MonitorId"].(string))
	}
	if param["Name"].(string) != "Example Domain" {
		t.Errorf("unexpected name, want %s, got %s", "Example Domain", param["Name"].(string))
	}
	if param["Type"].(string) != "external" {
		t.Errorf("unexpected type, want %s, got %s", "external", param["Type"].(string))
	}
}
