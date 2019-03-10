package cfn

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// Function is a custom resource function for CloudForamtion.
type Function struct {
	APIKey string

	mu     sync.Mutex
	client *mackerel.Client
	org    *mackerel.Org
}

type resource interface {
	create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
	update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
	delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
}

// Handle handles custom resource events of CloudForamtion.
func (f *Function) Handle(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	typ := strings.TrimPrefix(event.ResourceType, "Custom::")
	var r resource
	switch typ {
	case "Service":
		r = &service{
			Function: f,
			Event:    event,
		}
	case "Role":
		r = &role{
			Function: f,
			Event:    event,
		}
	case "Host":
		r = &host{
			Function: f,
			Event:    event,
		}
	case "Monitor":
		r = &monitor{
			Function: f,
			Event:    event,
		}
	default:
		return "", nil, nil // fmt.Errorf("unkdnown type: %s", typ)
	}
	switch event.RequestType {
	case cfn.RequestCreate:
		return r.create(ctx)
	case cfn.RequestUpdate:
		return r.update(ctx)
	case cfn.RequestDelete:
		return r.delete(ctx)
	}
	return "", nil, fmt.Errorf("unknown request type: %s", event.RequestType)
}

// LambdaWrap returns a CustomResourceLambdaFunction which is something lambda.Start()
// will understand.
func (f *Function) LambdaWrap() cfn.CustomResourceLambdaFunction {
	return cfn.LambdaWrap(f.Handle)
}

func (f *Function) getclient() *mackerel.Client {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.client == nil {
		f.client = &mackerel.Client{
			APIKey: f.APIKey,
		}
	}
	return f.client
}

func (f *Function) getorg(ctx context.Context) (*mackerel.Org, error) {
	c := f.getclient()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.org == nil {
		org, err := c.GetOrg(ctx)
		if err != nil {
			return nil, err
		}
		f.org = org
	}
	return f.org, nil
}

func (f *Function) buildID(ctx context.Context, resourceType string, ids ...string) (string, error) {
	org, err := f.getorg(ctx)
	if err != nil {
		return "", err
	}

	ret := append([]string{"mkr", org.Name, resourceType}, ids...)
	return strings.Join(ret, ":"), nil
}

func (f *Function) buildServiceID(ctx context.Context, serviceName string) (string, error) {
	return f.buildID(ctx, "service", serviceName)
}

func (f *Function) buildRoleID(ctx context.Context, serviceName, roleName string) (string, error) {
	return f.buildID(ctx, "role", serviceName, roleName)
}

func (f *Function) buildHostID(ctx context.Context, hostID string) (string, error) {
	return f.buildID(ctx, "host", hostID)
}

func (f *Function) buildMonitorID(ctx context.Context, monitorID string) (string, error) {
	return f.buildID(ctx, "monitor", monitorID)
}

// parseID parses ID of Mackerel resources.
func (f *Function) parseID(ctx context.Context, id string, n int) (string, []string, error) {
	org, err := f.getorg(ctx)
	if err != nil {
		return "", nil, err
	}

	ids := strings.Split(id, ":")
	if len(ids) < n+3 {
		return "", nil, fmt.Errorf("invalid mkr id: %s", id)
	}
	if ids[0] != "mkr" {
		return "", nil, fmt.Errorf("invalid mkr id: %s", id)
	}
	if ids[1] != org.Name {
		return "", nil, fmt.Errorf("invalid org name in id: %s", id)
	}
	return ids[2], ids[3:], nil
}

func (f *Function) parseServiceID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "service" {
		return "", fmt.Errorf("invalid type %s, expected service type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseRoleID(ctx context.Context, id string) (string, string, error) {
	typ, parts, err := f.parseID(ctx, id, 2)
	if err != nil {
		return "", "", err
	}
	if typ != "role" {
		return "", "", fmt.Errorf("invalid type %s, expected role type", typ)
	}
	return parts[0], parts[1], nil
}

func (f *Function) parseHostID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "host" {
		return "", fmt.Errorf("invalid type %s, expected host type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseMonitorID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "monitor" {
		return "", fmt.Errorf("invalid type %s, expected monitor type", typ)
	}
	return parts[0], nil
}
