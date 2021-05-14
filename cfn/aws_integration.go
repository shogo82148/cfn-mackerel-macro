package cfn

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type awsIntegration struct {
	Function *Function
	Event    cfn.Event
}

func (r *awsIntegration) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	param, err := r.convertToParam(ctx, r.Event.ResourceProperties)
	if err != nil {
		return "", nil, err
	}
	ret, err := c.CreateAWSIntegration(ctx, param)
	if err != nil {
		return "", nil, err
	}
	id, err := r.Function.buildAWSIntegrationID(ctx, ret.ID)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{}, nil
}

func (r *awsIntegration) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	param, err := r.convertToParam(ctx, r.Event.ResourceProperties)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}
	id, err := r.Function.parseAWSIntegrationID(ctx, r.Event.PhysicalResourceID)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}
	_, err = c.UpdateAWSIntegration(ctx, id, param)
	if err != nil {
		return r.Event.PhysicalResourceID, nil, err
	}
	return r.Event.PhysicalResourceID, map[string]interface{}{}, nil
}

func (r *awsIntegration) convertToParam(ctx context.Context, properties map[string]interface{}) (*mackerel.AWSIntegration, error) {
	var d dproxy.Drain
	in := dproxy.New(properties)

	param := &mackerel.AWSIntegration{
		Name:         d.String(in.M("Name")),
		Memo:         d.String(dproxy.Default(in.M("Memo"), "")),
		Key:          d.OptionalString(in.M("Key")),
		SecretKey:    d.OptionalString(in.M("SecretKey")),
		RoleArn:      d.OptionalString(in.M("RoleArn")),
		ExternalID:   d.OptionalString(in.M("ExternalID")),
		Region:       d.String(in.M("Region")),
		IncludedTags: r.convertTagList(&d, dproxy.Default(in.M("IncludedTags"), []interface{}{})),
		ExcludedTags: r.convertTagList(&d, dproxy.Default(in.M("ExcludedTags"), []interface{}{})),
		Services:     r.convertAWSServices(ctx, &d, in.M("Services")),
	}

	if err := d.CombineErrors(); err != nil {
		return nil, err
	}
	return param, nil
}

// convert the list of tags.
// https://mackerel.io/docs/entry/integrations/aws#tag
func (r *awsIntegration) convertTagList(d *dproxy.Drain, properties dproxy.Proxy) string {
	var tags []string
	for _, tag := range d.ProxyArray(properties.ProxySet()) {
		name := d.String(tag.M("Name"))
		value := d.String(tag.M("Value"))
		tags = append(tags, r.escapeTagValue(name)+":"+r.escapeTagValue(value))
	}
	return strings.Join(tags, ",")
}

// ':' and ',' have special meaning, so we need to escape them.
// https://mackerel.io/docs/entry/integrations/aws#tag
func (*awsIntegration) escapeTagValue(s string) string {
	if !strings.ContainsAny(s, ":, '\"") {
		// fast pass, no need to escape
		return s
	}
	hasDoubleQuote := strings.ContainsRune(s, '"')
	hasSingleQuote := strings.ContainsRune(s, '\'')
	if !hasDoubleQuote {
		return "\"" + s + "\""
	}
	if !hasSingleQuote {
		return "'" + s + "'"
	}
	// XXX: the string contains '"' and '\''. How should we handle it?
	return s
}

func (r *awsIntegration) convertAWSServices(ctx context.Context, d *dproxy.Drain, properties dproxy.Proxy) map[string]*mackerel.AWSIntegrationService {
	ret := map[string]*mackerel.AWSIntegrationService{}
	for _, s := range d.ProxyArray(properties.ProxySet()) {
		name := d.String(s.M("ServiceId"))
		exclude := dproxy.Default(s.M("ExcludedMetrics"), []interface{}{})

		var role *string
		roleID := d.OptionalString(s.M("Role"))
		if roleID != nil {
			service, name, err := r.Function.parseRoleID(ctx, *roleID)
			d.Put(err)
			fullname := service + ":" + name
			role = &fullname
		}

		ret[name] = &mackerel.AWSIntegrationService{
			Enable:              d.Bool(s.M("Enable")),
			Role:                role,
			ExcludedMetrics:     d.StringArray(exclude.ProxySet()),
			RetireAutomatically: d.Bool(dproxy.Default(s.M("RetireAutomatically"), false)),
		}
	}
	return ret
}

func (r *awsIntegration) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	physicalResourceID = r.Event.PhysicalResourceID
	id, err := r.Function.parseDowntimeID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as aws integration id: %s", physicalResourceID, err)
		err = nil // ignore it
		return
	}
	_, err = c.DeleteAWSIntegration(ctx, id)
	var merr mackerel.Error
	if errors.As(err, &merr) && merr.StatusCode() == http.StatusNotFound {
		log.Printf("It seems that the aws integration %q is already deleted, ignore the error: %s", physicalResourceID, err)
		err = nil // ignore it
	}
	return
}

type awsIntegrationExternalID struct {
	Function *Function
	Event    cfn.Event
}

func (r *awsIntegrationExternalID) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	ret, err := c.CreateAWSIntegrationExternalID(ctx)
	if err != nil {
		return "", nil, err
	}
	id, err := r.Function.buildAWSIntegrationExternalID(ctx, ret)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Id": ret,
	}, nil
}

func (r *awsIntegrationExternalID) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	c := r.Function.getclient()
	ret, err := c.CreateAWSIntegrationExternalID(ctx)
	if err != nil {
		return "", nil, err
	}
	id, err := r.Function.buildAWSIntegrationExternalID(ctx, ret)
	if err != nil {
		return "", nil, err
	}
	return id, map[string]interface{}{
		"Id": ret,
	}, nil
}

func (r *awsIntegrationExternalID) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = r.Event.PhysicalResourceID
	_, err = r.Function.parseAWSIntegrationExternalID(ctx, physicalResourceID)
	if err != nil {
		log.Printf("failed to parse %q as aws integration external id: %s", physicalResourceID, err)
		err = nil // ignore it
		return
	}

	// Mackerel doesn't provide to delete external ids.
	// nothing to do.
	return
}
