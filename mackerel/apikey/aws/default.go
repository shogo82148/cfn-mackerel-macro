package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel/apikey"
)

// LoadDefaultProvider returns default provider.
func LoadDefaultProvider(ctx context.Context) (mackerel.APIKeyProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	var decrypt bool
	if os.Getenv("MACKEREL_APIKEY_WITH_DECRYPT") != "" {
		decrypt = true
	}

	provider := apikey.NewEnvironment("MACKEREL_APIKEY")
	if decrypt {
		provider = &KMS{
			APIKeyProvider: provider,
			Config:         cfg,
		}
	}

	provider = apikey.NewChain(
		provider,
		&SSM{
			APIKeyProvider: apikey.NewEnvironment("MACKEREL_APIKEY_PARAMETER"),
			Config:         cfg,
			WithDecryption: decrypt,
		},
	)
	return provider, nil
}
