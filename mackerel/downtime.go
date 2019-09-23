package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Downtime is a downtime.
type Downtime struct {
	ID                   string              `json:"id,omitempty"`
	Name                 string              `json:"name"`
	Memo                 string              `json:"memo,omitempty"`
	Start                Timestamp           `json:"start"`
	Duration             int64               `json:"duration"`
	Recurrence           *DowntimeRecurrence `json:"recurrence,omitempty"`
	ServiceScopes        []string            `json:"serviceScopes,omitempty"`
	ServiceExcludeScopes []string            `json:"serviceExcludeScopes,omitempty"`
	RoleScopes           []string            `json:"roleScopes,omitempty"`
	RoleExcludeScopes    []string            `json:"roleExcludeScopes,omitempty"`
	MonitorScopes        []string            `json:"monitorScopes,omitempty"`
	MonitorExcludeScopes []string            `json:"monitorExcludeScopes,omitempty"`
}

// DowntimeRecurrenceType is a type of DowntimeRecurrence
type DowntimeRecurrenceType string

func (t DowntimeRecurrenceType) String() string {
	return string(t)
}

const (
	// DowntimeRecurrenceTypeHourly is a recurrence hourly.
	DowntimeRecurrenceTypeHourly DowntimeRecurrenceType = "hourly"

	// DowntimeRecurrenceTypeDaily is a recurrence daily.
	DowntimeRecurrenceTypeDaily DowntimeRecurrenceType = "daily"

	// DowntimeRecurrenceTypeWeekly is a recurrence weekly.
	DowntimeRecurrenceTypeWeekly DowntimeRecurrenceType = "weekly"

	// DowntimeRecurrenceTypeMonthly is a recurrence monthly.
	DowntimeRecurrenceTypeMonthly DowntimeRecurrenceType = "monthly"

	// DowntimeRecurrenceTypeYearly is a recurrence yearly.
	DowntimeRecurrenceTypeYearly DowntimeRecurrenceType = "yearly"
)

// DowntimeRecurrence is recurrence settings for a downtime.
type DowntimeRecurrence struct {
	Type     DowntimeRecurrenceType `json:"type,omitempty"`
	Interval int64                  `json:"interval,omitempty"`
	Weekdays []DowntimeWeekday      `json:"weekdays,omitempty"`
	Until    *Timestamp             `json:"until,omitempty"`
}

// DowntimeWeekday specifies a day of the week (Sunday = 0, ...)
type DowntimeWeekday time.Weekday

var days = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

// ParseDowntimeWeekday parses DowntimeWeekday.
func ParseDowntimeWeekday(name string) (DowntimeWeekday, error) {
	for i, s := range days {
		if s == name {
			return DowntimeWeekday(i), nil
		}
	}
	return 0, fmt.Errorf("unknown weekday: %s", name)
}

func (d DowntimeWeekday) String() string {
	return time.Weekday(d).String()
}

// Weekday converts d to time.Weekday.
func (d DowntimeWeekday) Weekday() time.Weekday {
	return time.Weekday(d)
}

// UnmarshalJSON implements json.Unmarshaler
func (d *DowntimeWeekday) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("type error: %s", s)
	}
	ret, err := ParseDowntimeWeekday(s[1 : len(s)-1])
	if err != nil {
		return fmt.Errorf("unknown downtime weekday: %s", s)
	}
	*d = ret
	return nil
}

// MarshalJSON implements json.Marshaler
func (d DowntimeWeekday) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 10)
	buf = append(buf, '"')
	s := d.String()
	buf = append(buf, s...)
	buf = append(buf, '"')
	return buf, nil
}

// FindDowntimes finds downtimes
func (c *Client) FindDowntimes(ctx context.Context) ([]*Downtime, error) {
	var ret []*Downtime
	_, err := c.do(ctx, http.MethodGet, "/api/v0/downtimes", nil, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CreateDowntime creates new downtime.
func (c *Client) CreateDowntime(ctx context.Context, param *Downtime) (*Downtime, error) {
	var ret Downtime
	_, err := c.do(ctx, http.MethodPost, "/api/v0/downtimes", param, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// UpdateDowntime updates a downtime.
func (c *Client) UpdateDowntime(ctx context.Context, downtimeID string, param *Downtime) (*Downtime, error) {
	var ret Downtime
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/downtimes/%s", downtimeID), param, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// DeleteDowntime deletes a downtime.
func (c *Client) DeleteDowntime(ctx context.Context, downtimeID string) (*Downtime, error) {
	var ret Downtime
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/downtimes/%s", downtimeID), nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
