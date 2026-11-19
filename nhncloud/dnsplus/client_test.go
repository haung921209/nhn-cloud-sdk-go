package dnsplus

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The API answers HTTP 200 with header.isSuccessful=false on failure.
func TestFailedHeaderIsAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"header":{"isSuccessful":false,"resultCode":4010001,"resultMessage":"Invalid appKey."}}`))
	}))
	defer srv.Close()
	c := NewClient("example-app-key", nil, false)
	c.baseURL = srv.URL
	_, err := c.ListZones(context.Background())
	if err == nil || !strings.Contains(err.Error(), "4010001") || !strings.Contains(err.Error(), "Invalid appKey.") {
		t.Fatalf("want an error with the result code and message, got %v", err)
	}
}

// Pools are top-level resources with their endpoints in the pool body, and
// deletes take the IDs in the query string (API guide).
func TestPoolRequestsFollowTheAPI(t *testing.T) {
	var reqs []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		reqs = append(reqs, r.Method+" "+r.URL.Path+"?"+r.URL.RawQuery+" "+string(b))
		w.Write([]byte(`{"header":{"isSuccessful":true},"pool":{"poolId":"p1","endpointList":[{"endpointAddress":"192.0.2.1","endpointWeight":1.0}]}}`))
	}))
	defer srv.Close()
	c := NewClient("ak", nil, false)
	c.baseURL = srv.URL
	ctx := context.Background()
	out, err := c.CreatePool(ctx, &CreatePoolInput{PoolName: "p", EndpointList: []Endpoint{{EndpointAddress: "192.0.2.1"}}})
	if err != nil || out.Pool.PoolID != "p1" || out.Pool.EndpointList[0].EndpointAddress != "192.0.2.1" {
		t.Fatalf("create: %v %+v", err, out)
	}
	c.ListPools(ctx)
	c.DeletePools(ctx, []string{"p1", "p2"})
	c.DeleteZones(ctx, []string{"z1"})
	want := []string{
		`POST /dnsplus/v1.0/appkeys/ak/pools? {"pool":{"poolName":"p","endpointList":[{"endpointAddress":"192.0.2.1"}]}}`,
		`GET /dnsplus/v1.0/appkeys/ak/pools? `,
		`DELETE /dnsplus/v1.0/appkeys/ak/pools?poolIdList=p1%2Cp2 `,
		`DELETE /dnsplus/v1.0/appkeys/ak/zones/async?zoneIdList=z1 `,
	}
	if strings.Join(reqs, "\n") != strings.Join(want, "\n") {
		t.Errorf("requests:\n%s\nwant:\n%s", strings.Join(reqs, "\n"), strings.Join(want, "\n"))
	}
}

// Record set and GSLB fields carry their prefixes (recordsetTtl, gslbTtl,
// gslbRoutingRule) in the API.
func TestRecordSetAndGSLBFieldNames(t *testing.T) {
	b, _ := json.Marshal(CreateRecordSetInput{RecordSetName: "a.", RecordSetType: "A", TTL: 300})
	g, _ := json.Marshal(CreateGSLBInput{GslbName: "g", RoutingType: "FAILOVER", TTL: 300})
	if !strings.Contains(string(b), `"recordsetTtl":300`) || !strings.Contains(string(g), `"gslbTtl":300`) || !strings.Contains(string(g), `"gslbRoutingRule":"FAILOVER"`) {
		t.Errorf("record set %s, gslb %s", b, g)
	}
	var out ListGSLBsOutput
	json.Unmarshal([]byte(`{"gslbList":[{"gslbId":"x","gslbTtl":300,"gslbRoutingRule":"GEOLOCATION"}]}`), &out)
	if out.GslbList[0].TTL != 300 || out.GslbList[0].RoutingType != "GEOLOCATION" {
		t.Errorf("gslb list %+v", out.GslbList)
	}
}
