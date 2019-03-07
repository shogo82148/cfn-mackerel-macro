package cfn

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/koron/go-dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type host struct {
	Function *Function
	Event    cfn.Event
}

func (h *host) handle(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	switch h.Event.RequestType {
	case cfn.RequestCreate:
		return h.create(ctx)
	case cfn.RequestUpdate:
		return h.update(ctx)
	case cfn.RequestDelete:
		return h.delete(ctx)
	}
	return "", nil, fmt.Errorf("unknown request type: %s", h.Event.RequestType)
}

func (h *host) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(h.Event.ResourceProperties)
	name := d.String(in.M("Name"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	c := h.Function.getclient()
	hostID, err := c.CreateHost(ctx, &mackerel.CreateHostParam{
		Name: name,
		// TODO: memo
	})
	if err != nil {
		return "", nil, err
	}

	id, err := h.Function.buildID(ctx, "host", hostID)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Name": name,
	}, nil
}

func (h *host) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(h.Event.ResourceProperties)
	old := dproxy.New(h.Event.OldResourceProperties)

	name := d.String(in.M("Name"))
	oldName := d.String(old.M("Name"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	// TODO: update information
	_ = oldName

	return h.Event.PhysicalResourceID, map[string]interface{}{
		"Name": name,
	}, nil
}

func (h *host) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	typ, id, err := h.Function.parseID(ctx, h.Event.PhysicalResourceID, 1)
	if err != nil {
		return "", nil, err
	}
	if typ != "host" {
		return "", nil, fmt.Errorf("invlid resource type: %s", typ)
	}

	c := h.Function.getclient()
	err = c.RetireHost(ctx, id[0])
	if err != nil {
		return "", nil, err
	}

	return h.Event.PhysicalResourceID, map[string]interface{}{}, nil
}
