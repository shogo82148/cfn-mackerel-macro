package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// SSM gets the api key from AWS Systems Manager Parameter Store.
type SSM struct {
	mackerel.APIKeyProvider
	Config         aws.Config
	WithDecryption bool
}

// MackerelAPIKey implements mackerel.APIKeyProvider
func (p *SSM) MackerelAPIKey(ctx context.Context) (string, error) {
	apikey, err := p.APIKeyProvider.MackerelAPIKey(ctx)
	if err != nil {
		return "", err
	}
	svc := ssm.New(p.Config)
	req := svc.GetParameterRequest(&ssm.GetParameterInput{
		Name:           aws.String(apikey),
		WithDecryption: aws.Bool(p.WithDecryption),
	})
	resp, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	return aws.StringValue(resp.Parameter.Value), nil
}
