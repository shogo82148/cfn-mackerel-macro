// The mackerel package is a API client library for mackerel.io

package mackerel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

// APIKeyProvider is an api key provider.
type APIKeyProvider interface {
	MackerelAPIKey(context.Context) (string, error)
}

// APIKeyProviderFunc type is an adapter to allow the use of ordinary functions.
type APIKeyProviderFunc func(context.Context) (string, error)

// MackerelAPIKey implements APIKeyProvider
func (f APIKeyProviderFunc) MackerelAPIKey(ctx context.Context) (string, error) {
	return f(ctx)
}

// Client is a client for mackerel.io
type Client struct {
	BaseURL        *url.URL
	APIKey         string
	APIKeyProvider APIKeyProvider
	UserAgent      string
	HTTPClient     *http.Client

	mu     sync.RWMutex
	apikey string // cached api key
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

func (c *Client) getAPIKey(ctx context.Context) (string, error) {
	// check static api key
	if c.APIKey != "" {
		return c.APIKey, nil
	}

	// check cached api key
	c.mu.RLock()
	if c.apikey != "" {
		key := c.apikey
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	// need to update api key
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.apikey != "" {
		return c.apikey, nil
	}

	provider := c.APIKeyProvider
	if provider == nil {
		return "", errors.New("api key is not found")
	}
	apikey, err := provider.MackerelAPIKey(ctx)
	if err != nil {
		return "", err
	}
	c.apikey = apikey
	return apikey, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.urlfor(path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	apikey, err := c.getAPIKey(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", apikey)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else {
		req.Header.Set("User-Agent", "cfn-mackerel-macro/main")
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) (http.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}

	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Header, handleError(resp)
	}

	if out == nil {
		// ignore the body
		io.Copy(io.Discard, resp.Body)
	} else {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(out); err != nil {
			return nil, err
		}
	}

	return resp.Header, nil
}

// Error is an error from the Mackerel.
type Error interface {
	StatusCode() int
	Message() string
}

type mkrError struct {
	statusCode int
	message    string
}

func (e mkrError) Error() string {
	return fmt.Sprintf("status: %d, %s", e.statusCode, e.message)
}

func (e mkrError) StatusCode() int {
	return e.statusCode
}

func (e mkrError) Message() string {
	return e.message
}

func handleError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct{ Error struct{ Message string } }
	err = json.Unmarshal(b, &data)
	if err != nil {
		return mkrError{
			statusCode: resp.StatusCode,
			message:    string(b),
		}
	}
	return mkrError{
		statusCode: resp.StatusCode,
		message:    data.Error.Message,
	}
}

// Timestamp is unix epoch time.
type Timestamp int64

// MarshalJSON implements the json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 20)
	buf = strconv.AppendInt(buf, int64(t), 10)
	return buf, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t Timestamp) MarshalText() ([]byte, error) {
	buf := make([]byte, 0, 20)
	buf = strconv.AppendInt(buf, int64(t), 10)
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(unix)
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Timestamp) UnmarshalText(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(unix)
	return nil
}

// Time converts t to time.Time type.
func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}
