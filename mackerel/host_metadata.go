package mackerel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// HostMetaMetaData is meta data of host meta data.
type HostMetaMetaData struct {
	LastModified time.Time
}

// GetHostMetaData gets host metadata and stores the result in the value pointed to by v.
// GetHostMetaData uses the json package for storing the result, see https://golang.org/pkg/encoding/json/#Unmarshal for decoding rules.
// https://mackerel.io/api-docs/entry/metadata#get
func (c *Client) GetHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) (*HostMetaMetaData, error) {
	h, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace), nil, v)
	if err != nil {
		return nil, err
	}

	ret := &HostMetaMetaData{}
	ret.LastModified, err = http.ParseTime(h.Get("Last-Modified"))
	return ret, nil
}

// GetHostMetaDataNameSpaces fetches namespaces of host metadata.
// https://mackerel.io/api-docs/entry/metadata#hostlist
func (c *Client) GetHostMetaDataNameSpaces(ctx context.Context, hostID string) ([]string, error) {
	var data struct {
		Metadata []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}
	_, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v0/hosts/%s/metadata", hostID), nil, &data)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(data.Metadata))
	for _, metadata := range data.Metadata {
		ret = append(ret, metadata.NameSpace)
	}
	return ret, nil
}

// PutHostMetaData creates or updates host metadata by the value of v.
// PutHostMetaData uses the json package for putting the metadata, see https://golang.org/pkg/encoding/json/#Marshal for encoding roles.
// https://mackerel.io/api-docs/entry/metadata#serviceput
func (c *Client) PutHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace), v, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}

// DeleteHostMetaData deletes host metadata by the value of v.
// https://mackerel.io/api-docs/entry/metadata#hostdelete
func (c *Client) DeleteHostMetaData(ctx context.Context, hostID, namespace string) error {
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace), nil, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("mackerel: unexpected response")
	}
	return nil
}
