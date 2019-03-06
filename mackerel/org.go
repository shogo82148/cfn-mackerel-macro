package mackerel

import (
	"context"
	"net/http"
)

// Org information
type Org struct {
	Name string `json:"name"`
}

// GetOrg get the org
func (c *Client) GetOrg(ctx context.Context) (*Org, error) {
	org := &Org{}
	err := c.do(ctx, http.MethodGet, "/api/v0/org", nil, org)
	if err != nil {
		return nil, err
	}
	return org, nil
}
