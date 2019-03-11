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
		return "", nil, err
	}
	id, err := d.Function.parseDashboardID(ctx, d.Event.PhysicalResourceID)
	if err != nil {
		return "", nil, err
	}
	_, err = c.UpdateDashboard(ctx, id, param)
	if err != nil {
		return "", nil, err
	}

	return d.Event.PhysicalResourceID, map[string]interface{}{}, nil
}

func (d *dashboard) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.Dashboard, error) {
	var dp dproxy.Drain
	in := dproxy.New(properties)
	param := &mackerel.Dashboard{
		Title:   dp.String(in.M("Title")),
		Memo:    dp.String(dproxy.Default(in.M("Memo"), "")),
		URLPath: dp.String(in.M("UrlPath")),
		Widgets: []mackerel.Widget{
			&mackerel.WidgetGraph{
				Type:  mackerel.WidgetTypeGraph,
				Title: "foobar",
				Graph: &mackerel.GraphHost{
					Type:   mackerel.GraphTypeHost,
					HostID: "3yAYEDLXKL5",
					Name:   "foobar",
				},
				Layout: &mackerel.Layout{
					X:      0,
					Y:      0,
					Width:  24,
					Height: 6,
				},
			},
		},
	}
	if err := dp.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

func (d *dashboard) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	id, err := d.Function.parseDashboardID(ctx, d.Event.PhysicalResourceID)
	if err != nil {
		return "", nil, err
	}

	c := d.Function.getclient()
	_, err = c.DeleteDashboard(ctx, id)
	if err != nil {
		return "", nil, err
	}

	return d.Event.PhysicalResourceID, map[string]interface{}{}, nil
}
