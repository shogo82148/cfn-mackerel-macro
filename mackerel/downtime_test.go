package mackerel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateDowntime(t *testing.T) {
	const (
		downtimeID = "9rxGOHfVF8F"
	)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		want := map[string]interface{}{
			"name":     "downtime name",
			"memo":     "memo",
			"start":    1234567890.0,
			"duration": 10.0,
		}
		if diff := cmp.Diff(data, want); diff != "" {
			t.Errorf("downtime differs: (-got +want)\n%s", diff)
		}
		w.WriteHeader(http.StatusOK)
		want["id"] = downtimeID
		json.NewEncoder(w).Encode(want)
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

	param := &Downtime{
		Name:     "downtime name",
		Memo:     "memo",
		Start:    1234567890,
		Duration: 10,
	}
	got, err := c.CreateDowntime(context.Background(), param)
	if err != nil {
		t.Error(err)
	}

	param.ID = downtimeID
	if diff := cmp.Diff(got, param); diff != "" {
		t.Errorf("downtime differs: (-got +want)\n%s", diff)
	}
}
