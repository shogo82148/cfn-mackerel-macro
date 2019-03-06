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

// Handle handles custom resource events of CloudForamtion.
func (f *Function) Handle(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	typ := strings.TrimPrefix(event.ResourceType, "Custom::")
	switch typ {
	case "Service":
		s := &service{
			Function: f,
			Event:    event,
		}
		return s.handle(ctx)
	}
	return "", nil, nil // fmt.Errorf("unkdnown type: %s", typ)
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
