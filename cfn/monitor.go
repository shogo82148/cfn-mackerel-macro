package cfn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type monitor struct {
	Function *Function
	Event    cfn.Event
}

func (m *monitor) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := m.Function.getclient()
	mm, err := m.convertToParam(ctx, m.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateMonitor(ctx, mm)
	if err != nil {
		return "", nil, err
	}

	id, err := m.Function.buildMonitorID(ctx, ret.MonitorID())
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"MonitorId": ret.MonitorID(),
		"Type":      ret.MonitorType().String(),
		"Name":      ret.MonitorName(),
	}, nil
}

func (m *monitor) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	id, err := m.Function.parseMonitorID(ctx, m.Event.PhysicalResourceID)
	if err != nil {
		return m.Event.PhysicalResourceID, nil, err
	}

	c := m.Function.getclient()
	mm, err := m.convertToParam(ctx, m.Event.ResourceProperties)
	if err != nil {
		return m.Event.PhysicalResourceID, nil, err
	}
	ret, err := c.UpdateMonitor(ctx, id, mm)
	if err != nil {
		return m.Event.PhysicalResourceID, nil, err
	}

	return m.Event.PhysicalResourceID, map[string]interface{}{
		"MonitorId": ret.MonitorID(),
		"Type":      ret.MonitorType().String(),
		"Name":      ret.MonitorName(),
	}, nil
}

func (m *monitor) convertToParam(ctx context.Context, properties map[string]interface{}) (mackerel.Monitor, error) {
	in := dproxy.New(properties)
	typ, err := in.M("Type").String()
	if err != nil {
		return nil, err
	}

	var d dproxy.Drain
	var mm mackerel.Monitor
	switch typ {
	case mackerel.MonitorTypeConnectivity.String():
		var scopes, excludeScopes []string
		for _, item := range d.Array(dproxy.Default(in.M("Scopes"), []interface{}{})) {
			s := d.String(dproxy.New(item))
			if serviceName, err := m.Function.parseServiceID(ctx, s); err == nil {
				scopes = append(scopes, serviceName)
			} else if serviceName, roleName, err := m.Function.parseRoleID(ctx, s); err == nil {
				scopes = append(scopes, serviceName+":"+roleName)
			} else {
				return nil, fmt.Errorf("scopes should be a service of a role: %s", s)
			}
		}
		for _, item := range d.Array(dproxy.Default(in.M("ExcludeScopes"), []interface{}{})) {
			s := d.String(dproxy.New(item))
			if serviceName, err := m.Function.parseServiceID(ctx, s); err == nil {
				excludeScopes = append(excludeScopes, serviceName)
			} else if serviceName, roleName, err := m.Function.parseRoleID(ctx, s); err == nil {
				excludeScopes = append(excludeScopes, serviceName+":"+roleName)
			} else {
				return nil, fmt.Errorf("excludeScopes should be a service of a role: %s", s)
			}
		}
		mm = &mackerel.MonitorConnectivity{
			Name:                 d.String(in.M("Name")),
			Memo:                 d.String(dproxy.Default(in.M("Memo"), "")),
			Scopes:               scopes,
			ExcludeScopes:        excludeScopes,
			NotificationInterval: uint64(d.Int64(dproxy.Default(in.M("NotificationInterval"), 0))),
		}
	case mackerel.MonitorTypeHostMetric.String():
		var scopes, excludeScopes []string
		for _, item := range d.Array(dproxy.Default(in.M("Scopes"), []interface{}{})) {
			s := d.String(dproxy.New(item))
			if serviceName, err := m.Function.parseServiceID(ctx, s); err == nil {
				scopes = append(scopes, serviceName)
			} else if serviceName, roleName, err := m.Function.parseRoleID(ctx, s); err == nil {
				scopes = append(scopes, serviceName+":"+roleName)
			} else {
				return nil, fmt.Errorf("scopes should be a service of a role: %s", s)
			}
		}
		for _, item := range d.Array(dproxy.Default(in.M("ExcludeScopes"), []interface{}{})) {
			s := d.String(dproxy.New(item))
			if serviceName, err := m.Function.parseServiceID(ctx, s); err == nil {
				excludeScopes = append(excludeScopes, serviceName)
			} else if serviceName, roleName, err := m.Function.parseRoleID(ctx, s); err == nil {
				excludeScopes = append(excludeScopes, serviceName+":"+roleName)
			} else {
				return nil, fmt.Errorf("excludeScopes should be a service of a role: %s", s)
			}
		}
		mm = &mackerel.MonitorHostMetric{
			Name:                 d.String(in.M("Name")),
			Memo:                 d.String(dproxy.Default(in.M("Memo"), "")),
			Duration:             d.Uint64(dproxy.Default(in.M("Duration"), 1)),
			Metric:               d.String(in.M("Metric")),
			Operator:             d.String(in.M("Operator")),
			Warning:              d.OptionalFloat64(in.M("Warning")),
			Critical:             d.OptionalFloat64(in.M("Critical")),
			MaxCheckAttempts:     uint64(d.Int64(dproxy.Default(in.M("MaxCheckAttempts"), 1))),
			Scopes:               scopes,
			ExcludeScopes:        excludeScopes,
			NotificationInterval: uint64(d.Int64(dproxy.Default(in.M("NotificationInterval"), 0))),
		}
	case mackerel.MonitorTypeServiceMetric.String():
		serviceName, err := m.Function.parseServiceID(ctx, d.String(in.M("Service")))
		if err != nil {
			return nil, err
		}
		mm = &mackerel.MonitorServiceMetric{
			Name:                    d.String(in.M("Name")),
			Memo:                    d.String(dproxy.Default(in.M("Memo"), "")),
			Duration:                d.Uint64(dproxy.Default(in.M("Duration"), 1)),
			Service:                 serviceName,
			Metric:                  d.String(in.M("Metric")),
			Operator:                d.String(in.M("Operator")),
			Warning:                 d.OptionalFloat64(in.M("Warning")),
			Critical:                d.OptionalFloat64(in.M("Critical")),
			MaxCheckAttempts:        d.Uint64(dproxy.Default(in.M("MaxCheckAttempts"), 1)),
			NotificationInterval:    d.Uint64(dproxy.Default(in.M("NotificationInterval"), 0)),
			MissingDurationWarning:  d.OptionalUint64(in.M("MissingDurationWarning")),
			MissingDurationCritical: d.OptionalUint64(in.M("MissingDurationCritical")),
		}
	case mackerel.MonitorTypeExternalHTTP.String():
		var serviceName string
		if s := d.OptionalString(in.M("Service")); s != nil {
			var err error
			serviceName, err = m.Function.parseServiceID(ctx, *s)
			if err != nil {
				return nil, err
			}
		}
		var headers []mackerel.HeaderField
		h, err := in.M("Headers").ProxySet().ProxyArray()
		if err == nil {
			headers = make([]mackerel.HeaderField, 0, len(h))
			for _, item := range h {
				headers = append(headers, mackerel.HeaderField{
					Name:  d.String(item.M("Name")),
					Value: d.String(item.M("Value")),
				})
			}
		} else if perr, ok := err.(dproxy.Error); !ok || perr.ErrorType() != dproxy.Enotfound {
			return nil, err
		}
		mm = &mackerel.MonitorExternalHTTP{
			Name:        d.String(in.M("Name")),
			Memo:        d.String(dproxy.Default(in.M("Memo"), "")),
			URL:         d.String(in.M("Url")),
			Method:      d.String(dproxy.Default(in.M("Method"), "GET")),
			RequestBody: d.String(dproxy.Default(in.M("RequestBody"), "")),

			Service:              serviceName,
			NotificationInterval: d.Uint64(dproxy.Default(in.M("NotificationInterval"), 0)),
			ResponseTimeWarning:  d.OptionalFloat64(in.M("ResponseTimeWarning")),
			ResponseTimeCritical: d.OptionalFloat64(in.M("ResponseTimeCritical")),
			ResponseTimeDuration: d.OptionalUint64(dproxy.Default(in.M("ResponseTimeDuration"), 1)),
			ContainsString:       d.String(dproxy.Default(in.M("ContainsString"), "")),
			MaxCheckAttempts:     d.Uint64(dproxy.Default(in.M("MaxCheckAttempts"), 1)),

			CertificationExpirationWarning:  d.OptionalUint64(in.M("CertificationExpirationWarning")),
			CertificationExpirationCritical: d.OptionalUint64(in.M("CertificationExpirationCritical")),
			SkipCertificateVerification:     d.Bool(dproxy.Default(in.M("SkipCertificateVerification"), false)),
			Headers:                         headers,
		}
	case mackerel.MonitorTypeExpression.String():
		mm = &mackerel.MonitorExpression{
			Name:                 d.String(in.M("Name")),
			Memo:                 d.String(dproxy.Default(in.M("Memo"), "")),
			Expression:           d.String(in.M("Expression")),
			Operator:             d.String(in.M("Operator")),
			Warning:              d.OptionalFloat64(in.M("Warning")),
			Critical:             d.OptionalFloat64(in.M("Critical")),
			NotificationInterval: d.Uint64(dproxy.Default(in.M("NotificationInterval"), 0)),
		}
	default:
		return nil, fmt.Errorf("unknown monitor type: %s", typ)
	}
	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return mm, nil
}

func (m *monitor) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = m.Event.PhysicalResourceID
	id, err := m.Function.parseHostID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as monitor id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := m.Function.getclient()
	_, err = c.DeleteMonitor(ctx, id)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
