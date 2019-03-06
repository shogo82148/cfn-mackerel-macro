package mackerel

import (
	"context"
	"net/http"
)

// Service represents the "service" of Mackerel.
type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

// CreateServiceParam parameters for CreateService.
type CreateServiceParam struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CreateService creates service
func (c *Client) CreateService(ctx context.Context, param *CreateServiceParam) (*Service, error) {
	service := &Service{}
	err := c.do(ctx, http.MethodPost, "/api/v0/services", param, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}
