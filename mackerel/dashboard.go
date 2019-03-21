package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Dashboard is a dashborad.
// https://mackerel.io/api-docs/entry/dashboards#create
type Dashboard struct {
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title,omitempty"`
	Memo      string   `json:"memo,omitempty"`
	URLPath   string   `json:"urlPath,omitempty"`
	Widgets   []Widget `json:"-"`
	CreatedAt int64    `json:"createdAt,omitempty"`
	UpdatedAt int64    `json:"updatedAt,omitempty"`
}

// MarshalJSON mashal JSON.
func (d *Dashboard) MarshalJSON() ([]byte, error) {
	type dashboard Dashboard
	var data struct {
		Widgets []widget `json:"widgets,omitempty"`
		dashboard
	}

	data.dashboard = dashboard(*d)
	data.Widgets = make([]widget, len(d.Widgets))
	for i, w := range d.Widgets {
		data.Widgets[i] = widget{w}
	}
	return json.Marshal(data)
}

// UnmarshalJSON unmashals JSON.
func (d *Dashboard) UnmarshalJSON(b []byte) error {
	type dashboard Dashboard
	var data struct {
		Widgets []widget `json:"widgets,omitempty"`
		dashboard
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	*d = Dashboard(data.dashboard)
	d.Widgets = make([]Widget, len(data.Widgets))
	for i, w := range data.Widgets {
		d.Widgets[i] = w.Widget
	}
	return nil
}

// WidgetType is a type of dashboard widget.
type WidgetType string

const (
	// WidgetTypeGraph is a graph widget.
	WidgetTypeGraph WidgetType = "graph"

	// WidgetTypeValue is a value widget.
	WidgetTypeValue WidgetType = "value"

	// WidgetTypeMarkdown is a markdown widget.
	WidgetTypeMarkdown WidgetType = "markdown"
)

func (t WidgetType) String() string {
	return string(t)
}

// Widget is a widget.
// https://mackerel.io/api-docs/entry/dashboards#widget
type Widget interface {
	json.Marshaler
	json.Unmarshaler
	WidgetType() WidgetType
	WidgetTitle() string
	WidgetLayout() *Layout
}

// Layout describes the layout of the widget.
// https://mackerel.io/api-docs/entry/dashboards#layout
type Layout struct {
	X      uint64 `json:"x"`
	Y      uint64 `json:"y"`
	Width  uint64 `json:"width"`
	Height uint64 `json:"height"`
}

type widget struct {
	Widget
}

func (w *widget) UnmarshalJSON(b []byte) error {
	var data struct {
		Type WidgetType `json:"type"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	switch data.Type {
	case WidgetTypeGraph:
		w.Widget = &WidgetGraph{}
	case WidgetTypeValue:
		w.Widget = &WidgetValue{}
	case WidgetTypeMarkdown:
		w.Widget = &WidgetMarkdown{}
	default:
		return fmt.Errorf("unknown widget type: %s", data.Type)
	}
	return json.Unmarshal(b, &w.Widget)
}

// WidgetGraph is a graph widget.
type WidgetGraph struct {
	Type   WidgetType `json:"type"`
	Title  string     `json:"title"`
	Graph  Graph      `json:"graph,omitempty"`
	Layout *Layout    `json:"layout,omitempty"`
}

var _ Widget = (*WidgetGraph)(nil)

// WidgetType returns WidgetTypeGraph.
func (w *WidgetGraph) WidgetType() WidgetType { return WidgetTypeGraph }

// WidgetTitle returns the title of the widget.
func (w *WidgetGraph) WidgetTitle() string { return w.Title }

// WidgetLayout returns the layout of the widget.
func (w *WidgetGraph) WidgetLayout() *Layout { return w.Layout }

// MarshalJSON implements the json.Marshaler.
func (w *WidgetGraph) MarshalJSON() ([]byte, error) {
	type widgetGraph WidgetGraph
	data := *(*widgetGraph)(w)
	data.Type = WidgetTypeGraph
	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *WidgetGraph) UnmarshalJSON(b []byte) error {
	// wrap Graph with graph type to use custom UnmarshalJSON func
	type widgetGraph WidgetGraph
	data := (*widgetGraph)(w)
	g := &graph{}
	data.Graph = g

	if err := json.Unmarshal(b, data); err != nil {
		return err
	}

	w.Type = WidgetTypeGraph
	w.Graph = g.Graph // unwrap Graph
	return nil
}

// WidgetValue is a value widget.
type WidgetValue struct {
	Type   WidgetType `json:"type"`
	Title  string     `json:"title"`
	Metric Metric     `json:"metric,omitempty"`
	Layout *Layout    `json:"layout,omitempty"`
}

var _ Widget = (*WidgetValue)(nil)

// WidgetType returns WidgetTypeValue.
func (w *WidgetValue) WidgetType() WidgetType { return WidgetTypeValue }

// WidgetTitle returns the title of the widget.
func (w *WidgetValue) WidgetTitle() string { return w.Title }

// WidgetLayout returns the layout of the widget.
func (w *WidgetValue) WidgetLayout() *Layout { return w.Layout }

// MarshalJSON implements the json.Marshaler.
func (w *WidgetValue) MarshalJSON() ([]byte, error) {
	type widgetValue WidgetValue
	data := *(*widgetValue)(w)
	data.Type = WidgetTypeValue
	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *WidgetValue) UnmarshalJSON(b []byte) error {
	// wrap Metric with metric type to use custom UnmarshalJSON func
	type widgetValue WidgetValue
	data := (*widgetValue)(w)
	m := &metric{}
	data.Metric = m

	if err := json.Unmarshal(b, data); err != nil {
		return err
	}

	w.Type = WidgetTypeValue
	w.Metric = m.Metric // unwrap metric
	return nil
}

// WidgetMarkdown is a markdown widget for dashboards.
type WidgetMarkdown struct {
	Type     WidgetType `json:"type"`
	Title    string     `json:"title"`
	Markdown string     `json:"markdown,omitempty"`
	Layout   *Layout    `json:"layout,omitempty"`
}

var _ Widget = (*WidgetMarkdown)(nil)

// WidgetType returns WidgetTypeMarkdown.
func (w *WidgetMarkdown) WidgetType() WidgetType { return WidgetTypeMarkdown }

// WidgetTitle returns the title of the widget.
func (w *WidgetMarkdown) WidgetTitle() string { return w.Title }

// WidgetLayout returns the layout of the widget.
func (w *WidgetMarkdown) WidgetLayout() *Layout { return w.Layout }

// MarshalJSON implements the json.Marshaler.
func (w *WidgetMarkdown) MarshalJSON() ([]byte, error) {
	type widgetMarkdown WidgetMarkdown
	data := *(*widgetMarkdown)(w)
	data.Type = WidgetTypeMarkdown
	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *WidgetMarkdown) UnmarshalJSON(b []byte) error {
	// wrap Metric with metric type to use custom UnmarshalJSON func
	type widgetMarkdown WidgetMarkdown
	data := (*widgetMarkdown)(w)

	if err := json.Unmarshal(b, data); err != nil {
		return err
	}

	w.Type = WidgetTypeMarkdown
	return nil
}

// GraphType is a type of a graph widget.
type GraphType string

func (t GraphType) String() string {
	return string(t)
}

const (
	// GraphTypeHost is a host metric graph.
	GraphTypeHost GraphType = "host"

	// GraphTypeRole is a role metric graph.
	GraphTypeRole GraphType = "role"

	// GraphTypeService is a service metric graph.
	GraphTypeService GraphType = "service"

	// GraphTypeExpression is an expression graph.
	GraphTypeExpression GraphType = "expression"

	// GraphTypeUnknown is an unknown graph.
	GraphTypeUnknown GraphType = "unknown"
)

// Graph is a graph definition of a graph widget.
// https://mackerel.io/api-docs/entry/dashboards#graph
type Graph interface {
	GraphType() GraphType
}

type graph struct {
	Graph
}

func (g *graph) UnmarshalJSON(b []byte) error {
	var data struct {
		Type GraphType `json:"type"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	switch data.Type {
	case GraphTypeHost:
		g.Graph = &GraphHost{}
	case GraphTypeRole:
		g.Graph = &GraphRole{}
	case GraphTypeService:
		g.Graph = &GraphService{}
	case GraphTypeExpression:
		g.Graph = &GraphExpression{}
	case GraphTypeUnknown:
		g.Graph = &GraphUnknown{}
	default:
		return fmt.Errorf("unknown widget type: %s", data.Type)
	}
	return json.Unmarshal(b, &g.Graph)
}

// GraphHost is a host metric graph.
type GraphHost struct {
	Type   GraphType `json:"type"`
	HostID string    `json:"hostId"`
	Name   string    `json:"name"`
}

var _ Graph = (*GraphHost)(nil)

// GraphType returns GraphTypeHost.
func (g *GraphHost) GraphType() GraphType { return GraphTypeHost }

// MarshalJSON implements json.Marshaler.
func (g *GraphHost) MarshalJSON() ([]byte, error) {
	type graph GraphHost
	data := *(*graph)(g)
	data.Type = g.GraphType()
	return json.Marshal(data)
}

// GraphRole is a role metric graph.
type GraphRole struct {
	Type         GraphType `json:"type"`
	RoleFullname string    `json:"roleFullname"`
	Name         string    `json:"name"`
	IsStacked    bool      `json:"isStacked,omitempty"`
}

var _ Graph = (*GraphRole)(nil)

// GraphType returns GraphTypeRole.
func (g *GraphRole) GraphType() GraphType { return GraphTypeRole }

// MarshalJSON implements json.Marshaler.
func (g *GraphRole) MarshalJSON() ([]byte, error) {
	type graph GraphRole
	data := *(*graph)(g)
	data.Type = g.GraphType()
	return json.Marshal(data)
}

// GraphService is a service metric graph.
type GraphService struct {
	Type        GraphType `json:"type"`
	ServiceName string    `json:"serviceName"`
	Name        string    `json:"name"`
}

var _ Graph = (*GraphService)(nil)

// GraphType returns GraphTypeService.
func (g *GraphService) GraphType() GraphType { return GraphTypeService }

// MarshalJSON implements json.Marshaler.
func (g *GraphService) MarshalJSON() ([]byte, error) {
	type graph GraphService
	data := *(*graph)(g)
	data.Type = g.GraphType()
	return json.Marshal(data)
}

// GraphExpression is an expression metric graph.
type GraphExpression struct {
	Type       GraphType `json:"type"`
	Expression string    `json:"expression"`
}

var _ Graph = (*GraphExpression)(nil)

// GraphType returns GraphTypeExpression.
func (g *GraphExpression) GraphType() GraphType { return GraphTypeExpression }

// MarshalJSON implements json.Marshaler.
func (g *GraphExpression) MarshalJSON() ([]byte, error) {
	type graph GraphExpression
	data := *(*graph)(g)
	data.Type = g.GraphType()
	return json.Marshal(data)
}

// GraphUnknown is an unknown graph.
type GraphUnknown struct {
	Type GraphType `json:"type"`
}

// GraphType returns GraphTypeUnknown.
func (g *GraphUnknown) GraphType() GraphType { return GraphTypeUnknown }

// MarshalJSON implements json.Marshaler.
func (g *GraphUnknown) MarshalJSON() ([]byte, error) {
	type graph GraphUnknown
	data := *(*graph)(g)
	data.Type = g.GraphType()
	return json.Marshal(data)
}

// MetricType is a type of a metric.
type MetricType string

const (
	// MetricTypeHost is a host metric.
	MetricTypeHost MetricType = "host"

	// MetricTypeService is a service metric graph.
	MetricTypeService MetricType = "service"

	// MetricTypeExpression is an expression graph.
	MetricTypeExpression MetricType = "expression"

	// MetricTypeUnknown is an unknown graph.
	MetricTypeUnknown MetricType = "unknown"
)

func (t MetricType) String() string {
	return string(t)
}

// Metric is a metric of metric widget.
// https://mackerel.io/api-docs/entry/dashboards#metric
type Metric interface {
	MetricType() MetricType
}

type metric struct {
	Metric
}

func (m *metric) UnmarshalJSON(b []byte) error {
	var data struct {
		Type MetricType `json:"type"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	switch data.Type {
	case MetricTypeHost:
		m.Metric = &MetricHost{}
	case MetricTypeService:
		m.Metric = &MetricService{}
	case MetricTypeExpression:
		m.Metric = &MetricExpression{}
	case MetricTypeUnknown:
		m.Metric = &MetricUnknown{}
	default:
		return fmt.Errorf("unknown metric type: %s", data.Type)
	}
	return json.Unmarshal(b, &m.Metric)
}

// MetricHost is a host metric.
type MetricHost struct {
	Type   MetricType `json:"type"`
	HostID string     `json:"hostId"`
	Name   string     `json:"name"`
}

var _ Metric = (*MetricHost)(nil)

// MetricType returns MetricTypeHost.
func (m *MetricHost) MetricType() MetricType { return MetricTypeHost }

// MarshalJSON implements json.Marshaler.
func (m *MetricHost) MarshalJSON() ([]byte, error) {
	type metric MetricHost
	data := *(*metric)(m)
	data.Type = m.MetricType()
	return json.Marshal(data)
}

// MetricService is a service metric.
type MetricService struct {
	Type        MetricType `json:"type"`
	ServiceName string     `json:"serviceName"`
	Name        string     `json:"name"`
}

var _ Metric = (*MetricService)(nil)

// MetricType returns MetricTypeService.
func (m *MetricService) MetricType() MetricType { return MetricTypeService }

// MarshalJSON implements json.Marshaler.
func (m *MetricService) MarshalJSON() ([]byte, error) {
	type metric MetricService
	data := *(*metric)(m)
	data.Type = m.MetricType()
	return json.Marshal(data)
}

// MetricExpression is an expression metric Metric.
type MetricExpression struct {
	Type       MetricType `json:"type"`
	Expression string     `json:"expression"`
}

var _ Metric = (*MetricExpression)(nil)

// MetricType returns MetricTypeExpression.
func (m *MetricExpression) MetricType() MetricType { return MetricTypeExpression }

// MarshalJSON implements json.Marshaler.
func (m *MetricExpression) MarshalJSON() ([]byte, error) {
	type metric MetricExpression
	data := *(*metric)(m)
	data.Type = m.MetricType()
	return json.Marshal(data)
}

// MetricUnknown is an unknown Metric.
type MetricUnknown struct {
	Type MetricType `json:"type"`
}

// MetricType returns MetricTypeUnknown.
func (m *MetricUnknown) MetricType() MetricType { return MetricTypeUnknown }

// MarshalJSON implements json.Marshaler.
func (m *MetricUnknown) MarshalJSON() ([]byte, error) {
	type metric MetricUnknown
	data := *(*metric)(m)
	data.Type = m.MetricType()
	return json.Marshal(data)
}

// FindDashboards finds dashboards.
// https://mackerel.io/api-docs/entry/dashboards#list
func (c *Client) FindDashboards(ctx context.Context) ([]*Dashboard, error) {
	ret := []*Dashboard{}
	err := c.do(ctx, http.MethodGet, "/api/v0/dashboards", nil, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// FindDashboard find the dashboard.
// https://mackerel.io/api-docs/entry/dashboards#get
func (c *Client) FindDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CreateDashboard creates a new dashboard.
// https://mackerel.io/api-docs/entry/dashboards#create
func (c *Client) CreateDashboard(ctx context.Context, param *Dashboard) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodPost, "/api/v0/dashboards", param, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateDashboard deletes a service
// https://mackerel.io/api-docs/entry/dashboards#update
func (c *Client) UpdateDashboard(ctx context.Context, dashboardID string, param *Dashboard) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), param, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteDashboard deletes a service
// https://mackerel.io/api-docs/entry/dashboards#delete
func (c *Client) DeleteDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	dashboard := &Dashboard{}
	err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), nil, dashboard)
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}
