package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// AWSIntegration is an AWS integration.
// https://mackerel.io/api-docs/entry/aws-integration
type AWSIntegration struct {
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	Key          *string `json:"key,omitempty"`
	SecretKey    *string `json:"secretKey,omitempty"`
	RoleArn      *string `json:"roleArn,omitempty"`
	ExternalID   *string `json:"externalId,omitempty"`
	Region       string  `json:"region,omitempty"`
	IncludedTags string  `json:"includedTags,omitempty"`
	ExcludedTags string  `json:"excludedTags,omitempty"`

	Services map[string]*AWSIntegrationService `json:"services,omitempty"`
}

// AWSIntegrationService is an AWS service.
type AWSIntegrationService struct {
	Enable              bool     `json:"enable,omitempty"`
	Role                string   `json:"role,omitempty"`
	ExcludedMetrics     []string `json:"excludedMetrics,omitempty"`
	RetireAutomatically bool     `json:"retireAutomatically,omitempty"`
}

// FindAWSIntegrations returns a list of aws integrations.
func (c *Client) FindAWSIntegrations(ctx context.Context) ([]*AWSIntegration, error) {
	var integrations struct {
		AWSIntegrations []*AWSIntegration `json:"aws_integrations"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/aws-integrations", nil, &integrations)
	if err != nil {
		return nil, err
	}
	return integrations.AWSIntegrations, nil
}

// FindAWSIntegration returns an aws integration.
func (c *Client) FindAWSIntegration(ctx context.Context, awsIntegrationID string) (*AWSIntegration, error) {
	var awsIntegration AWSIntegration
	_, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID), nil, &awsIntegration)
	if err != nil {
		return nil, err
	}
	return &awsIntegration, nil
}

// CreateAWSIntegration creates a new aws integration.
func (c *Client) CreateAWSIntegration(ctx context.Context, param *AWSIntegration) (*AWSIntegration, error) {
	var awsIntegration AWSIntegration
	_, err := c.do(ctx, http.MethodPost, "/api/v0/aws-integrations", param, &awsIntegration)
	if err != nil {
		return nil, err
	}
	return &awsIntegration, nil
}

// UpdateAWSIntegration updates an aws integration.
func (c *Client) UpdateAWSIntegration(ctx context.Context, awsIntegrationID string, param *AWSIntegration) (*AWSIntegration, error) {
	var awsIntegration AWSIntegration
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID), param, &awsIntegration)
	if err != nil {
		return nil, err
	}
	return &awsIntegration, nil
}

// DeleteAWSIntegration deletes an aws integration.
func (c *Client) DeleteAWSIntegration(ctx context.Context, awsIntegrationID string) (*AWSIntegration, error) {
	var awsIntegration AWSIntegration
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID), nil, &awsIntegration)
	if err != nil {
		return nil, err
	}
	return &awsIntegration, nil
}

// CreateAWSIntegrationExternalID creates an external id for aws integrations.
func (c *Client) CreateAWSIntegrationExternalID(ctx context.Context) (string, error) {
	var resp struct {
		ExternalID string `json:"externalId"`
	}
	_, err := c.do(ctx, http.MethodPost, "/api/v0/aws-integrations-external-id", nil, &resp)
	if err != nil {
		return "", err
	}
	return resp.ExternalID, nil
}

// FindAWSIntegrationsExcludableMetrics list excludable metrics for AWS Integration.
func (c *Client) FindAWSIntegrationsExcludableMetrics(ctx context.Context) (map[string][]string, error) {
	var resp map[string][]string
	_, err := c.do(ctx, http.MethodGet, "/api/v0/aws-integrations-excludable-metrics", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
