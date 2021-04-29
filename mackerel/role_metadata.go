package mackerel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// RoleMetaMetaData is meta data of Role meta data.
type RoleMetaMetaData struct {
	LastModified time.Time
}

// GetRoleMetaData gets role metadata and stores the result in the value pointed to by v.
// GetRoleMetaData uses the json package for storing the result, see https://golang.org/pkg/encoding/json/#Unmarshal for decoding rules.
// https://mackerel.io/api-docs/entry/metadata#roleget
func (c *Client) GetRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) (*RoleMetaMetaData, error) {
	h, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace), nil, v)
	if err != nil {
		return nil, err
	}

	ret := &RoleMetaMetaData{}
	ret.LastModified, err = http.ParseTime(h.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetRoleMetaDataNameSpaces fetches namespaces of role metadata.
// https://mackerel.io/api-docs/entry/metadata#rolelist
func (c *Client) GetRoleMetaDataNameSpaces(ctx context.Context, serviceName, roleName string) ([]string, error) {
	var data struct {
		Metadata []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}
	_, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata", serviceName, roleName), nil, &data)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(data.Metadata))
	for _, metadata := range data.Metadata {
		ret = append(ret, metadata.NameSpace)
	}
	return ret, nil
}

// PutRoleMetaData creates or updates Role metadata by the value of v.
// PutRoleMetaData uses the json package for putting the metadata, see https://golang.org/pkg/encoding/json/#Marshal for encoding roles.
// https://mackerel.io/api-docs/entry/metadata#roleput
func (c *Client) PutRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace), v, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}

// DeleteRoleMetaData deletes Role metadata by the value of v.
// https://mackerel.io/api-docs/entry/metadata#roledelete
func (c *Client) DeleteRoleMetaData(ctx context.Context, serviceName, roleName, namespace string) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace), nil, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}
