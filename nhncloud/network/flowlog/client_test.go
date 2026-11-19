package flowlog

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

type testToken struct{}

func (testToken) GetToken(context.Context) (string, error)    { return "tok", nil }
func (testToken) SetAuthHeader(r *http.Request, token string) { r.Header.Set("X-Auth-Token", token) }

// The API takes filter and connection_setup_only, returns status, and writes
// dates as "2006-01-02 15:04:05".
func TestLoggerUsesAPIFieldNames(t *testing.T) {
	srv, rec := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"flowlog_logger":{"id":"l1","filter":"ACCEPT","connection_setup_only":true,
			"status":"BUILD","created_at":"2024-08-04 04:46:33","updated_at":"2024-08-04 04:46:33"}}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	out, err := c.CreateLogger(context.Background(), &CreateLoggerInput{
		Name: "l", ResourceType: "SUBNET", ResourceID: "s", Filter: "ACCEPT",
		ConnectionSetupOnly: true, StorageType: "OBS", StorageURL: "https://obs/c", LogFormat: "CSV",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body struct {
		Logger map[string]interface{} `json:"flowlog_logger"`
	}
	if err := json.Unmarshal(rec.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.Logger["filter"] != "ACCEPT" || body.Logger["connection_setup_only"] != true {
		t.Errorf("request body: %s", rec.Body)
	}
	for _, k := range []string{"filter_type", "connection_action"} {
		if _, ok := body.Logger[k]; ok {
			t.Errorf("request sends %s, which the API does not have", k)
		}
	}
	l := out.Logger
	if l.Filter != "ACCEPT" || !l.ConnectionSetupOnly || l.Status != "BUILD" || l.CreatedAt != "2024-08-04 04:46:33" {
		t.Errorf("decoded logger: %+v", l)
	}
}

func TestLoggingPortDecodesAPIDate(t *testing.T) {
	srv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"flowlog_logging_ports":[{"id":"p1","logger_id":"l1","port_id":"x","network_id":"n","created_at":"2024-07-29 09:21:09"}]}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	out, err := c.ListLoggingPorts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.LoggingPorts) != 1 || out.LoggingPorts[0].NetworkID != "n" {
		t.Errorf("decoded ports: %+v", out.LoggingPorts)
	}
}
