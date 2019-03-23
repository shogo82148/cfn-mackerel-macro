package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Monitor represents interface to which each monitor type must confirm to.
type Monitor interface {
	MonitorType() MonitorType
	MonitorID() string
	MonitorName() string

	isMonitor()
}

// MonitorType is a type of monitors.
type MonitorType string

const (
	// MonitorTypeConnectivity is a type for host connectivity monitoring.
	MonitorTypeConnectivity MonitorType = "connectivity"

	// MonitorTypeHostMeric is a type for Host metric monitoring.
	MonitorTypeHostMeric MonitorType = "host"

	// MonitorTypeServiceMetric is a type for Service metric monitoring.
	MonitorTypeServiceMetric MonitorType = "service"

	// MonitorTypeExternalHTTP is type for External monitoring.
	MonitorTypeExternalHTTP MonitorType = "external"

	// MonitorTypeExpression is a type for Expression monitoring.
	MonitorTypeExpression MonitorType = "expression"
)

func (t MonitorType) String() string {
	return string(t)
}

// Ensure each monitor type conforms to the Monitor interface.
var (
	_ Monitor = (*MonitorConnectivity)(nil)
	_ Monitor = (*MonitorHostMetric)(nil)
	_ Monitor = (*MonitorServiceMetric)(nil)
	_ Monitor = (*MonitorExternalHTTP)(nil)
	_ Monitor = (*MonitorExpression)(nil)
)

// Ensure only monitor types defined in this package can be assigned to the
// Monitor interface.
func (m *MonitorConnectivity) isMonitor()  {}
func (m *MonitorHostMetric) isMonitor()    {}
func (m *MonitorServiceMetric) isMonitor() {}
func (m *MonitorExternalHTTP) isMonitor()  {}
func (m *MonitorExpression) isMonitor()    {}

// MonitorConnectivity represents connectivity monitor.
type MonitorConnectivity struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorConnectivity) MonitorType() MonitorType { return MonitorTypeConnectivity }

// MonitorName returns monitor name.
func (m *MonitorConnectivity) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorConnectivity) MonitorID() string { return m.ID }

// MonitorHostMetric represents host metric monitor.
type MonitorHostMetric struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Metric           string   `json:"metric,omitempty"`
	Operator         string   `json:"operator,omitempty"`
	Warning          *float64 `json:"warning"`
	Critical         *float64 `json:"critical"`
	Duration         uint64   `json:"duration,omitempty"`
	MaxCheckAttempts uint64   `json:"maxCheckAttempts,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorHostMetric) MonitorType() MonitorType { return MonitorTypeHostMeric }

// MonitorName returns monitor name.
func (m *MonitorHostMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorHostMetric) MonitorID() string { return m.ID }

// MonitorServiceMetric represents service metric monitor.
type MonitorServiceMetric struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Service          string   `json:"service,omitempty"`
	Metric           string   `json:"metric,omitempty"`
	Operator         string   `json:"operator,omitempty"`
	Warning          *float64 `json:"warning"`
	Critical         *float64 `json:"critical"`
	Duration         uint64   `json:"duration,omitempty"`
	MaxCheckAttempts uint64   `json:"maxCheckAttempts,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorServiceMetric) MonitorType() MonitorType { return MonitorTypeServiceMetric }

// MonitorName returns monitor name.
func (m *MonitorServiceMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorServiceMetric) MonitorID() string { return m.ID }

// MonitorExternalHTTP represents external HTTP monitor.
type MonitorExternalHTTP struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Method                          string   `json:"method,omitempty"`
	URL                             string   `json:"url,omitempty"`
	MaxCheckAttempts                uint64   `json:"maxCheckAttempts,omitempty"`
	Service                         string   `json:"service,omitempty"`
	ResponseTimeCritical            *float64 `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning             *float64 `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration            *uint64  `json:"responseTimeDuration,omitempty"`
	RequestBody                     string   `json:"requestBody,omitempty"`
	ContainsString                  string   `json:"containsString,omitempty"`
	CertificationExpirationCritical *uint64  `json:"certificationExpirationCritical,omitempty"`
	CertificationExpirationWarning  *uint64  `json:"certificationExpirationWarning,omitempty"`
	SkipCertificateVerification     bool     `json:"skipCertificateVerification,omitempty"`
	// Empty list of headers and nil are different. You have to specify empty
	// list as headers explicitly if you want to remove all headers instead of
	// using nil.
	Headers []HeaderField `json:"headers"`
}

// HeaderField represents key-value pairs in an HTTP header for external http
// monitoring.
type HeaderField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MonitorType returns monitor type.
func (m *MonitorExternalHTTP) MonitorType() MonitorType { return MonitorTypeExternalHTTP }

// MonitorName returns monitor name.
func (m *MonitorExternalHTTP) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExternalHTTP) MonitorID() string { return m.ID }

// MonitorExpression represents expression monitor.
type MonitorExpression struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Expression string   `json:"expression,omitempty"`
	Operator   string   `json:"operator,omitempty"`
	Warning    *float64 `json:"warning"`
	Critical   *float64 `json:"critical"`
}

// MonitorType returns monitor type.
func (m *MonitorExpression) MonitorType() MonitorType { return MonitorTypeExpression }

// MonitorName returns monitor name.
func (m *MonitorExpression) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExpression) MonitorID() string { return m.ID }

// CreateMonitor creates a new monitoring.
func (c *Client) CreateMonitor(ctx context.Context, param Monitor) (Monitor, error) {
	var resp json.RawMessage
	_, err := c.do(ctx, http.MethodPost, "/api/v0/monitors", param, &resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(resp)
}

// UpdateMonitor updates a monitoring.
func (c *Client) UpdateMonitor(ctx context.Context, monitorID string, param Monitor) (Monitor, error) {
	var resp json.RawMessage
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/monitors/%s", monitorID), param, &resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(resp)
}

// DeleteMonitor deletes a monitoring.
func (c *Client) DeleteMonitor(ctx context.Context, monitorID string) (Monitor, error) {
	var resp json.RawMessage
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/monitors/%s", monitorID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(resp)
}

func decodeMonitor(mes json.RawMessage) (Monitor, error) {
	var typeData struct {
		Type MonitorType `json:"type"`
	}
	if err := json.Unmarshal(mes, &typeData); err != nil {
		return nil, err
	}
	var m Monitor
	switch typeData.Type {
	case MonitorTypeConnectivity:
		m = &MonitorConnectivity{}
	case MonitorTypeHostMeric:
		m = &MonitorHostMetric{}
	case MonitorTypeServiceMetric:
		m = &MonitorServiceMetric{}
	case MonitorTypeExternalHTTP:
		m = &MonitorExternalHTTP{}
	case MonitorTypeExpression:
		m = &MonitorExpression{}
	}
	if err := json.Unmarshal(mes, m); err != nil {
		return nil, err
	}
	return m, nil
}
