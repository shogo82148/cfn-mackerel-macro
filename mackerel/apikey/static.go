package apikey

import (
	"context"

	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// NewStatic returns a provider which provides a static api key.
func NewStatic(apikey string) mackerel.APIKeyProvider {
	return staticAPIKey(apikey)
}

type staticAPIKey string

func (s staticAPIKey) MackerelAPIKey(ctx context.Context) (string, error) {
	return string(s), nil
}
