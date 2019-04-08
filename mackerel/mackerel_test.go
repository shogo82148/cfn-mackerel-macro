package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	const apiKey = "DUMMY-APi-KEY"
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != apiKey {
			t.Errorf("unexpected api key, want %s, got %s", apiKey, r.Header.Get("X-Api-Key"))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("static api key", func(t *testing.T) {
		c := &Client{
			BaseURL:    u,
			APIKey:     apiKey,
			HTTPClient: ts.Client(),
		}
		_, err := c.do(context.Background(), http.MethodGet, "/foo/bar", nil, nil)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("dynamic api key", func(t *testing.T) {
		var cnt int32
		c := &Client{
			BaseURL: u,
			APIKeyProvider: APIKeyProviderFunc(func(ctx context.Context) (string, error) {
				atomic.AddInt32(&cnt, 1)
				time.Sleep(100 * time.Millisecond)
				return apiKey, nil
			}),
			HTTPClient: ts.Client(),
		}

		// send requests concurrently
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, err := c.do(context.Background(), http.MethodGet, "/foo/bar", nil, nil)
			if err != nil {
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			_, err := c.do(context.Background(), http.MethodGet, "/foo/bar", nil, nil)
			if err != nil {
				t.Error(err)
			}
		}()
		wg.Wait()

		if got := atomic.LoadInt32(&cnt); got != 1 {
			t.Errorf("unexpected call count of APIKeyProvider: want %d, got %d", 1, got)
		}
	})
}

func TestError(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error": {"message": "ERROR MESSAGE HERE"}}`)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		BaseURL:    u,
		APIKey:     "DUMMY-API-KEY",
		HTTPClient: ts.Client(),
	}

	_, err = c.do(context.Background(), http.MethodGet, "/foo/bar", nil, nil)
	merr, ok := err.(Error)
	if !ok {
		t.Errorf("want mackerel.Error, got %t", err)
		return
	}

	if merr.Message() != "ERROR MESSAGE HERE" {
		t.Errorf("unexpected error mesage: want %s, got %s", "ERROR MESSAGE HERE", merr.Message())
	}
	if merr.StatusCode() != http.StatusBadRequest {
		t.Errorf("unexpected status code: want %d, got %d", http.StatusBadRequest, merr.StatusCode())
	}
}
