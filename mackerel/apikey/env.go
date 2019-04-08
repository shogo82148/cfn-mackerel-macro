package apikey

import (
	"context"
	"fmt"
	"os"

	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// NewEnvironment returns a provider which gets a static api key from the environment values.
func NewEnvironment(name string) mackerel.APIKeyProvider {
	return environmentAPIKey(name)
}

type environmentAPIKey string

func (s environmentAPIKey) MackerelAPIKey(ctx context.Context) (string, error) {
	apikey, ok := os.LookupEnv(string(s))
	if !ok {
		return "", fmt.Errorf("environment value %s not found", string(s))
	}
	return apikey, nil
}
