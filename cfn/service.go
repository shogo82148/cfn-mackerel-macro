package cfn

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/koron/go-dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type service struct {
	Function *Function
	Event    cfn.Event
}

func (s *service) handle(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	switch s.Event.RequestType {
	case cfn.RequestCreate:
		return s.create(ctx)
	case cfn.RequestUpdate:
		return s.update(ctx)
	case cfn.RequestDelete:
		return s.delete(ctx)
	}
	return "", nil, fmt.Errorf("unknown request type: %s", s.Event.RequestType)
}

func (s *service) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	in := dproxy.New(s.Event.ResourceProperties)
	name, err := in.M("Name").String()
	if err != nil {
		return "", nil, err
	}

	c := s.Function.getclient()
	ss, err := c.CreateService(ctx, &mackerel.CreateServiceParam{
		Name: name,
		// TODO: memo
	})
	if err != nil {
		return "", nil, err
	}

	id, err := s.Function.buildID(ctx, "service", ss.Name)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Name": ss.Name,
	}, nil
}

func (s *service) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(s.Event.ResourceProperties)
	old := dproxy.New(s.Event.OldResourceProperties)

	name := d.String(in.M("Name"))
	oldName := d.String(old.M("Name"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	if name == oldName {
		// No update is needed.
		return s.Event.PhysicalResourceID, map[string]interface{}{
			"Name": name,
		}, nil
	}

	// need to create a new service.
	return s.create(ctx)
}

func (s *service) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	_, id, err := s.Function.parseID(ctx, s.Event.PhysicalResourceID, 1)
	if err != nil {
		return "", nil, err
	}

	c := s.Function.getclient()
	ss, err := c.DeleteService(ctx, id[0])
	if err != nil {
		return "", nil, err
	}

	return s.Event.PhysicalResourceID, map[string]interface{}{
		"Name": ss.Name,
	}, nil
}
