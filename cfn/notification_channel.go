package cfn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"errors"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type notificationChannel struct {
	Function *Function
	Event    cfn.Event
}

func (ch *notificationChannel) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := ch.Function.getclient()
	param, err := ch.convertToParam(ctx, ch.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateNotificationChannel(ctx, param)
	if err != nil {
		return "", nil, err
	}

	id, err := ch.Function.buildNotificationChannelID(ctx, ret.NotificationChannelID())
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Name": ret.NotificationChannelName(),
	}, nil
}

func (ch *notificationChannel) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return ch.create(ctx) // create new one and replace
}

func (ch *notificationChannel) convertToParam(ctx context.Context, properties map[string]interface{}) (mackerel.NotificationChannel, error) {
	var ret mackerel.NotificationChannel
	var d dproxy.Drain
	in := dproxy.New(properties)
	typ := d.String(in.M("Type"))
	switch typ {
	case mackerel.NotificationChannelTypeEmail.String():
		ret = &mackerel.NotificationChannelEmail{
			Name:    d.String(in.M("Name")),
			Emails:  d.StringArray(in.M("Emails").ProxySet()),
			UserIDs: []string{}, // TODO: support user ids
			Events:  ch.convertEvents(&d, in.M("Events")),
		}
	case mackerel.NotificationChannelTypeSlack.String():
		ret = &mackerel.NotificationChannelSlack{
			Name: d.String(in.M("Name")),
			URL:  d.String(in.M("Url")),
			Mentions: mackerel.NotificationChannelSlackMentions{
				OK:       d.String(dproxy.Default(in.M("Mentions").M("Ok"), "")),
				Warning:  d.String(dproxy.Default(in.M("Mentions").M("Warning"), "")),
				Critical: d.String(dproxy.Default(in.M("Mentions").M("Critical"), "")),
			},
			EnabledGraphImage: d.Bool(dproxy.Default(in.M("EnabledGraphImage"), false)),
			Events:            ch.convertEvents(&d, in.M("Events")),
		}
	case mackerel.NotificationChannelTypeWebHook.String():
		ret = &mackerel.NotificationChannelWebHook{
			Name:   d.String(in.M("Name")),
			URL:    d.String(in.M("Url")),
			Events: ch.convertEvents(&d, in.M("Events")),
		}
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (ch *notificationChannel) convertEvents(d *dproxy.Drain, in dproxy.Proxy) []mackerel.NotificationEvent {
	events := d.StringArray(in.ProxySet())
	ret := make([]mackerel.NotificationEvent, 0, len(events))
	for _, e := range events {
		ret = append(ret, mackerel.NotificationEvent(e))
	}
	return ret
}

func (ch *notificationChannel) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = ch.Event.PhysicalResourceID
	id, err := ch.Function.parseNotificationChannelID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as monitor id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := ch.Function.getclient()
	_, err = c.DeleteNotificationChannel(ctx, id)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
