package mackerel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ServiceMetaMetaData is meta data of service meta data.
type ServiceMetaMetaData struct {
	LastModified time.Time
}

// GetServiceMetaData gets Service metadata and stores the result in the value pointed to by v.
// GetServiceMetaData uses the json package for storing the result, see https://golang.org/pkg/encoding/json/#Unmarshal for decoding rules.
// https://mackerel.io/api-docs/entry/metadata#get
func (c *Client) GetServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) (*ServiceMetaMetaData, error) {
	h, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace), nil, v)
	if err != nil {
		return nil, err
	}

	ret := &ServiceMetaMetaData{}
	ret.LastModified, err = http.ParseTime(h.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetServiceMetaDataNameSpaces fetches namespaces of Service metadata.
// https://mackerel.io/api-docs/entry/metadata#Servicelist
func (c *Client) GetServiceMetaDataNameSpaces(ctx context.Context, serviceName string) ([]string, error) {
	var data struct {
		Metadata []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}
	_, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/services/%s/metadata", serviceName), nil, &data)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(data.Metadata))
	for _, metadata := range data.Metadata {
		ret = append(ret, metadata.NameSpace)
	}
	return ret, nil
}

// PutServiceMetaData creates or updates Service metadata by the value of v.
// PutServiceMetaData uses the json package for putting the metadata, see https://golang.org/pkg/encoding/json/#Marshal for encoding roles.
// https://mackerel.io/api-docs/entry/metadata#serviceput
func (c *Client) PutServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace), v, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}

// DeleteServiceMetaData deletes Service metadata by the value of v.
// https://mackerel.io/api-docs/entry/metadata#Servicedelete
func (c *Client) DeleteServiceMetaData(ctx context.Context, serviceName, namespace string) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace), nil, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}
