package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// NotificationGroup is a notification group.
type NotificationGroup struct {
	ID                        string                     `json:"id,omitempty"`
	Name                      string                     `json:"name"`
	NotificationLevel         NotificationLevel          `json:"notificationLevel"`
	ChildNotificationGroupIDs []string                   `json:"childNotificationGroupIds"`
	ChildChannelIDs           []string                   `json:"childChannelIds"`
	Monitors                  []NotificationGroupMonitor `json:"monitors,omitempty"`
	Services                  []NotificationGroupService `json:"services,omitempty"`
}

// NotificationLevel is the notification level.
type NotificationLevel string

const (
	// NotificationLevelAll receives all notifications.
	NotificationLevelAll NotificationLevel = "all"

	// NotificationLevelCritical receives critical notifications.
	NotificationLevelCritical NotificationLevel = "critical"
)

func (level NotificationLevel) String() string {
	return string(level)
}

// NotificationGroupMonitor is a monitor setting for a notification group.
type NotificationGroupMonitor struct {
	ID          string `json:"id"`
	SkipDefault bool   `json:"skipDefault"`
}

// NotificationGroupService is a service setting for a notification group.
type NotificationGroupService struct {
	Name string `json:"name"`
}

// FindNotificationGroups returns the list of notification groups.
func (c *Client) FindNotificationGroups(ctx context.Context) ([]*NotificationGroup, error) {
	var data struct {
		NotificationGroups []*NotificationGroup `json:"notificationGroups"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/notification-groups", nil, &data)
	if err != nil {
		return nil, err
	}
	return data.NotificationGroups, nil
}

// CreateNotificationGroup creates a new notification group.
func (c *Client) CreateNotificationGroup(ctx context.Context, group *NotificationGroup) (*NotificationGroup, error) {
	// fill required fields
	in := *group
	if in.ChildNotificationGroupIDs == nil {
		in.ChildNotificationGroupIDs = []string{}
	}
	if in.ChildChannelIDs == nil {
		in.ChildChannelIDs = []string{}
	}

	var data NotificationGroup
	_, err := c.do(ctx, http.MethodPost, "/api/v0/notification-groups", in, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateNotificationGroup creates a new notification group.
func (c *Client) UpdateNotificationGroup(ctx context.Context, groupID string, group *NotificationGroup) (*NotificationGroup, error) {
	// fill required fields
	in := *group
	if in.ChildNotificationGroupIDs == nil {
		in.ChildNotificationGroupIDs = []string{}
	}
	if in.ChildChannelIDs == nil {
		in.ChildChannelIDs = []string{}
	}

	var data NotificationGroup
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/notification-groups/%s", groupID), in, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteNotificationGroup deletes a notification group.
func (c *Client) DeleteNotificationGroup(ctx context.Context, groupID string) (*NotificationGroup, error) {
	var data NotificationGroup
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/notification-groups/%s", groupID), nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
