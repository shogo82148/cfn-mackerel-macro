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
	_, err = c.CreateService(ctx, &mackerel.CreateServiceParam{
		Name: name,
		// TODO: memo
	})
	if err != nil {
		merr, ok := err.(mackerel.Error)
		if !ok {
			return "", nil, err
		}
		if merr.StatusCode() != http.StatusBadRequest {
			return "", nil, err
		}

		// the service may already exist. try to continue.
	}
	creationErr := err

	id, err := s.Function.buildServiceID(ctx, name)
	if err != nil {
		return "", nil, err
	}

	meta := getmetadata(s.Event)
	if err := c.PutServiceMetaData(ctx, name, "cloudformation", meta); err != nil {
		if creationErr != nil {
			return "", nil, creationErr
		}
		return id, nil, err
	}

	return id, map[string]interface{}{
		"Name": name,
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
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the service %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
