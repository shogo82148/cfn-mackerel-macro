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
