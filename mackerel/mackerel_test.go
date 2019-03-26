package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestError(t *testing.T) {
	const apiKey = "DUMMY-APi-KEY"
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != apiKey {
			t.Errorf("unexpected api key, want %s, got %s", apiKey, r.Header.Get("X-Api-Key"))
		}
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
		APIKey:     apiKey,
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
