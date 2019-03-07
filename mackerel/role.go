package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// Role represents the "role" of Mackerel.
type Role struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CreateRoleParam parameters for CreateRole.
type CreateRoleParam struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CreateRole creates a new role
func (c *Client) CreateRole(ctx context.Context, serviceName string, param *CreateRoleParam) (*Role, error) {
	role := &Role{}
	err := c.do(ctx, http.MethodPost, fmt.Sprintf("/api/v0/services/%s/roles", serviceName), param, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRole deletes a role
func (c *Client) DeleteRole(ctx context.Context, serviceName, roleName string) (*Role, error) {
	role := &Role{}
	err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/services/%s/roles/%s", serviceName, roleName), nil, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
