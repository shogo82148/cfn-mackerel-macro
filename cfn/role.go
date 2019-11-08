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

type role struct {
	Function *Function
	Event    cfn.Event
}

func (r *role) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(r.Event.ResourceProperties)
	name := d.String(in.M("Name"))
	service := d.String(in.M("Service"))
	err = d.CombineErrors()
	if err != nil {
		return
	}

	serviceName, err := r.Function.parseServiceID(ctx, service)
	if err != nil {
		err = fmt.Errorf("failed to parse %q as service id: %w", service, err)
		return
	}

	c := r.Function.getclient()
	_, err = c.CreateRole(ctx, serviceName, &mackerel.CreateRoleParam{
		Name: name,
	})
	var creationErr error
	if err != nil {
		creationErr = fmt.Errorf("failed to create role: %w", err)
		var merr mackerel.Error
		if errors.As(err, &merr) {
			if merr.StatusCode() != http.StatusBadRequest {
				err = creationErr
				return
			}
		}

		log.Println("the role may already exist. try to override it.")
		err = nil
	}

	physicalResourceID, err = r.Function.buildRoleID(ctx, serviceName, name)
	if err != nil {
		return
	}
	meta := getmetadata(r.Event)
	if err := c.PutRoleMetaData(ctx, serviceName, name, "cloudformation", meta); err != nil {
		if creationErr != nil {
			return physicalResourceID, nil, creationErr
		}
		return physicalResourceID, nil, fmt.Errorf("failed to put role metadata: %w", err)
	}

	data = map[string]interface{}{
		"Name":     name,
		"FullName": serviceName + ":" + name,
	}
	return
}

func (r *role) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(r.Event.ResourceProperties)
	old := dproxy.New(r.Event.OldResourceProperties)

	name := d.String(in.M("Name"))
	service := d.String(in.M("Service"))
	oldName := d.String(old.M("Name"))
	oldService := d.String(old.M("Service"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	serviceName, err := r.Function.parseServiceID(ctx, service)
	if err != nil {
		err = fmt.Errorf("failed to parse %q as service id: %s", service, err)
		return
	}
	if name == oldName && service == oldService {
		// No update is needed.
		return r.Event.PhysicalResourceID, map[string]interface{}{
			"Name":     name,
			"FullName": serviceName + ":" + name,
		}, nil
	}

	// need to create a new role.
	return r.create(ctx)
}

func (r *role) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = r.Event.PhysicalResourceID
	serviceName, roleName, err := r.Function.parseRoleID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as role id: %s", physicalResourceID, err)
		err = nil
		return
	}

	c := r.Function.getclient()
	_, err = c.DeleteRole(ctx, serviceName, roleName)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the role %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil
	}
	return
}
