package s3credential

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

type fakeTransport struct {
	hosts []string
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.hosts = append(f.hosts, r.URL.Host)
	body := `{"credentials":[{"access":"ak","secret":"sk","user_id":"u1","tenant_id":"t1","created_at":"2024-10-19T08:24:46.000000Z"}]}`
	if r.URL.Path == "/v2.0/tokens" {
		body = `{"access":{"token":{"id":"tok","expires":"` + time.Now().Add(time.Hour).Format(time.RFC3339) + `"}}}`
	}
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
}

// The credential API lives on the identity host, which has no region part,
// and the list comes wrapped in {"credentials": [...]}.
func TestListCredentialsUsesIdentityHostAndDecodesList(t *testing.T) {
	ft := &fakeTransport{}
	old := http.DefaultTransport
	http.DefaultTransport = ft
	defer func() { http.DefaultTransport = old }()

	c := NewClient("kr1", credentials.NewStaticIdentity("u", "p", "t"), nil, false)
	out, err := c.ListCredentials(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range ft.hosts {
		if h != "api-identity-infrastructure.nhncloudservice.com" {
			t.Errorf("request went to %q", h)
		}
	}
	if len(out.Credentials) != 1 || out.Credentials[0].Access != "ak" {
		t.Errorf("decoded %+v", out)
	}
}
