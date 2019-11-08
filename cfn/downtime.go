package cfn

import (
	"context"
	"fmt"

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

	recurrence := in.M("Recurrence")
	if _, err := recurrence.Map(); err == nil {
		var typ mackerel.DowntimeRecurrenceType
		switch d.String(recurrence.M("Type")) {
		case "hourly":
			typ = mackerel.DowntimeRecurrenceTypeHourly
		case "daily":
			typ = mackerel.DowntimeRecurrenceTypeDaily
		case "weekly":
			typ = mackerel.DowntimeRecurrenceTypeWeekly
		case "monthly":
			typ = mackerel.DowntimeRecurrenceTypeMonthly
		case "yearly":
			typ = mackerel.DowntimeRecurrenceTypeYearly
		default:
			d.Put(fmt.Errorf("unknown recurrence type: %s", d.String(recurrence.M("Type"))))
		}

		var weekdays []mackerel.DowntimeWeekday
		if days, err := recurrence.M("Weekdays").ProxySet().StringArray(); err == nil {
			if typ != mackerel.DowntimeRecurrenceTypeWeekly {
				d.Put(fmt.Errorf("weekdays are available with weekly type, but it is %s type", typ))
			}
			weekdays = make([]mackerel.DowntimeWeekday, 0, len(days))
			for _, day := range days {
				if ret, err := mackerel.ParseDowntimeWeekday(day); err != nil {
					d.Put(err)
				} else {
					weekdays = append(weekdays, ret)
				}
			}
		} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
			d.Put(err)
		}

		param.Recurrence = &mackerel.DowntimeRecurrence{
			Type:     typ,
			Interval: d.Int64(recurrence.M("Interval")),
			Weekdays: weekdays,
			Until:    (*mackerel.Timestamp)(d.OptionalInt64(recurrence.M("Until"))),
		}
	} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
		d.Put(err)
	}

	// Service Scopes
	if scopes, err := in.M("ServiceScopes").ProxySet().StringArray(); err == nil {
		services := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			name, err := r.Function.parseServiceID(ctx, scope)
			d.Put(err)
			services = append(services, name)
		}
		param.ServiceScopes = services
	} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
		d.Put(err)
	}

	// Service Exclude Scopes
	if scopes, err := in.M("ServiceExcludeScopes").ProxySet().StringArray(); err == nil {
		services := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			name, err := r.Function.parseServiceID(ctx, scope)
			d.Put(err)
			services = append(services, name)
		}
		param.ServiceExcludeScopes = services
	} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
		d.Put(err)
	}

	// Role Scopes
	if scopes, err := in.M("RoleScopes").ProxySet().StringArray(); err == nil {
		roles := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			role, service, err := r.Function.parseRoleID(ctx, scope)
			d.Put(err)
			roles = append(roles, role+":"+service)
		}
		param.RoleScopes = roles
	} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
		d.Put(err)
	}

	// Role Exclude Scopes
	if scopes, err := in.M("RoleExcludeScopes").ProxySet().StringArray(); err == nil {
		roles := make([]string, 0, len(scopes))
		for _, scope := range scopes {
			role, service, err := r.Function.parseRoleID(ctx, scope)
			d.Put(err)
			roles = append(roles, role+":"+service)
		}
		param.RoleExcludeScopes = roles
	} else if !dproxy.IsErrorCode(err, dproxy.ErrorCodeNotFound) {
		d.Put(err)
	}

	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return &param, nil
}

func (r *downtime) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	return
}
