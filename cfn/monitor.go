package cfn

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/koron/go-dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type monitor struct {
	Function *Function
	Event    cfn.Event
}

func (m *monitor) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := m.Function.getclient()
	mm, err := m.convertToParam(ctx, m.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateMonitor(ctx, mm)
	if err != nil {
		return "", nil, err
	}

	id, err := m.Function.buildMonitorID(ctx, ret.MonitorID())
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Type": ret.MonitorType().String(),
		"Name": ret.MonitorName(),
	}, nil
}

func (m *monitor) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	id, err := m.Function.parseMonitorID(ctx, m.Event.PhysicalResourceID)
	if err != nil {
		return "", nil, err
	}

	c := m.Function.getclient()
	mm, err := m.convertToParam(ctx, m.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.UpdateMonitor(ctx, id, mm)
	if err != nil {
		return "", nil, err
	}

	return id, map[string]interface{}{
		"Type": ret.MonitorType().String(),
		"Name": ret.MonitorName(),
	}, nil
}

func (m *monitor) convertToParam(ctx context.Context, properties map[string]interface{}) (mackerel.Monitor, error) {
	in := dproxy.New(properties)
	typ, err := in.M("Type").String()
	if err != nil {
		return nil, err
	}

	var d dproxy.Drain
	var mm mackerel.Monitor
	switch typ {
	case mackerel.MonitorTypeConnectivity.String():
		mm = &mackerel.MonitorConnectivity{
			Type: mackerel.MonitorTypeConnectivity,
			Name: d.String(in.M("Name")),
		}
	default:
		return nil, fmt.Errorf("unknown monitor type: %s", typ)
	}
	return mm, nil
}

func (m *monitor) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	monitorID, err := m.Function.parseMonitorID(ctx, m.Event.PhysicalResourceID)
	if err != nil {
		return "", nil, err
	}

	c := m.Function.getclient()
	_, err = c.DeleteMonitor(ctx, monitorID)
	if err != nil {
		return "", nil, err
	}

	return m.Event.PhysicalResourceID, map[string]interface{}{}, nil
}
