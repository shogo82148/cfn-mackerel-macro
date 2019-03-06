// The mackerel package is a API client library for mackerel.io

package mackerel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var defaultBaseURL *url.URL

func init() {
	var err error
	defaultBaseURL, err = url.Parse("https://api.mackerelio.com/")
	if err != nil {
		panic(err)
	}
}

// Client is a client for mackerel.io
type Client struct {
	BaseURL    *url.URL
	APIKey     string
	UserAgent  string
	HTTPClient *http.Client
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Client) urlfor(path string) string {
	base := c.BaseURL
	if base == nil {
		base = defaultBaseURL
	}

	// shallow copy
	u := new(url.URL)
	*u = *base

	u.Path = path
	return u.String()
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.urlfor(path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("X-Api-Key", c.APIKey)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else {
		agent := fmt.Sprintf("cfn-mackerel-macro/0.0.0")
		req.Header.Set("User-Agent", agent)
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return handleError(resp)
	}

	if out == nil {
		// ignore the body
		io.Copy(ioutil.Discard, resp.Body)
	} else {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(out); err != nil {
			return err
		}
	}

	return nil
}

// Error is an error from the Mackerel.
type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("status: %d, %s", e.StatusCode, e.Message)
}

func handleError(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return Error{
		StatusCode: resp.StatusCode,
		Message:    string(b),
	}
}
