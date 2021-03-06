package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// NotificationChannel represents interface to which each notification type must confirm to.
type NotificationChannel interface {
	json.Marshaler
	json.Unmarshaler
	NotificationChannelType() NotificationChannelType
	NotificationChannelID() string
	NotificationChannelName() string
}

// NotificationChannelType is a type of notification channel.
type NotificationChannelType string

const (
	// NotificationChannelTypeEmail is email type.
	NotificationChannelTypeEmail NotificationChannelType = "email"

	// NotificationChannelTypeSlack is slack type.
	NotificationChannelTypeSlack NotificationChannelType = "slack"

	// NotificationChannelTypeWebHook is web hook type.
	NotificationChannelTypeWebHook NotificationChannelType = "webhook"

	// NotificationChannelTypeLine is LINE type.
	NotificationChannelTypeLine NotificationChannelType = "line"

	// NotificationChannelTypeWebHook is Chatwork type.
	NotificationChannelTypeChatwork NotificationChannelType = "chatwork"

	// NotificationChannelTypeTypetalk is Typetalk type.
	NotificationChannelTypeTypetalk NotificationChannelType = "typetalk"

	// NotificationChannelTypeHipchat is Hipchat type.
	NotificationChannelTypeHipchat NotificationChannelType = "hipchat"

	// NotificationChannelTypeHipchat is twilio type.
	NotificationChannelTypeTwilio NotificationChannelType = "twilio"

	// NotificationChannelTypeReactio is Hipchat type.
	NotificationChannelTypeReactio NotificationChannelType = "reactio"

	// NotificationChannelTypePagerduty is Hipchat type.
	NotificationChannelTypePagerduty NotificationChannelType = "pagerduty"

	// NotificationChannelTypeOpsgenie is Opsgenie type.
	NotificationChannelTypeOpsgenie NotificationChannelType = "opsgenie"

	// NotificationChannelTypeYammer is Yammer type.
	NotificationChannelTypeYammer NotificationChannelType = "yammer"

	// NotificationChannelTypeMicrosoftTeams is Hipchat type.
	NotificationChannelTypeMicrosoftTeams NotificationChannelType = "microsoft-teams"

	// NotificationChannelTypeAmazonEventBridge is Amazon Event Bridge type.
	NotificationChannelTypeAmazonEventBridge NotificationChannelType = "amazon-event-bridge"
)

func (t NotificationChannelType) String() string {
	return string(t)
}

// NotificationEvent is an event type to notify.
type NotificationEvent string

const (
	// NotificationEventAlert notifies alert events.
	NotificationEventAlert NotificationEvent = "alert"

	// NotificationEventAlertGroup notifies alert group events.
	NotificationEventAlertGroup NotificationEvent = "alertGroup"

	// NotificationEventHostStatus notifies host status change events.
	NotificationEventHostStatus NotificationEvent = "hostStatus"

	// NotificationEventHostRegister notifies host registration events.
	NotificationEventHostRegister NotificationEvent = "hostRegister"

	// NotificationEventHostRetire notifies host retirement events.
	NotificationEventHostRetire NotificationEvent = "hostRetire"

	// NotificationEventMonitor notifies monitor events.
	NotificationEventMonitor NotificationEvent = "monitor"
)

func (e NotificationEvent) String() string {
	return string(e)
}

// NotificationChannelBase is base type of notification channel.
type NotificationChannelBase struct {
	Type NotificationChannelType `json:"type"`
	ID   string                  `json:"id,omitempty"`
	Name string                  `json:"name"`
}

// NotificationChannelType returns the type.
func (c *NotificationChannelBase) NotificationChannelType() NotificationChannelType {
	return c.Type
}

// NotificationChannelID returns the id of the channel.
func (c *NotificationChannelBase) NotificationChannelID() string {
	return c.ID
}

// NotificationChannelName returns the name of the channel.
func (c *NotificationChannelBase) NotificationChannelName() string {
	return c.Name
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelBase) UnmarshalJSON(b []byte) error {
	type channel NotificationChannelBase
	var data channel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	*c = NotificationChannelBase(data)
	return nil
}

// MarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelBase) MarshalJSON() ([]byte, error) {
	type channel NotificationChannelBase
	data := *(*channel)(c)
	return json.Marshal(data)
}

// NotificationChannelEmail is an email notification channel.
type NotificationChannelEmail struct {
	Type    NotificationChannelType `json:"type"`
	ID      string                  `json:"id,omitempty"`
	Name    string                  `json:"name"`
	Emails  []string                `json:"emails"`
	UserIDs []string                `json:"userIds"`
	Events  []NotificationEvent     `json:"events"`
}

// NotificationChannelType returns NotificationChannelTypeEmail
func (c *NotificationChannelEmail) NotificationChannelType() NotificationChannelType {
	return NotificationChannelTypeEmail
}

// NotificationChannelID returns the id of the channel.
func (c *NotificationChannelEmail) NotificationChannelID() string {
	return c.ID
}

// NotificationChannelName returns the name of the channel.
func (c *NotificationChannelEmail) NotificationChannelName() string {
	return c.Name
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelEmail) UnmarshalJSON(b []byte) error {
	type channel NotificationChannelEmail
	var data channel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	data.Type = NotificationChannelTypeEmail
	*c = NotificationChannelEmail(data)
	return nil
}

// MarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelEmail) MarshalJSON() ([]byte, error) {
	type channel NotificationChannelEmail
	data := *(*channel)(c)
	data.Type = NotificationChannelTypeEmail
	return json.Marshal(data)
}

// NotificationChannelSlack is a slack notification channel.
type NotificationChannelSlack struct {
	Type              NotificationChannelType          `json:"type"`
	ID                string                           `json:"id,omitempty"`
	Name              string                           `json:"name"`
	URL               string                           `json:"url"`
	EnabledGraphImage bool                             `json:"enabledGraphImage"`
	Mentions          NotificationChannelSlackMentions `json:"mentions"`
	Events            []NotificationEvent              `json:"events"`
}

// NotificationChannelSlackMentions is mentions in slack notification.
type NotificationChannelSlackMentions struct {
	OK       string `json:"ok,omitempty"`
	Warning  string `json:"warning,omitempty"`
	Critical string `json:"critical,omitempty"`
}

// NotificationChannelType returns NotificationChannelTypeSlack
func (c *NotificationChannelSlack) NotificationChannelType() NotificationChannelType {
	return NotificationChannelTypeSlack
}

// NotificationChannelID returns the id of the channel.
func (c *NotificationChannelSlack) NotificationChannelID() string {
	return c.ID
}

// NotificationChannelName returns the name of the channel.
func (c *NotificationChannelSlack) NotificationChannelName() string {
	return c.Name
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelSlack) UnmarshalJSON(b []byte) error {
	type channel NotificationChannelSlack
	var data channel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	data.Type = NotificationChannelTypeSlack
	*c = NotificationChannelSlack(data)
	return nil
}

// MarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelSlack) MarshalJSON() ([]byte, error) {
	type channel NotificationChannelSlack
	data := *(*channel)(c)
	data.Type = NotificationChannelTypeSlack
	return json.Marshal(data)
}

// NotificationChannelWebHook is a web hook notification channel.
type NotificationChannelWebHook struct {
	Type   NotificationChannelType `json:"type"`
	ID     string                  `json:"id,omitempty"`
	Name   string                  `json:"name"`
	URL    string                  `json:"url"`
	Events []NotificationEvent     `json:"events"`
}

// NotificationChannelType returns NotificationChannelTypeWebHook
func (c *NotificationChannelWebHook) NotificationChannelType() NotificationChannelType {
	return NotificationChannelTypeWebHook
}

// NotificationChannelID returns the id of the channel.
func (c *NotificationChannelWebHook) NotificationChannelID() string {
	return c.ID
}

// NotificationChannelName returns the name of the channel.
func (c *NotificationChannelWebHook) NotificationChannelName() string {
	return c.Name
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelWebHook) UnmarshalJSON(b []byte) error {
	type channel NotificationChannelWebHook
	var data channel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	data.Type = NotificationChannelTypeWebHook
	*c = NotificationChannelWebHook(data)
	return nil
}

// MarshalJSON implements json.Unmarshaler.
func (c *NotificationChannelWebHook) MarshalJSON() ([]byte, error) {
	type channel NotificationChannelWebHook
	data := *(*channel)(c)
	data.Type = NotificationChannelTypeWebHook
	return json.Marshal(data)
}

type notificationChannel struct {
	NotificationChannel
}

func (c *notificationChannel) UnmarshalJSON(b []byte) error {
	var data struct {
		Type NotificationChannelType `json:"type"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	switch data.Type {
	case NotificationChannelTypeEmail:
		c.NotificationChannel = &NotificationChannelEmail{}
	case NotificationChannelTypeSlack:
		c.NotificationChannel = &NotificationChannelSlack{}
	case NotificationChannelTypeWebHook:
		c.NotificationChannel = &NotificationChannelWebHook{}
	default:
		c.NotificationChannel = &NotificationChannelBase{}
	}
	return c.NotificationChannel.UnmarshalJSON(b)
}

// FindNotificationChannels returns rhe list of notification channels.
func (c *Client) FindNotificationChannels(ctx context.Context) ([]NotificationChannel, error) {
	var data struct {
		Channels []notificationChannel `json:"channels"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/channels", nil, &data)
	if err != nil {
		return nil, err
	}

	ret := make([]NotificationChannel, 0, len(data.Channels))
	for _, c := range data.Channels {
		ret = append(ret, c.NotificationChannel)
	}
	return ret, nil
}

// CreateNotificationChannel creates a new notification channel.
func (c *Client) CreateNotificationChannel(ctx context.Context, ch NotificationChannel) (NotificationChannel, error) {
	var ret notificationChannel
	_, err := c.do(ctx, http.MethodPost, "/api/v0/channels", ch, &ret)
	if err != nil {
		return nil, err
	}
	return ret.NotificationChannel, nil
}

// DeleteNotificationChannel deletes a notification channel.
func (c *Client) DeleteNotificationChannel(ctx context.Context, channelID string) (NotificationChannel, error) {
	var ret notificationChannel
	_, err := c.do(ctx, http.MethodPost, fmt.Sprintf("/api/v0/channels/%s", channelID), nil, &ret)
	if err != nil {
		return nil, err
	}
	return ret.NotificationChannel, nil
}
