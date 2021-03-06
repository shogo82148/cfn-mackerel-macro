package cfn

import (
	"context"
	"log"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type host struct {
	Function *Function
	Event    cfn.Event
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

	meta := getmetadata(h.Event)
	if err := c.PutHostMetaData(ctx, hostID, "cloudformation", meta); err != nil {
		return hostID, nil, err
	}

	id, err := h.Function.buildHostID(ctx, hostID)
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
		return h.Event.PhysicalResourceID, nil, err
	}
	id, err := h.Function.parseHostID(ctx, h.Event.PhysicalResourceID)
	if err != nil {
		return h.Event.PhysicalResourceID, nil, err
	}
	_, err = c.UpdateHost(ctx, id, (*mackerel.UpdateHostParam)(param))
	if err != nil {
		return h.Event.PhysicalResourceID, nil, err
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
		serviceName, roleName, err := h.Function.parseRoleID(ctx, id)
		if err != nil {
			return nil, err
		}
		param.RoleFullnames = append(param.RoleFullnames, serviceName+":"+roleName)
	}

	return &param, nil
}

func (h *host) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = h.Event.PhysicalResourceID
	id, err := h.Function.parseHostID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as host id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := h.Function.getclient()
	err = c.RetireHost(ctx, id)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
