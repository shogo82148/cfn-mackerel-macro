package cfn

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
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
	name, ok := r.Event.ResourceProperties["Name"].(string)
	if !ok {
		return "", nil, errors.New("Name is missing")
	}
	service, ok := r.Event.ResourceProperties["Service"].(string)
	if !ok {
		return "", nil, errors.New("Service is missing")
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
	name, ok := r.Event.ResourceProperties["Name"].(string)
	if !ok {
		return "", nil, errors.New("Name is missing")
	}
	oldName, ok := r.Event.OldResourceProperties["Name"].(string)
	if !ok {
		return "", nil, errors.New("Name is missing")
	}
	service, ok := r.Event.OldResourceProperties["Service"].(string)
	if !ok {
		return "", nil, errors.New("Service is missing")
	}
	oldService, ok := r.Event.OldResourceProperties["Service"].(string)
	if !ok {
		return "", nil, errors.New("Service is missing")
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
