package cfn

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type service struct {
	Function *Function
	Event    cfn.Event
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

	id, err := s.Function.buildServiceID(ctx, ss.Name)
	if err != nil {
		return "", nil, err
	}

	meta := getmetadata(s.Event)
	if err := c.PutServiceMetaData(ctx, ss.Name, "cloudformation", meta); err != nil {
		return id, nil, err
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
	physicalResourceID = s.Event.PhysicalResourceID
	serviceName, err := s.Function.parseServiceID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as service id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := s.Function.getclient()
	_, err = c.DeleteService(ctx, serviceName)
	return
}
