package mackerel

import (
	"context"
	"fmt"
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

// FindServices return the list of services.
func (c *Client) FindServices(ctx context.Context) ([]*Service, error) {
	var resp struct {
		Services []*Service `json:"services"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/services", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Services, nil
}

// CreateService creates a new service
func (c *Client) CreateService(ctx context.Context, param *CreateServiceParam) (*Service, error) {
	service := &Service{}
	_, err := c.do(ctx, http.MethodPost, "/api/v0/services", param, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// DeleteService deletes a service
func (c *Client) DeleteService(ctx context.Context, serviceName string) (*Service, error) {
	service := &Service{}
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/services/%s", serviceName), nil, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}
