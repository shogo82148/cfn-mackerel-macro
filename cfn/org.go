package cfn

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
)

type org struct {
	Function *Function
	Event    cfn.Event
}

func (o *org) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := o.Function.getclient()
	ret, err := c.GetOrg(ctx)
	if err != nil {
		return
	}

	physicalResourceID = "mkr:" + ret.Name
	data = map[string]interface{}{
		"Name": ret.Name,
	}
	return
}

func (o *org) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return o.create(ctx)
}

func (o *org) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = o.Event.PhysicalResourceID
	return
}
