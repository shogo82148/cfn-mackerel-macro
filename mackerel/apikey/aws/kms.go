package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// KMS decrypts the api key with AWS Key Management Service.
type KMS struct {
	mackerel.APIKeyProvider
	Config aws.Config
}

// MackerelAPIKey implements mackerel.APIKeyProvider
func (p *KMS) MackerelAPIKey(ctx context.Context) (string, error) {
	apikey, err := p.APIKeyProvider.MackerelAPIKey(ctx)
	if err != nil {
		return "", err
	}
	b, err := base64.StdEncoding.DecodeString(apikey)
	if err != nil {
		return "", err
	}
	svc := kms.NewFromConfig(p.Config)
	resp, err := svc.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: b,
	})
	if err != nil {
		return "", err
	}
	return string(resp.Plaintext), nil
}
