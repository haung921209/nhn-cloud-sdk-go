package mysql

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

type seenRequest struct {
	method, path string
	query        url.Values
	body         map[string]interface{}
}

// fakeClient returns a Client whose requests go to an httptest server that
// answers with reply and records the last request.
func fakeClient(t *testing.T, reply string) (*Client, *seenRequest) {
	t.Helper()
	seen := &seenRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.method, seen.path, seen.query = r.Method, r.URL.Path, r.URL.Query()
		seen.body = nil
		if b, _ := io.ReadAll(r.Body); len(b) > 0 {
			json.Unmarshal(b, &seen.body)
		}
		io.WriteString(w, reply)
	}))
	t.Cleanup(srv.Close)
	target, _ := url.Parse(srv.URL)
	old := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		r.URL.Scheme, r.URL.Host = target.Scheme, target.Host
		return old.RoundTrip(r)
	})
	t.Cleanup(func() { http.DefaultTransport = old })
	return &Client{core: core.NewClient("rds.example.com", nil, nil)}, seen
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

const okJob = `{"header":{"resultCode":0,"resultMessage":"SUCCESS","isSuccessful":true},"jobId":"job-1"}`

// v4.0 creates a backup with POST /v4.0/backups and the instance in the body.
func TestCreateBackupUsesV4BackupsPath(t *testing.T) {
	c, seen := fakeClient(t, okJob)
	_, err := c.CreateBackup(context.Background(), "inst-1", &CreateBackupRequest{BackupName: "bk", BackupMethodType: "FULL"})
	if err != nil {
		t.Fatal(err)
	}
	if seen.method != "POST" || seen.path != "/v4.0/backups" {
		t.Errorf("request %s %s", seen.method, seen.path)
	}
	want := map[string]interface{}{"backupName": "bk", "backupMethodType": "FULL", "dbInstanceId": "inst-1"}
	for k, v := range want {
		if seen.body[k] != v {
			t.Errorf("body %s = %v, want %v (body %v)", k, seen.body[k], v, seen.body)
		}
	}
}

// Enabling HA needs the standby's name.
func TestEnableHASendsCandidateName(t *testing.T) {
	c, seen := fakeClient(t, okJob)
	_, err := c.EnableHA(context.Background(), "inst-1", &EnableHARequest{UseHighAvailability: true, DBInstanceCandidateName: "standby"})
	if err != nil {
		t.Fatal(err)
	}
	if seen.body["dbInstanceCandidateName"] != "standby" {
		t.Errorf("body %v", seen.body)
	}
}

// A replica needs network.availabilityZone.
func TestCreateReplicaSendsAvailabilityZone(t *testing.T) {
	c, seen := fakeClient(t, okJob)
	_, err := c.CreateReplica(context.Background(), "inst-1", &CreateReplicaRequest{
		DBInstanceName: "rr", Network: ReplicaNetworkConfig{AvailabilityZone: "kr-pub-b"}})
	if err != nil {
		t.Fatal(err)
	}
	if seen.path != "/v4.0/db-instances/inst-1/replicate" {
		t.Errorf("path %s", seen.path)
	}
	nw, _ := seen.body["network"].(map[string]interface{})
	if nw["availabilityZone"] != "kr-pub-b" {
		t.Errorf("body %v", seen.body)
	}
}

// measureNames is a required array query parameter; from/to keep their '+'
// offset; values come back as [timestamp, "value"] pairs.
func TestGetMetricStatisticsQueryAndValues(t *testing.T) {
	c, seen := fakeClient(t, `{"header":{"resultCode":0,"isSuccessful":true},
		"metricStatistics":[{"measureName":"CPU_USAGE","unit":"%","values":[[1679298540,"7.5%"]]}]}`)
	res, err := c.GetMetricStatistics(context.Background(), "inst-1", []string{"CPU_USAGE", "MEMORY_USAGE"},
		"2026-11-18T00:00:00+09:00", "2026-11-18T01:00:00+09:00", 0)
	if err != nil {
		t.Fatal(err)
	}
	q := seen.query
	if got := q["measureNames"]; len(got) != 2 || got[0] != "CPU_USAGE" || got[1] != "MEMORY_USAGE" {
		t.Errorf("measureNames %v", got)
	}
	if q.Get("from") != "2026-11-18T00:00:00+09:00" || q.Get("to") != "2026-11-18T01:00:00+09:00" {
		t.Errorf("from/to %q %q", q.Get("from"), q.Get("to"))
	}
	if q.Get("dbInstanceId") != "inst-1" || q.Has("interval") {
		t.Errorf("query %v", q)
	}
	v := res.MetricStatistics[0].Values
	if len(v) != 1 || v[0].Timestamp != 1679298540 || v[0].Value != "7.5%" {
		t.Errorf("values %+v", v)
	}
}

// selectAll (not selectAllYN) adds every project member; memberIds is a
// required array, so it is sent even when empty.
func TestCreateUserGroupBody(t *testing.T) {
	c, seen := fakeClient(t, `{"header":{"resultCode":0,"isSuccessful":true},"userGroupId":"ug-1"}`)
	if _, err := c.CreateUserGroup(context.Background(), &CreateUserGroupRequest{UserGroupName: "ug", SelectAllYN: true}); err != nil {
		t.Fatal(err)
	}
	if seen.body["selectAll"] != true {
		t.Errorf("body %v", seen.body)
	}
	if m, ok := seen.body["memberIds"].([]interface{}); !ok || len(m) != 0 {
		t.Errorf("memberIds %v", seen.body["memberIds"])
	}
}
