package transithub

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

func requestObject(t *testing.T, body []byte, wrapper string) map[string]interface{} {
	t.Helper()
	var m map[string]map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatal(err)
	}
	return m[wrapper]
}

// The create body must carry the three required booleans, also when false.
func TestCreateTransitHubSendsRequiredBooleans(t *testing.T) {
	srv, rec := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"transithub":{"id":"h1","multicast_enable":true,"default_association_enable":false,"default_propagation_enable":true}}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	out, err := c.CreateTransitHub(context.Background(), &CreateTransitHubInput{Name: "h", MulticastEnable: true, DefaultPropagationEnable: true})
	if err != nil {
		t.Fatal(err)
	}
	got := requestObject(t, rec.Body, "transithub")
	want := map[string]interface{}{"multicast_enable": true, "default_association_enable": false, "default_propagation_enable": true}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("request %s = %v, want %v", k, got[k], v)
		}
	}
	if !out.TransitHub.MulticastEnable || out.TransitHub.DefaultAssociationEnable {
		t.Errorf("decoded hub: %+v", out.TransitHub)
	}
}

func TestCreateAttachmentSendsSubnet(t *testing.T) {
	srv, rec := testutil.NewTestHTTPServer(nil)
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	if _, err := c.CreateAttachment(context.Background(), &CreateAttachmentInput{Name: "a", TransitHubID: "h", ResourceID: "vpc", SubnetID: "sub"}); err != nil {
		t.Fatal(err)
	}
	got := requestObject(t, rec.Body, "transithub_attachment")
	if got["subnet_id"] != "sub" || got["resource_id"] != "vpc" || got["transithub_id"] != "h" {
		t.Errorf("request body: %s", rec.Body)
	}
	if _, ok := got["resource_type"]; ok {
		t.Errorf("request sends resource_type, which the API does not have: %s", rec.Body)
	}
}

func TestRoutingRuleUsesAPIFieldNames(t *testing.T) {
	srv, rec := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"transithub_routing_rule":{"id":"r1","cidr":"1.1.1.0/24","action":"FORWARD","rule_type":"STATIC","attachment_id":"a1","routing_table_id":"t1"}}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	out, err := c.CreateRoutingRule(context.Background(), &CreateRoutingRuleInput{RoutingTableID: "t1", CIDR: "1.1.1.0/24", AttachmentID: "a1"})
	if err != nil {
		t.Fatal(err)
	}
	got := requestObject(t, rec.Body, "transithub_routing_rule")
	if got["cidr"] != "1.1.1.0/24" || got["attachment_id"] != "a1" || got["routing_table_id"] != "t1" {
		t.Errorf("request body: %s", rec.Body)
	}
	if _, ok := got["action"]; ok {
		t.Errorf("empty action must be left to the API default (FORWARD): %s", rec.Body)
	}
	r := out.Rule
	if r.CIDR != "1.1.1.0/24" || r.Action != "FORWARD" || r.AttachmentID != "a1" || r.RuleType != "STATIC" {
		t.Errorf("decoded rule: %+v", r)
	}
}
