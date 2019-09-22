package cfn

import (
	"context"
	"errors"
	"log"
	"net/http"

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
		Name: d.String(in.M("Name")),
		NotificationLevel: mackerel.NotificationLevel(
			d.String(dproxy.Default(in.M("NotificationLevel"), mackerel.NotificationLevelAll.String())),
		),
	}
	for _, id := range d.StringArray(dproxy.Default(in.M("ChildNotificationGroupIds"), []interface{}{}).ProxySet()) {
		groupID, err := g.Function.parseNotificationGroupID(ctx, id)
		d.Put(err)
		if err != nil {
			continue
		}
		param.ChildNotificationGroupIDs = append(param.ChildNotificationGroupIDs, groupID)
	}
	for _, id := range d.StringArray(dproxy.Default(in.M("ChildChannelIds"), []interface{}{}).ProxySet()) {
		channelID, err := g.Function.parseNotificationChannelID(ctx, id)
		d.Put(err)
		if err != nil {
			continue
		}
		param.ChildChannelIDs = append(param.ChildChannelIDs, channelID)
	}
	for _, m := range d.ProxyArray(dproxy.Default(in.M("Services"), []interface{}{}).ProxySet()) {
		var serviceName string
		id, err := m.M("Id").String()
		d.Put(err)
		if err == nil {
			serviceName, err = g.Function.parseServiceID(ctx, id)
			d.Put(err)
		}
		param.Services = append(param.Services, mackerel.NotificationGroupService{
			Name: serviceName,
		})
	}
	for _, m := range d.ProxyArray(dproxy.Default(in.M("Monitors"), []interface{}{}).ProxySet()) {
		var monitorID string
		id, err := m.M("Id").String()
		d.Put(err)
		if err == nil {
			monitorID, err = g.Function.parseMonitorID(ctx, id)
			d.Put(err)
		}
		param.Monitors = append(param.Monitors, mackerel.NotificationGroupMonitor{
			ID:          monitorID,
			SkipDefault: d.Bool(dproxy.Default(m.M("SkipDefault"), false)),
		})
	}
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

func (g *notificationGroup) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = g.Event.PhysicalResourceID
	groupID, err := g.Function.parseNotificationGroupID(ctx, physicalResourceID)
	if err != nil {
		return
	}
	c := g.Function.getclient()
	param, err := g.convertToParam(ctx, g.Event.ResourceProperties)
	if err != nil {
		return
	}
	ret, err := c.UpdateNotificationGroup(ctx, groupID, param)
	if err != nil {
		return
	}
	data = map[string]interface{}{
		"Name": ret.Name,
	}
	return
}

func (g *notificationGroup) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = g.Event.PhysicalResourceID
	groupID, err := g.Function.parseNotificationGroupID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as notification group id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := g.Function.getclient()
	_, err = c.DeleteNotificationGroup(ctx, groupID)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
