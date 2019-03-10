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
	c := h.Function.getclient()
	param, err := h.convertToParam(ctx, h.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	hostID, err := c.CreateHost(ctx, param)
	if err != nil {
		return "", nil, err
	}

	id, err := h.Function.buildID(ctx, "host", hostID)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Name": param.Name,
	}, nil
}

func (h *host) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := h.Function.getclient()
	param, err := h.convertToParam(ctx, h.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	_, id, err := h.Function.parseID(ctx, h.Event.PhysicalResourceID, 1)
	if err != nil {
		return "", nil, err
	}
	_, err = c.UpdateHost(ctx, id[0], (*mackerel.UpdateHostParam)(param))
	if err != nil {
		return "", nil, err
	}

	return h.Event.PhysicalResourceID, map[string]interface{}{
		"Name": param.Name,
	}, nil
}

func (h *host) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.CreateHostParam, error) {
	var param mackerel.CreateHostParam
	var d dproxy.Drain
	in := dproxy.New(properties)
	param.Name = d.String(in.M("Name"))
	roles := d.Array(in.M("Roles"))
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}

	for _, r := range roles {
		id, err := dproxy.New(r).String()
		if err != nil {
			return nil, err
		}
		typ, name, err := h.Function.parseID(ctx, id, 2)
		if err != nil {
			return nil, err
		}
		if typ != "role" {
			return nil, fmt.Errorf("invalid type: %s", typ)
		}
		param.RoleFullnames = append(param.RoleFullnames, name[0]+":"+name[1])
	}

	return &param, nil
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
