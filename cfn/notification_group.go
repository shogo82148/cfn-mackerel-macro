package cfn

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type notificationGroup struct {
	Function *Function
	Event    cfn.Event
}

func (g *notificationGroup) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := g.Function.getclient()
	param, err := g.convertToParam(ctx, g.Event.ResourceProperties)
	if err != nil {
		return
	}
	ret, err := c.CreateNotificationGroup(ctx, param)
	if err != nil {
		return
	}
	data = map[string]interface{}{
		"Name": ret.Name,
	}
	physicalResourceID, err = g.Function.buildNotificationGroupID(ctx, ret.ID)
	return
}

func (g *notificationGroup) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.NotificationGroup, error) {
	var d dproxy.Drain
	in := dproxy.New(properties)
	param := &mackerel.NotificationGroup{
		Name:              d.String(in.M("Name")),
		NotificationLevel: mackerel.NotificationLevel(d.String(in.M("NotificationLevel"))),
	}
	for _, id := range d.StringArray(in.M("ChildNotificationGroupIds").ProxySet()) {
		groupID, err := g.Function.parseNotificationChannelID(ctx, id)
		if err != nil {
			d.Put(err)
			continue
		}
		param.ChildNotificationGroupIDs = append(param.ChildNotificationGroupIDs, groupID)
	}
	for _, id := range d.StringArray(in.M("ChildChannelIds").ProxySet()) {
		channelID, err := g.Function.parseNotificationChannelID(ctx, id)
		if err != nil {
			d.Put(err)
			continue
		}
		param.ChildChannelIDs = append(param.ChildChannelIDs, channelID)
	}
	for _, id := range d.StringArray(in.M("Services").ProxySet()) {
		serviceName, err := g.Function.parseServiceID(ctx, id)
		if err != nil {
			d.Put(err)
			continue
		}
		param.Services = append(param.Services, serviceName)
	}
	for _, m := range d.ProxyArray(in.M("Monitors").ProxySet()) {
		var monitorID string
		id, err := m.M("Id").String()
		if err == nil {
			monitorID, err = g.Function.parseMonitorID(ctx, id)
			d.Put(err)
		}
		param.Monitors = append(param.Monitors, mackerel.NotificationGroupMonitor{
			ID:          monitorID,
			SkipDefault: d.Bool(m.M("SkipDefault")),
		})
	}
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

func (g *notificationGroup) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}

func (g *notificationGroup) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}
