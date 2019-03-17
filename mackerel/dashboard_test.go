package mackerel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateDashboard(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(map[string]interface{}{
			"id":      "foobar",
			"title":   "title",
			"urlPath": "url path",
			"widgets": []map[string]interface{}{
				{
					"type":  "graph",
					"title": "graph title",
					"graph": map[string]string{
						"type":   "host",
						"hostId": "host-foobar",
						"name":   "host-graph",
					},
				},
			},
			"createdAt": 1234567890,
			"updatedAt": 1234567890,
		})
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		BaseURL:    u,
		HTTPClient: ts.Client(),
	}

	ret, err := c.CreateDashboard(context.Background(), &Dashboard{
		Title:   "title",
		Memo:    "memo",
		URLPath: "url path",
		Widgets: []Widget{},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
