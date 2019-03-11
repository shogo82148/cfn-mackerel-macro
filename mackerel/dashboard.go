package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// Dashboard is a dashborad.
type Dashboard struct {
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title,omitempty"`
	Memo      string   `json:"memo,omitempty"`
	URLPath   string   `json:"urlPath,omitempty"`
	Widgets   []Widget `json:"widgets,omitempty"`
	CreatedAt int64    `json:"createdAt,omitempty"`
	UpdatedAt int64    `json:"updatedAt,omitempty"`
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

// GraphType is a type of a graph widget.
type GraphType string

const (
	// GraphTypeHost is a host metric graph.
	GraphTypeHost GraphType = "host"
)

// Graph is a graph definition of a graph widget.
type Graph interface {
	GraphType() GraphType
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
