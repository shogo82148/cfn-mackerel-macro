package cfn

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type role struct {
	Function *Function
	Event    cfn.Event
}

func (r *role) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(r.Event.ResourceProperties)
	name := d.String(in.M("Name"))
	service := d.String(in.M("Service"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	serviceName, err := r.Function.parseServiceID(ctx, service)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to parse Service ID: %s", err)
	}

	c := r.Function.getclient()
	ss, err := c.CreateRole(ctx, serviceName, &mackerel.CreateRoleParam{
		Name: name,
	})
	if err != nil {
		return "", nil, err
	}

	id, err := r.Function.buildRoleID(ctx, serviceName, ss.Name)
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
	serviceName, roleName, err := r.Function.parseRoleID(ctx, r.Event.PhysicalResourceID)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}

	c := r.Function.getclient()
	ss, err := c.DeleteRole(ctx, serviceName, roleName)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}

	return r.Event.PhysicalResourceID, map[string]interface{}{
		"Name": ss.Name,
	}, nil
}
