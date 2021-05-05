package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/pointer"
)

func TestFindAWSIntegrations(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"aws_integrations": [
{
	"id": "46vGJ7uUsp3",
	"name": "shogo82148",
	"memo": "",
	"key": null,
	"roleArn": "arn:aws:iam::123456789012:role/foobar",
	"externalId": "hogehoge",
	"region": "ap-northeast-1",
	"includedTags": "",
	"excludedTags": "",
	"services": {
		"S3": {
			"enable": false,
			"role": null,
			"excludedMetrics": []
		}
	}
}
]}`)
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

	integrations, err := c.FindAWSIntegrations(context.Background())
	if err != nil {
		t.Error(err)
	}
	want := []*AWSIntegration{
		{
			ID:           "46vGJ7uUsp3",
			Name:         "shogo82148",
			Memo:         "",
			RoleArn:      pointer.String("arn:aws:iam::123456789012:role/foobar"),
			ExternalID:   pointer.String("hogehoge"),
			Region:       "ap-northeast-1",
			IncludedTags: "",
			ExcludedTags: "",
			Services: map[string]*AWSIntegrationService{
				"S3": {
					Enable:          false,
					Role:            "",
					ExcludedMetrics: []string{},
				},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}
