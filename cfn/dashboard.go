package cfn

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type dashboard struct {
	Function *Function
	Event    cfn.Event
}

func (d *dashboard) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := d.Function.getclient()
	param, err := d.convertToParam(ctx, d.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateDashboard(ctx, param)
	if err != nil {
		return "", nil, err
	}

	id, err := d.Function.buildDashboardID(ctx, ret.ID)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{}, nil
}

func (d *dashboard) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := d.Function.getclient()
	param, err := d.convertToParam(ctx, d.Event.ResourceProperties)
	if err != nil {
		return d.Event.PhysicalResourceID, nil, err
	}
	id, err := d.Function.parseDashboardID(ctx, d.Event.PhysicalResourceID)
	if err != nil {
		return d.Event.PhysicalResourceID, nil, err
	}
	_, err = c.UpdateDashboard(ctx, id, param)
	if err != nil {
		return d.Event.PhysicalResourceID, nil, err
	}

	return d.Event.PhysicalResourceID, map[string]interface{}{}, nil
}

func (d *dashboard) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.Dashboard, error) {
	var dp dproxy.Drain
	in := dproxy.New(properties)

	widgets := []mackerel.Widget{}
	for _, w := range dp.ProxyArray(in.M("Widgets").ProxySet()) {
		widgets = append(widgets, d.convertWidget(ctx, &dp, w))
	}

	param := &mackerel.Dashboard{
		Title:   dp.String(in.M("Title")),
		Memo:    dp.String(dproxy.Default(in.M("Memo"), "")),
		URLPath: dp.String(in.M("UrlPath")),
		Widgets: widgets,
	}
	if err := dp.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

func (d *dashboard) convertWidget(ctx context.Context, dp *dproxy.Drain, properties dproxy.Proxy) mackerel.Widget {
	typ, err := properties.M("Type").String()
	if err != nil {
		dp.Put(err)
		return nil
	}
	switch typ {
	case mackerel.WidgetTypeGraph.String():
		return &mackerel.WidgetGraph{
			Type:   mackerel.WidgetTypeGraph,
			Title:  dp.String(dproxy.Default(properties.M("Title"), "")),
			Graph:  d.convertGraph(ctx, dp, properties.M("Graph")),
			Layout: d.convertLayout(dp, properties.M("Layout")),
		}
	}
	return nil
}

func (d *dashboard) convertGraph(ctx context.Context, dp *dproxy.Drain, properties dproxy.Proxy) mackerel.Graph {
	typ, err := properties.M("Type").String()
	if err != nil {
		dp.Put(err)
		return nil
	}
	switch typ {
	case mackerel.GraphTypeHost.String():
		id, err := properties.M("Host").String()
		dp.Put(err)
		hostID, err := d.Function.parseHostID(ctx, id)
		dp.Put(err)
		return &mackerel.GraphHost{
			HostID: hostID,
			Name:   dp.String(properties.M("Name")),
		}
	case mackerel.GraphTypeRole.String():
		id, err := properties.M("Role").String()
		dp.Put(err)
		serviceName, roleName, err := d.Function.parseRoleID(ctx, id)
		dp.Put(err)
		return &mackerel.GraphRole{
			RoleFullname: serviceName + ":" + roleName,
			Name:         dp.String(properties.M("Name")),
			IsStacked:    dp.Bool(dproxy.Default(properties.M("IsStacked"), false)),
		}
	case mackerel.GraphTypeService.String():
		id, err := properties.M("Service").String()
		dp.Put(err)
		serviceName, err := d.Function.parseServiceID(ctx, id)
		dp.Put(err)
		return &mackerel.GraphService{
			ServiceName: serviceName,
			Name:        dp.String(properties.M("Name")),
		}
	case mackerel.GraphTypeExpression.String():
		return &mackerel.GraphExpression{
			Expression: dp.String(properties.M("Expression")),
		}
	}
	return nil
}

func (d *dashboard) convertLayout(dp *dproxy.Drain, properties dproxy.Proxy) *mackerel.Layout {
	return &mackerel.Layout{
		X:      dp.Uint64(properties.M("X")),
		Y:      dp.Uint64(properties.M("Y")),
		Width:  dp.Uint64(properties.M("Width")),
		Height: dp.Uint64(properties.M("Height")),
	}
}

func (d *dashboard) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	id, err := d.Function.parseDashboardID(ctx, d.Event.PhysicalResourceID)
	if err != nil {
		return d.Event.PhysicalResourceID, nil, err
	}

	c := d.Function.getclient()
	_, err = c.DeleteDashboard(ctx, id)
	if err != nil {
		return d.Event.PhysicalResourceID, nil, err
	}

	return d.Event.PhysicalResourceID, map[string]interface{}{}, nil
}
