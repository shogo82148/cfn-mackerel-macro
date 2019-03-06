package cfn

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
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
		return s.create(ctx) // TODO: in place update
	case cfn.RequestDelete:
		// TODO: delete
		return s.Event.PhysicalResourceID, nil, nil
	}
	return "", nil, fmt.Errorf("unknown request type: %s", s.Event.RequestType)
}

func (s *service) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	name, ok := s.Event.ResourceProperties["Name"].(string)
	if !ok {
		return "", nil, errors.New("Name is missing")
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
		"Memo": ss.Memo,
	}, nil
}
