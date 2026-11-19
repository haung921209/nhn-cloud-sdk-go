package mariadb

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// newTestClient points a client at a TLS test server that answers every
// request with body and records the request URL.
func newTestClient(t *testing.T, body string) (*Client, *url.URL) {
	t.Helper()
	var got url.URL
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = *r.URL
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	old := http.DefaultTransport
	http.DefaultTransport = srv.Client().Transport
	defer func() { http.DefaultTransport = old }()
	return &Client{core: core.NewClient(strings.TrimPrefix(srv.URL, "https://"), nil, nil)}, &got
}

const okHeader = `"header":{"isSuccessful":true,"resultCode":0,"resultMessage":"SUCCESS"}`

func TestListStorageTypesDecodesStrings(t *testing.T) {
	c, _ := newTestClient(t, `{`+okHeader+`,"storageTypes":["General SSD","General HDD"]}`)
	res, err := c.ListStorageTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := json.Marshal(res); !strings.Contains(string(b), `"storageTypes":["General SSD","General HDD"]`) {
		t.Errorf("got %s", b)
	}
}

func TestGetParameterGroupReadsTopLevelFields(t *testing.T) {
	c, _ := newTestClient(t, `{`+okHeader+`,"parameterGroupId":"pg-1","parameterGroupName":"n","dbVersion":null,
		"parameters":[{"parameterId":"p-1","parameterName":"wait_timeout","value":"28800"}]}`)
	res, err := c.GetParameterGroup(context.Background(), "pg-1")
	if err != nil {
		t.Fatal(err)
	}
	// The API sends the group's fields at the top level, next to header.
	b, _ := json.Marshal(res)
	if !strings.Contains(string(b), `"parameterGroupId":"pg-1"`) || !strings.Contains(string(b), `"parameterId":"p-1"`) {
		t.Errorf("got %s", b)
	}
}

func TestGetMetricStatisticsQueryAndValues(t *testing.T) {
	c, got := newTestClient(t, `{`+okHeader+`,"metricStatistics":[{"measureName":"CPU_USAGE","unit":"%","values":[[1679298540,"7.5%"]]}]}`)
	res, err := c.GetMetricStatistics(context.Background(), "i-1", []string{"CPU_USAGE", "MEMORY_USAGE"},
		"2026-11-18T00:00:00.000+09:00", "2026-11-18T01:00:00.000+09:00", 0)
	if err != nil {
		t.Fatal(err)
	}
	q := got.Query()
	if q.Get("from") != "2026-11-18T00:00:00.000+09:00" || q.Get("dbInstanceId") != "i-1" {
		t.Errorf("query not encoded: %s", got.RawQuery)
	}
	if m := q["measureNames"]; len(m) != 2 || m[1] != "MEMORY_USAGE" {
		t.Errorf("measureNames = %v", m)
	}
	if q.Has("interval") {
		t.Errorf("interval sent when unset: %s", got.RawQuery)
	}
	if v := res.MetricStatistics[0].Values; len(v) != 1 || v[0][1] != "7.5%" {
		t.Errorf("values = %v", v)
	}
}

// The API field is selectAll; an unknown selectAllYN is ignored and the
// group is refused as empty (API error 17002).
func TestCreateUserGroupSendsSelectAll(t *testing.T) {
	b, _ := json.Marshal(CreateUserGroupRequest{UserGroupName: "g", SelectAllYN: true})
	if !strings.Contains(string(b), `"selectAll":true`) {
		t.Errorf("body %s", b)
	}
}
