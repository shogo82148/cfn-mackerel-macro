package mackerel

import (
	"context"
	"encoding/json"
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
		if r.Method != http.MethodGet {
			t.Errorf("invalid method: want %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations", r.URL.Path)
		}
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
					Role:            nil,
					ExcludedMetrics: []string{},
				},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}

func TestFindAWSIntegration(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("invalid method: want %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations/46vGJ7uUsp3" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations/46vGJ7uUsp3", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
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
}`)
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

	integrations, err := c.FindAWSIntegration(context.Background(), "46vGJ7uUsp3")
	if err != nil {
		t.Error(err)
	}
	want := &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}

func TestCreateAWSIntegration(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("invalid method: want %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations", r.URL.Path)
		}
		dec := json.NewDecoder(r.Body)
		var got interface{}
		if err := dec.Decode(&got); err != nil {
			t.Error(err)
		}
		want := map[string]interface{}{
			"name":         "shogo82148",
			"memo":         "",
			"externalId":   "hogehoge",
			"region":       "ap-northeast-1",
			"includedTags": "",
			"excludedTags": "",
			"roleArn":      "arn:aws:iam::123456789012:role/foobar",
			"services": map[string]interface{}{
				"S3": map[string]interface{}{
					"enable":          false,
					"role":            nil,
					"excludedMetrics": []interface{}{},
				},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("request body missmatch: (-want/+got)\n%s", diff)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
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
}`)
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

	integrations, err := c.CreateAWSIntegration(context.Background(), &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	want := &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}

func TestUpdateAWSIntegration(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("invalid method: want %s, got %s", http.MethodPut, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations/46vGJ7uUsp3" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations/46vGJ7uUsp3", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
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
}`)
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

	integrations, err := c.UpdateAWSIntegration(context.Background(), "46vGJ7uUsp3", &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	want := &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}

func TestDeleteAWSIntegration(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("invalid method: want %s, got %s", http.MethodDelete, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations/46vGJ7uUsp3" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations/46vGJ7uUsp3", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
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
}`)
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

	integrations, err := c.DeleteAWSIntegration(context.Background(), "46vGJ7uUsp3")
	if err != nil {
		t.Error(err)
	}
	want := &AWSIntegration{
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
				Role:            nil,
				ExcludedMetrics: []string{},
			},
		},
	}
	if diff := cmp.Diff(want, integrations); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}

func TestCreateAWSIntegrationExternalID(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("invalid method: want %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations-external-id" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations-external-id", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"externalId": "hogehoge"}`)
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

	externalID, err := c.CreateAWSIntegrationExternalID(context.Background())
	if err != nil {
		t.Error(err)
	}
	if externalID != "hogehoge" {
		t.Errorf("want %q, got %q", "hogehoge", externalID)
	}
}

func TestFindAWSIntegrationsExcludableMetrics(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("invalid method: want %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v0/aws-integrations-excludable-metrics" {
			t.Errorf("invalid path: want %s, got %s", "/api/v0/aws-integrations-excludable-metrics", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"S3": [
				"s3.requests.all_requests",
				"s3.requests.get_requests",
				"s3.requests.put_requests",
				"s3.requests.delete_requests",
				"s3.requests.head_requests",
				"s3.requests.post_requests",
				"s3.requests.select_requests",
				"s3.requests.list_requests",
				"s3.select_bytes.scanned",
				"s3.select_bytes.returned",
				"s3.bytes.downloaded",
				"s3.bytes.uploaded",
				"s3.errors.4xx",
				"s3.errors.5xx",
				"s3.latency.#.minimum",
				"s3.latency.#.average",
				"s3.latency.#.maximum",
				"s3.bucket_size.*",
				"s3.number_of_objects.count"
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

	got, err := c.FindAWSIntegrationsExcludableMetrics(context.Background())
	if err != nil {
		t.Error(err)
	}
	want := map[string][]string{
		"S3": {
			"s3.requests.all_requests",
			"s3.requests.get_requests",
			"s3.requests.put_requests",
			"s3.requests.delete_requests",
			"s3.requests.head_requests",
			"s3.requests.post_requests",
			"s3.requests.select_requests",
			"s3.requests.list_requests",
			"s3.select_bytes.scanned",
			"s3.select_bytes.returned",
			"s3.bytes.downloaded",
			"s3.bytes.uploaded",
			"s3.errors.4xx",
			"s3.errors.5xx",
			"s3.latency.#.minimum",
			"s3.latency.#.average",
			"s3.latency.#.maximum",
			"s3.bucket_size.*",
			"s3.number_of_objects.count",
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want/+got):\n%s", diff)
	}
}
