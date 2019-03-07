package cfn

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/koron/go-dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type role struct {
	Function *Function
	Event    cfn.Event
}

func (r *role) handle(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	switch r.Event.RequestType {
	case cfn.RequestCreate:
		return r.create(ctx)
	case cfn.RequestUpdate:
		return r.update(ctx)
	case cfn.RequestDelete:
		return r.delete(ctx)
	}
	return "", nil, fmt.Errorf("unknown request type: %s", r.Event.RequestType)
}

func (r *role) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(r.Event.ResourceProperties)
	name := d.String(in.M("Name"))
	service := d.String(in.M("Service"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	typ, serviceName, err := r.Function.parseID(ctx, service, 1)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to parse Service ID: %s", err)
	}
	if typ != "service" {
		return "", nil, fmt.Errorf("Invlid type for Service: %s", service)
	}

	c := r.Function.getclient()
	ss, err := c.CreateRole(ctx, serviceName[0], &mackerel.CreateRoleParam{
		Name: name,
	})
	if err != nil {
		return "", nil, err
	}

	id, err := r.Function.buildID(ctx, "role", serviceName[0], ss.Name)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Name": ss.Name,
	}, nil
}

func (r *role) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(r.Event.ResourceProperties)
	old := dproxy.New(r.Event.OldResourceProperties)

	name := d.String(in.M("Name"))
	service := d.String(in.M("Service"))
	oldName := d.String(old.M("Name"))
	oldService := d.String(old.M("Service"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	if name == oldName && service == oldService {
		// No update is needed.
		return r.Event.PhysicalResourceID, map[string]interface{}{
			"Name": name,
		}, nil
	}

	// need to create a new role.
	return r.create(ctx)
}

func (r *role) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	_, id, err := r.Function.parseID(ctx, r.Event.PhysicalResourceID, 1)
	if err != nil {
		return "", nil, err
	}

	c := r.Function.getclient()
	ss, err := c.DeleteRole(ctx, id[0], id[1])
	if err != nil {
		return "", nil, err
	}

	return r.Event.PhysicalResourceID, map[string]interface{}{
		"Name": ss.Name,
	}, nil
}
