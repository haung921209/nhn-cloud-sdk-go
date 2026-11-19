package networkacl

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

// The API names the rule fields src_ip, dst_ip and *_port_range_min/max;
// src_ip and dst_ip are required.
func TestACLRuleUsesAPIFieldNames(t *testing.T) {
	srv, rec := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"acl_rule":{"id":"r1","src_ip":"10.0.0.0/8","dst_ip":"10.0.1.0/24",
			"src_port_range_min":null,"dst_port_range_min":22,"dst_port_range_max":22,"order":200}}`)
	})
	defer srv.Close()
	c := &Client{httpClient: client.NewClient(srv.URL, testToken{})}

	port := 22
	out, err := c.CreateACLRule(context.Background(), &CreateACLRuleInput{
		ACLID: "acl", SrcIPPrefix: "10.0.0.0/8", DstIPPrefix: "10.0.1.0/24",
		DstPortMin: &port, DstPortMax: &port, Policy: "allow", OrderNum: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	var body struct {
		ACLRule map[string]interface{} `json:"acl_rule"`
	}
	if err := json.Unmarshal(rec.Body, &body); err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"src_ip": "10.0.0.0/8", "dst_ip": "10.0.1.0/24",
		"dst_port_range_min": 22.0, "dst_port_range_max": 22.0}
	for k, v := range want {
		if body.ACLRule[k] != v {
			t.Errorf("request %s = %v, want %v (body %s)", k, body.ACLRule[k], v, rec.Body)
		}
	}
	r := out.ACLRule
	if r.SrcIPPrefix != "10.0.0.0/8" || r.DstIPPrefix != "10.0.1.0/24" || r.DstPortMin == nil || *r.DstPortMin != 22 {
		t.Errorf("decoded rule: %+v", r)
	}
}
