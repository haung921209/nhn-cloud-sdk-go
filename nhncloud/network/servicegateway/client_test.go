package servicegateway

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

type testToken struct{}

func (testToken) GetToken(context.Context) (string, error)    { return "tok", nil }
func (testToken) SetAuthHeader(r *http.Request, token string) { r.Header.Set("X-Auth-Token", token) }

// The API returns the gateway address as fixed_ip, with its port_id.
func TestServiceGatewayDecodesFixedIPAndPort(t *testing.T) {
	srv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"servicegateway":{"id":"g1","fixed_ip":"192.168.0.82","port_id":"p1"}}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	out, err := c.GetServiceGateway(context.Background(), "g1")
	if err != nil {
		t.Fatal(err)
	}
	if out.ServiceGateway.FixedIP != "192.168.0.82" || out.ServiceGateway.PortID != "p1" {
		t.Errorf("decoded gateway: %+v", out.ServiceGateway)
	}
}
