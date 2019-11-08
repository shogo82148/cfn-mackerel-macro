package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Monitor represents interface to which each monitor type must confirm to.
type Monitor interface {
	json.Unmarshaler
	json.Marshaler
	MonitorType() MonitorType
	MonitorID() string
	MonitorName() string
}

// MonitorType is a type of monitors.
type MonitorType string

const (
	// MonitorTypeConnectivity is a type for host connectivity monitoring.
	MonitorTypeConnectivity MonitorType = "connectivity"

	// MonitorTypeHostMetric is a type for Host metric monitoring.
	MonitorTypeHostMetric MonitorType = "host"

	// MonitorTypeServiceMetric is a type for Service metric monitoring.
	MonitorTypeServiceMetric MonitorType = "service"

	// MonitorTypeExternalHTTP is a type for External monitoring.
	MonitorTypeExternalHTTP MonitorType = "external"

	// MonitorTypeExpression is a type for Expression monitoring.
	MonitorTypeExpression MonitorType = "expression"

	// MonitorTypeAnomalyDetection is a type for anomaly detection.
	MonitorTypeAnomalyDetection MonitorType = "anomalyDetection"
)

func (t MonitorType) String() string {
	return string(t)
}

type monitor struct {
	Monitor
}

func (m *monitor) UnmarshalJSON(b []byte) error {
	var data struct {
		Type MonitorType `json:"type"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	switch data.Type {
	case MonitorTypeConnectivity:
		m.Monitor = &MonitorConnectivity{}
	case MonitorTypeHostMetric:
		m.Monitor = &MonitorHostMetric{}
	case MonitorTypeServiceMetric:
		m.Monitor = &MonitorServiceMetric{}
	case MonitorTypeExternalHTTP:
		m.Monitor = &MonitorExternalHTTP{}
	case MonitorTypeExpression:
		m.Monitor = &MonitorExpression{}
	case MonitorTypeAnomalyDetection:
		m.Monitor = &MonitorAnomalyDetection{}
	default:
		return fmt.Errorf("unknown monitor type: %s", data.Type)
	}
	return json.Unmarshal(b, m.Monitor)
}

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

// UnmarshalJSON implements json.Unmarshaler.
func (m *MonitorConnectivity) UnmarshalJSON(b []byte) error {
	type monitor MonitorConnectivity
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeConnectivity
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m *MonitorConnectivity) MarshalJSON() ([]byte, error) {
	type monitor MonitorConnectivity
	data := (*monitor)(m)
	data.Type = MonitorTypeConnectivity
	return json.Marshal(data)
}

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
func (m *MonitorHostMetric) MonitorType() MonitorType { return MonitorTypeHostMetric }

// MonitorName returns monitor name.
func (m *MonitorHostMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorHostMetric) MonitorID() string { return m.ID }

// UnmarshalJSON implements json.Unmarshal.
func (m *MonitorHostMetric) UnmarshalJSON(b []byte) error {
	type monitor MonitorHostMetric
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeHostMetric
	return nil
}

// MarshalJSON implements json.Marshal.
func (m *MonitorHostMetric) MarshalJSON() ([]byte, error) {
	type monitor MonitorHostMetric
	data := (*monitor)(m)
	data.Type = MonitorTypeHostMetric
	return json.Marshal(data)
}

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

	MissingDurationWarning  *uint64 `json:"missingDurationWarning"`
	MissingDurationCritical *uint64 `json:"missingDurationCritical"`
}

// MonitorType returns monitor type.
func (m *MonitorServiceMetric) MonitorType() MonitorType { return MonitorTypeServiceMetric }

// MonitorName returns monitor name.
func (m *MonitorServiceMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorServiceMetric) MonitorID() string { return m.ID }

// UnmarshalJSON implements json.Unmarshaler.
func (m *MonitorServiceMetric) UnmarshalJSON(b []byte) error {
	type monitor MonitorServiceMetric
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeServiceMetric
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m *MonitorServiceMetric) MarshalJSON() ([]byte, error) {
	type monitor MonitorServiceMetric
	data := (*monitor)(m)
	data.Type = MonitorTypeServiceMetric
	return json.Marshal(data)
}

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

// UnmarshalJSON implements json.Unmarshaler.
func (m *MonitorExternalHTTP) UnmarshalJSON(b []byte) error {
	type monitor MonitorExternalHTTP
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeExternalHTTP
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m *MonitorExternalHTTP) MarshalJSON() ([]byte, error) {
	type monitor MonitorExternalHTTP
	data := (*monitor)(m)
	data.Type = MonitorTypeExternalHTTP
	return json.Marshal(data)
}

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

// UnmarshalJSON implements json.Unmarshaler.
func (m *MonitorExpression) UnmarshalJSON(b []byte) error {
	type monitor MonitorExpression
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeExpression
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m *MonitorExpression) MarshalJSON() ([]byte, error) {
	type monitor MonitorExpression
	data := (*monitor)(m)
	data.Type = MonitorTypeExpression
	return json.Marshal(data)
}

// AnomalyDetectionSensitivityType is a type for sensitivity of anomaly detection monitorings.
type AnomalyDetectionSensitivityType string

const (
	// AnomalyDetectionSensitivityInsensitive is insensitive
	AnomalyDetectionSensitivityInsensitive AnomalyDetectionSensitivityType = "insensitive"

	// AnomalyDetectionSensitivityNormal is normal
	AnomalyDetectionSensitivityNormal AnomalyDetectionSensitivityType = "normal"

	// AnomalyDetectionSensitivitySensitive is sensitive
	AnomalyDetectionSensitivitySensitive AnomalyDetectionSensitivityType = "sensitive"
)

// MonitorAnomalyDetection represents anomaly detection monitor.
type MonitorAnomalyDetection struct {
	ID                   string      `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Memo                 string      `json:"memo,omitempty"`
	Type                 MonitorType `json:"type,omitempty"`
	IsMute               bool        `json:"isMute,omitempty"`
	NotificationInterval uint64      `json:"notificationInterval,omitempty"`

	Scopes              []string                        `json:"scopes"`
	WarningSensitivity  AnomalyDetectionSensitivityType `json:"warningSensitivity,omitempty"`
	CriticalSensitivity AnomalyDetectionSensitivityType `json:"criticalSensitivity,omitempty"`
	MaxCheckAttempts    uint64                          `json:"maxCheckAttempts,omitempty"`
	TrainingPeriodFrom  Timestamp                       `json:"trainingPeriodFrom,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorAnomalyDetection) MonitorType() MonitorType { return MonitorTypeAnomalyDetection }

// MonitorName returns monitor name.
func (m *MonitorAnomalyDetection) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorAnomalyDetection) MonitorID() string { return m.ID }

// UnmarshalJSON implements json.Unmarshaler.
func (m *MonitorAnomalyDetection) UnmarshalJSON(b []byte) error {
	type monitor MonitorAnomalyDetection
	data := (*monitor)(m)
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	m.Type = MonitorTypeAnomalyDetection
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m *MonitorAnomalyDetection) MarshalJSON() ([]byte, error) {
	type monitor MonitorAnomalyDetection
	data := (*monitor)(m)
	data.Type = MonitorTypeAnomalyDetection
	return json.Marshal(data)
}

// FindMonitors returns monitoring settings.
func (c *Client) FindMonitors(ctx context.Context) ([]Monitor, error) {
	var resp []monitor
	_, err := c.do(ctx, http.MethodGet, "/api/v0/monitors", nil, &resp)
	if err != nil {
		return nil, err
	}

	ret := make([]Monitor, 0, len(resp))
	for _, m := range resp {
		ret = append(ret, m.Monitor)
	}
	return ret, nil
}

// FindMonitor returns a monitoring setting.
func (c *Client) FindMonitor(ctx context.Context, monitorID string) (Monitor, error) {
	var resp monitor
	_, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/monitors/%s", monitorID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}

// CreateMonitor creates a new monitoring.
func (c *Client) CreateMonitor(ctx context.Context, param Monitor) (Monitor, error) {
	var resp monitor
	_, err := c.do(ctx, http.MethodPost, "/api/v0/monitors", param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}

// UpdateMonitor updates a monitoring.
func (c *Client) UpdateMonitor(ctx context.Context, monitorID string, param Monitor) (Monitor, error) {
	var resp monitor
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/monitors/%s", monitorID), param, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}

// DeleteMonitor deletes a monitoring.
func (c *Client) DeleteMonitor(ctx context.Context, monitorID string) (Monitor, error) {
	var resp monitor
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/monitors/%s", monitorID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}
