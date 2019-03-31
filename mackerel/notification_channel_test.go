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

func TestCreateNotificationChannel(t *testing.T) {
	tests := []struct {
		in   NotificationChannel
		want map[string]interface{}
	}{
		{
			in: &NotificationChannelEmail{
				Name:    "notification-test",
				Emails:  []string{"my.address@example.com"},
				UserIDs: []string{"userId"},
				Events:  []NotificationEvent{NotificationEventAlert},
			},
			want: map[string]interface{}{
				"name":    "notification-test",
				"type":    "email",
				"emails":  []interface{}{"my.address@example.com"},
				"userIds": []interface{}{"userId"},
				"events":  []interface{}{"alert"},
			},
		},
		{
			in: &NotificationChannelSlack{
				Name: "notification-test",
				URL:  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
				Mentions: NotificationChannelSlackMentions{
					OK:      "ok message",
					Warning: "warning message",
				},
				EnabledGraphImage: true,
				Events:            []NotificationEvent{NotificationEventAlert},
			},
			want: map[string]interface{}{
				"name": "notification-test",
				"type": "slack",
				"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
				"mentions": map[string]interface{}{
					"ok":      "ok message",
					"warning": "warning message",
				},
				"enabledGraphImage": true,
				"events":            []interface{}{"alert"},
			},
		},
		{
			in: &NotificationChannelWebHook{
				Name:   "notification-test",
				URL:    "https://example.com/webhook",
				Events: []NotificationEvent{NotificationEventAlert},
			},
			want: map[string]interface{}{
				"name":   "notification-test",
				"type":   "webhook",
				"url":    "https://example.com/webhook",
				"events": []interface{}{"alert"},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("CreateNotificationChannel-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method, want %s, got %s", http.MethodPost, r.Method)
				}
				var data map[string]interface{}
				dec := json.NewDecoder(r.Body)
				if err := dec.Decode(&data); err != nil {
					t.Error(err)
					return
				}
				if diff := cmp.Diff(data, tc.want); diff != "" {
					t.Errorf("CreateNotification differs: (-got +want)\n%s", diff)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				data["id"] = "channelId"
				enc := json.NewEncoder(w)
				if err := enc.Encode(data); err != nil {
					t.Error(err)
				}
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

			ch, err := c.CreateNotificationChannel(context.Background(), tc.in)
			if err != nil {
				t.Error(err)
				return
			}
			if ch.NotificationChannelID() != "channelId" {
				t.Errorf("unexpected channel id: want %s, got %s", "channelId", ch.NotificationChannelID())
			}
		})
	}
}
