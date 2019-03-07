package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// Host is host information
type Host struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	DisplayName      string `json:"displayName,omitempty"`
	CustomIdentifier string `json:"customIdentifier,omitempty"`
	Type             string `json:"type"`
	Status           string `json:"status"`
	Memo             string `json:"memo"`
	// Roles            Roles       `json:"roles"`
	IsRetired bool     `json:"isRetired"`
	CreatedAt int32    `json:"createdAt"`
	Meta      HostMeta `json:"meta"`
	// Interfaces       []Interface `json:"interfaces"`
}

// HostMeta host meta informations
type HostMeta struct {
	AgentRevision string `json:"agent-revision,omitempty"`
	AgentVersion  string `json:"agent-version,omitempty"`
	AgentName     string `json:"agent-name,omitempty"`
	// BlockDevice   BlockDevice `json:"block_device,omitempty"`
	// CPU           CPU         `json:"cpu,omitempty"`
	// Filesystem    FileSystem  `json:"filesystem,omitempty"`
	// Kernel        Kernel      `json:"kernel,omitempty"`
	// Memory        Memory      `json:"memory,omitempty"`
	// Cloud         *Cloud      `json:"cloud,omitempty"`
}

// CreateHostParam parameters for CreateHost
type CreateHostParam struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName,omitempty"`
	Meta        HostMeta `json:"meta"`
	// Interfaces       []Interface   `json:"interfaces,omitempty"`
	RoleFullnames []string `json:"roleFullnames,omitempty"`
	// Checks           []CheckConfig `json:"checks,omitempty"`
	CustomIdentifier string `json:"customIdentifier,omitempty"`
}

// CreateHost creates new host
func (c *Client) CreateHost(ctx context.Context, param *CreateHostParam) (string, error) {
	var data struct {
		ID string `json:"id"`
	}
	err := c.do(ctx, http.MethodPost, "/api/v0/hosts", param, &data)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

// RetireHost make the host retuired.
func (c *Client) RetireHost(ctx context.Context, id string) error {
	param := map[string]string{}
	return c.do(ctx, http.MethodPost, fmt.Sprintf("/api/v0/hosts/%s/retire", id), param, nil)
}
