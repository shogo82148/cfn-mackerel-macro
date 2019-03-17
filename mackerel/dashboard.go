package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Dashboard is a dashborad.
type Dashboard struct {
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Memo      string `json:"memo,omitempty"`
	URLPath   string `json:"urlPath,omitempty"`
	Widgets   []Widget
	CreatedAt int64 `json:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty"`
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
)

func (t WidgetType) String() string {
	return string(t)
}

// Widget is a widget.
type Widget interface {
	WidgetType() WidgetType
	WidgetTitle() string
	WidgetLayout() *Layout
}

// Layout describes the layout of the widget.
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

// UnmarshalJSON unmashal JSON.
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

// GraphType is a type of a graph widget.
type GraphType string

func (t GraphType) String() string {
	return string(t)
}

const (
	// GraphTypeHost is a host metric graph.
	GraphTypeHost GraphType = "host"
)

// Graph is a graph definition of a graph widget.
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

// FindDashboards finds dashboards.
func (c *Client) FindDashboards(ctx context.Context) ([]*Dashboard, error) {
	ret := []*Dashboard{}
	err := c.do(ctx, http.MethodGet, "/api/v0/dashboards", nil, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// FindDashboard find the dashboard.
func (c *Client) FindDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), nil, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CreateDashboard creates a new dashboard.
func (c *Client) CreateDashboard(ctx context.Context, param *Dashboard) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodPost, "/api/v0/dashboards", param, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateDashboard deletes a service
func (c *Client) UpdateDashboard(ctx context.Context, dashboardID string, param *Dashboard) (*Dashboard, error) {
	ret := &Dashboard{}
	err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), param, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteDashboard deletes a service
func (c *Client) DeleteDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	dashboard := &Dashboard{}
	err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), nil, dashboard)
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}
