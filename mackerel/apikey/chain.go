package apikey

import (
	"context"
	"fmt"
	"strings"

	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// NewChain searches a valid provider from providers.
func NewChain(providers ...mackerel.APIKeyProvider) mackerel.APIKeyProvider {
	return chainAPIKey(providers)
}

type chainAPIKey []mackerel.APIKeyProvider

func (s chainAPIKey) MackerelAPIKey(ctx context.Context) (string, error) {
	var errs []string
	for _, p := range s {
		apikey, err := p.MackerelAPIKey(ctx)
		if err == nil {
			return apikey, nil
		}
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		errs = append(errs, err.Error())
	}

	return "", fmt.Errorf("no valid providers: %s", strings.Join(errs, ", "))
}
