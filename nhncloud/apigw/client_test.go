package apigw

import (
	"context"
	"encoding/json"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/auth"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// rewrite sends every request to the test server.
type rewrite struct{ to *url.URL }

func (r rewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme, req.URL.Host = r.to.Scheme, r.to.Host
	return http.DefaultTransport.RoundTrip(req)
}

// The API answers HTTP 200 with header.isSuccessful=false on failure.
func TestFailedHeaderIsAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"header":{"isSuccessful":false,"resultCode":4010001,"resultMessage":"Invalid appKey."}}`))
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	c := NewClient("kr1", "example-app-key", "", "", &http.Client{Transport: rewrite{u}}, false)
	_, err := c.ListServices(context.Background())
	if err == nil || !strings.Contains(err.Error(), "4010001") || !strings.Contains(err.Error(), "Invalid appKey.") {
		t.Fatalf("want an error with the result code and message, got %v", err)
	}
}

// API v2.0 takes a User Access Key token, not the X-TC-AUTHENTICATION pair.
func TestRequestsUseV2AndBearerToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".nhncloud"), 0700)
	cache, _ := json.Marshal(map[string]interface{}{"token": "tok", "expires_at": time.Now().Add(time.Hour)})
	os.WriteFile(auth.TokenCacheFile(home, "id"), cache, 0600)

	var got *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		w.Write([]byte(`{"header":{"isSuccessful":true,"resultCode":0}}`))
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	c := NewClient("kr1", "example-app-key", "id", "secret", &http.Client{Transport: rewrite{u}}, false)
	if _, err := c.ListServices(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got.URL.Path != "/v2.0/appkeys/example-app-key/services" || got.Header.Get("X-NHN-AUTHORIZATION") != "Bearer tok" || got.Header.Get("X-TC-AUTHENTICATION-SECRET") != "" {
		t.Errorf("request %s, headers %v", got.URL.Path, got.Header)
	}
}

// The service list is "apigwServiceList" (API guide).
func TestListServicesDecodesServiceList(t *testing.T) {
	var out ListServicesOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"apigwServiceList":[{"apigwServiceId":"s1","apigwServiceName":"n"}]}`), &out)
	if len(out.Services) != 1 || out.Services[0].ID != "s1" {
		t.Fatalf("decoded %+v", out)
	}
}
