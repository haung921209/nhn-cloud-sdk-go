package ncs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	iclient "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type staticToken struct{}

func (staticToken) GetToken(context.Context) (string, error) { return "t", nil }
func (staticToken) SetAuthHeader(r *http.Request, tok string) {
	r.Header.Set("x-nhn-authorization", "Bearer "+tok)
}

func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &Client{httpClient: iclient.NewClient(srv.URL, staticToken{})}
}

const okHeader = `"header":{"isSuccessful":true,"resultCode":0,"resultMessage":"SUCCESS"}`

// The API wants {"workload": {name, templateId, desired, loadBalancing}} and
// answers {"header":..., "workload": {...}}.
func TestCreateWorkloadWrapsBodyAndDecodesWorkload(t *testing.T) {
	var got map[string]map[string]interface{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/workloads" {
			t.Errorf("request %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &got)
		io.WriteString(w, `{`+okHeader+`,"workload":{"id":"wl-1","name":"web","templateId":"tpl-1","desired":1,"status":"Creating"}}`)
	})
	out, err := c.CreateWorkload(context.Background(), &CreateWorkloadInput{
		Name: "web", TemplateID: "tpl-1", Desired: 1, LoadBalancing: &LoadBalancing{},
	})
	if err != nil {
		t.Fatal(err)
	}
	wl := got["workload"]
	if wl == nil || wl["name"] != "web" || wl["templateId"] != "tpl-1" || wl["desired"] != float64(1) {
		t.Errorf("body %v", got)
	}
	if lb, _ := wl["loadBalancing"].(map[string]interface{}); lb == nil || lb["enabled"] != false || lb["floatingIp"] != false {
		t.Errorf("loadBalancing %v", wl["loadBalancing"])
	}
	if out.Workload.ID != "wl-1" || out.Workload.TemplateID != "tpl-1" {
		t.Errorf("decoded %+v", out.Workload)
	}
}

func TestGetWorkloadAndTemplateDecodeWrappedObject(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/workloads/wl-1":
			io.WriteString(w, `{`+okHeader+`,"workload":{"id":"wl-1","name":"web","desired":2,"available":1}}`)
		case "/templates/tpl-1":
			io.WriteString(w, `{`+okHeader+`,"template":{"id":"tpl-1","name":"nginx","version":"v1"}}`)
		default:
			t.Errorf("path %s", r.URL.Path)
		}
	})
	wl, err := c.GetWorkload(context.Background(), "wl-1")
	if err != nil {
		t.Fatal(err)
	}
	if wl.Workload.ID != "wl-1" || wl.Workload.Desired != 2 || wl.Workload.Available != 1 {
		t.Errorf("workload %+v", wl.Workload)
	}
	tpl, err := c.GetTemplate(context.Background(), "tpl-1")
	if err != nil {
		t.Fatal(err)
	}
	if tpl.Template.ID != "tpl-1" || tpl.Template.Name != "nginx" || tpl.Template.Version != "v1" {
		t.Errorf("template %+v", tpl.Template)
	}
}
