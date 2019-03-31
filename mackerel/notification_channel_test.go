package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindNotificationChannels(t *testing.T) {
	tests := []struct {
		resp map[string]interface{} // the response of the mackerel api
		want []NotificationChannel
	}{
		// email type
		{
			resp: map[string]interface{}{
				"channels": []interface{}{
					map[string]interface{}{
						"id":      "ch-foobar",
						"name":    "notification-test",
						"type":    "email",
						"emails":  []interface{}{"john.doe@example.com"},
						"userIds": []interface{}{"user-john-doe"},
						"events":  []interface{}{"alert", "alertGroup"},
					},
				},
			},
			want: []NotificationChannel{
				&NotificationChannelEmail{
					ID:      "ch-foobar",
					Name:    "notification-test",
					Type:    NotificationChannelTypeEmail,
					Emails:  []string{"john.doe@example.com"},
					UserIDs: []string{"user-john-doe"},
					Events:  []NotificationEvent{NotificationEventAlert, NotificationEventAlertGroup},
				},
			},
		},

		// slack
		{
			resp: map[string]interface{}{
				"channels": []interface{}{
					map[string]interface{}{
						"id":                "ch-foobar",
						"name":              "notification-test",
						"type":              "slack",
						"url":               "http://example.com",
						"enabledGraphImage": true,
						"events":            []interface{}{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"},
					},
				},
			},
			want: []NotificationChannel{
				&NotificationChannelSlack{
					ID:                "ch-foobar",
					Name:              "notification-test",
					Type:              NotificationChannelTypeSlack,
					URL:               "http://example.com",
					EnabledGraphImage: true,
					Events: []NotificationEvent{
						NotificationEventAlert, NotificationEventAlertGroup, NotificationEventHostStatus, NotificationEventHostRegister, NotificationEventHostRetire, NotificationEventMonitor,
					},
				},
			},
		},

		// webhook
		{
			resp: map[string]interface{}{
				"channels": []interface{}{
					map[string]interface{}{
						"id":     "ch-foobar",
						"name":   "notification-test",
						"type":   "webhook",
						"url":    "http://example.com",
						"events": []interface{}{"alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"},
					},
				},
			},
			want: []NotificationChannel{
				&NotificationChannelWebHook{
					ID:   "ch-foobar",
					Name: "notification-test",
					Type: NotificationChannelTypeWebHook,
					URL:  "http://example.com",
					Events: []NotificationEvent{
						NotificationEventAlert, NotificationEventAlertGroup, NotificationEventHostStatus, NotificationEventHostRegister, NotificationEventHostRetire, NotificationEventMonitor,
					},
				},
			},
		},

		// other notification channel types
		{
			resp: map[string]interface{}{
				"channels": []interface{}{
					map[string]interface{}{
						"id":   "ch-foobar",
						"name": "notification-test",
						"type": "line",
					},
				},
			},
			want: []NotificationChannel{
				&NotificationChannelBase{
					ID:   "ch-foobar",
					Name: "notification-test",
					Type: NotificationChannelType("line"),
				},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("FindNotificationChannel-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/channels" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/channels", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode(tc.resp)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			c := &Client{
				BaseURL:    u,
				HTTPClient: ts.Client(),
			}
			got, err := c.FindNotificationChannels(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("FindNotificationChannels differs: (-got +want)\n%s", diff)
			}
		})
	}
}
