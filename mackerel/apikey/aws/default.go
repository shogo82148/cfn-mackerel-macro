package aws

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel/apikey"
)

// LoadDefaultProvider returns default provider.
func LoadDefaultProvider() (mackerel.APIKeyProvider, error) {
	cfg, err := external.LoadDefaultAWSConfig()
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
