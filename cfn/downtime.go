package cfn

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type downtime struct {
	Function *Function
	Event    cfn.Event
}

func (r *downtime) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	param, err := r.convertToParam(ctx, r.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateDowntime(ctx, param)
	if err != nil {
		return "", nil, err
	}

	physicalResourceID = "mkr:test-org:downtime:" + ret.ID
	return
}

func (r *downtime) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}

func (r *downtime) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.Downtime, error) {
	var param mackerel.Downtime
	var d dproxy.Drain
	in := dproxy.New(properties)

	param.Name = d.String(in.M("Name"))
	param.Memo = d.String(dproxy.Default(in.M("Memo"), ""))
	param.Start = mackerel.Timestamp(d.Int64(in.M("Start")))
	param.Duration = d.Int64(in.M("Duration"))

	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return &param, nil
}

func (r *downtime) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}
