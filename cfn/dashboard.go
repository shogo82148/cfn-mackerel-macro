package cfn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type dashboard struct {
	Function *Function
	Event    cfn.Event
}

func (r *dashboard) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	param, err := r.convertToParam(ctx, r.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateDashboard(ctx, param)
	if err != nil {
		return "", nil, err
	}

	id, err := r.Function.buildDashboardID(ctx, ret.ID)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{}, nil
}

func (r *dashboard) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	param, err := r.convertToParam(ctx, r.Event.ResourceProperties)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}
	id, err := r.Function.parseDashboardID(ctx, r.Event.PhysicalResourceID)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}
	_, err = c.UpdateDashboard(ctx, id, param)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}

	return r.Event.PhysicalResourceID, map[string]interface{}{}, nil
}

func (r *dashboard) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.Dashboard, error) {
	var d dproxy.Drain
	in := dproxy.New(properties)

	widgets := []mackerel.Widget{}
	for _, w := range d.ProxyArray(in.M("Widgets").ProxySet()) {
		widgets = append(widgets, r.convertWidget(ctx, &d, w))
	}

	param := &mackerel.Dashboard{
		Title:   d.String(in.M("Title")),
		Memo:    d.String(dproxy.Default(in.M("Memo"), "")),
		URLPath: d.String(in.M("UrlPath")),
		Widgets: widgets,
	}
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

func (r *dashboard) convertWidget(ctx context.Context, d *dproxy.Drain, properties dproxy.Proxy) mackerel.Widget {
	typ, err := properties.M("Type").String()
	if err != nil {
		d.Put(err)
		return nil
	}
	switch typ {
	case mackerel.WidgetTypeGraph.String():
		return &mackerel.WidgetGraph{
			Title:  d.String(dproxy.Default(properties.M("Title"), "")),
			Graph:  r.convertGraph(ctx, d, properties.M("Graph")),
			Range:  r.convertRange(ctx, d, properties.M("Range")),
			Layout: r.convertLayout(d, properties.M("Layout")),
		}
	case mackerel.WidgetTypeValue.String():
		return &mackerel.WidgetValue{
			Title:  d.String(dproxy.Default(properties.M("Title"), "")),
			Metric: r.convertMetric(ctx, d, properties.M("Metric")),
			Layout: r.convertLayout(d, properties.M("Layout")),
		}
	case mackerel.WidgetTypeMarkdown.String():
		return &mackerel.WidgetMarkdown{
			Title:    d.String(dproxy.Default(properties.M("Title"), "")),
			Markdown: d.String(dproxy.Default(properties.M("Markdown"), "")),
			Layout:   r.convertLayout(d, properties.M("Layout")),
		}
	case mackerel.WidgetTypeAlertStatus.String():
		id, err := properties.M("Role").String()
		d.Put(err)
		serviceName, roleName, err := r.Function.parseRoleID(ctx, id)
		d.Put(err)
		roleFullname := serviceName + ":" + roleName
		return &mackerel.WidgetAlertStatus{
			Title:        d.String(dproxy.Default(properties.M("Title"), "")),
			RoleFullname: &roleFullname,
			Layout:       r.convertLayout(d, properties.M("Layout")),
		}
	}
	d.Put(fmt.Errorf("unknown widget type: %s", typ))
	return nil
}

func (r *dashboard) convertGraph(ctx context.Context, d *dproxy.Drain, properties dproxy.Proxy) mackerel.Graph {
	typ, err := properties.M("Type").String()
	if err != nil {
		d.Put(err)
		return nil
	}
	switch typ {
	case mackerel.GraphTypeHost.String():
		id, err := properties.M("Host").String()
		d.Put(err)
		hostID, err := r.Function.parseHostID(ctx, id)
		d.Put(err)
		return &mackerel.GraphHost{
			HostID: hostID,
			Name:   d.String(properties.M("Name")),
		}
	case mackerel.GraphTypeRole.String():
		id, err := properties.M("Role").String()
		d.Put(err)
		serviceName, roleName, err := r.Function.parseRoleID(ctx, id)
		d.Put(err)
		return &mackerel.GraphRole{
			RoleFullname: serviceName + ":" + roleName,
			Name:         d.String(properties.M("Name")),
			IsStacked:    d.Bool(dproxy.Default(properties.M("IsStacked"), false)),
		}
	case mackerel.GraphTypeService.String():
		id, err := properties.M("Service").String()
		d.Put(err)
		serviceName, err := r.Function.parseServiceID(ctx, id)
		d.Put(err)
		return &mackerel.GraphService{
			ServiceName: serviceName,
			Name:        d.String(properties.M("Name")),
		}
	case mackerel.GraphTypeExpression.String():
		return &mackerel.GraphExpression{
			Expression: d.String(properties.M("Expression")),
		}
	}
	d.Put(fmt.Errorf("unknown graph type: %s", typ))
	return nil
}

func (r *dashboard) convertMetric(ctx context.Context, d *dproxy.Drain, properties dproxy.Proxy) mackerel.Metric {
	typ, err := properties.M("Type").String()
	if err != nil {
		d.Put(err)
		return nil
	}
	switch typ {
	case mackerel.MetricTypeHost.String():
		id, err := properties.M("Host").String()
		d.Put(err)
		hostID, err := r.Function.parseHostID(ctx, id)
		d.Put(err)
		return &mackerel.MetricHost{
			HostID: hostID,
			Name:   d.String(properties.M("Name")),
		}
	case mackerel.MetricTypeService.String():
		id, err := properties.M("Service").String()
		d.Put(err)
		serviceName, err := r.Function.parseServiceID(ctx, id)
		d.Put(err)
		return &mackerel.MetricService{
			ServiceName: serviceName,
			Name:        d.String(properties.M("Name")),
		}
	case mackerel.MetricTypeExpression.String():
		return &mackerel.MetricExpression{
			Expression: d.String(properties.M("Expression")),
		}
	}
	d.Put(fmt.Errorf("unknown metric type: %s", typ))
	return nil
}

func (r *dashboard) convertRange(ctx context.Context, d *dproxy.Drain, properties dproxy.Proxy) mackerel.GraphRange {
	if dproxy.IsError(properties, dproxy.ErrorCodeNotFound) {
		return nil
	}
	typ, err := properties.M("Type").String()
	if err != nil {
		d.Put(err)
		return nil
	}
	switch typ {
	case mackerel.GraphRangeTypeRelative.String():
		return &mackerel.GraphRangeRelative{
			Period: d.Int64(properties.M("Period")),
			Offset: d.Int64(properties.M("Offset")),
		}
	case mackerel.GraphRangeTypeAbsolute.String():
		return &mackerel.GraphRangeAbsolute{
			Start: mackerel.Timestamp(d.Int64(properties.M("Start"))),
			End:   mackerel.Timestamp(d.Int64(properties.M("End"))),
		}
	}
	d.Put(fmt.Errorf("unknown graph range type: %s", typ))
	return nil
}

func (r *dashboard) convertLayout(d *dproxy.Drain, properties dproxy.Proxy) *mackerel.Layout {
	return &mackerel.Layout{
		X:      d.Uint64(properties.M("X")),
		Y:      d.Uint64(properties.M("Y")),
		Width:  d.Uint64(properties.M("Width")),
		Height: d.Uint64(properties.M("Height")),
	}
}

func (r *dashboard) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = r.Event.PhysicalResourceID
	id, err := r.Function.parseDashboardID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as dashboard id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := r.Function.getclient()
	_, err = c.DeleteDashboard(ctx, id)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
