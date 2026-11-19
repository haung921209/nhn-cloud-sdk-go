package mariadb

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
	body         map[string]interface{}
}

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// fakeRequestClient sends the client's requests to an httptest server that
// answers with a job response and records the last request.
func fakeRequestClient(t *testing.T) (*Client, *seenRequest) {
	t.Helper()
	seen := &seenRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.method, seen.path, seen.body = r.Method, r.URL.Path, nil
		if b, _ := io.ReadAll(r.Body); len(b) > 0 {
			json.Unmarshal(b, &seen.body)
		}
		io.WriteString(w, `{"header":{"resultCode":0,"isSuccessful":true},"jobId":"job-1","userGroupId":"ug-1"}`)
	}))
	t.Cleanup(srv.Close)
	target, _ := url.Parse(srv.URL)
	old := http.DefaultTransport
	http.DefaultTransport = rtFunc(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		r.URL.Scheme, r.URL.Host = target.Scheme, target.Host
		return old.RoundTrip(r)
	})
	t.Cleanup(func() { http.DefaultTransport = old })
	return &Client{core: core.NewClient("rds.example.com", nil, nil)}, seen
}

// Guide v4.0: the same request shapes as MySQL for these calls.
func TestMariaDBRequestShapes(t *testing.T) {
	ctx := context.Background()
	c, seen := fakeRequestClient(t)

	if _, err := c.CreateBackup(ctx, "inst-1", &CreateBackupRequest{BackupName: "bk", BackupMethodType: "FULL"}); err != nil {
		t.Fatal(err)
	}
	if seen.method != "POST" || seen.path != "/v4.0/backups" || seen.body["backupMethodType"] != "FULL" || seen.body["dbInstanceId"] != "inst-1" {
		t.Errorf("backup: %s %s %v", seen.method, seen.path, seen.body)
	}

	if _, err := c.EnableHA(ctx, "inst-1", &EnableHARequest{UseHighAvailability: true, DBInstanceCandidateName: "sb"}); err != nil {
		t.Fatal(err)
	}
	if seen.body["dbInstanceCandidateName"] != "sb" {
		t.Errorf("enable HA body %v", seen.body)
	}

	if _, err := c.CreateReplica(ctx, "inst-1", &CreateReplicaRequest{DBInstanceName: "rr", Network: ReplicaNetworkConfig{AvailabilityZone: "kr-pub-a"}}); err != nil {
		t.Fatal(err)
	}
	if n, _ := seen.body["network"].(map[string]interface{}); n["availabilityZone"] != "kr-pub-a" {
		t.Errorf("replica body %v", seen.body)
	}

	if _, err := c.CreateUserGroup(ctx, &CreateUserGroupRequest{UserGroupName: "ug"}); err != nil {
		t.Fatal(err)
	}
	if ids, ok := seen.body["memberIds"].([]interface{}); !ok || len(ids) != 0 {
		t.Errorf("user group memberIds must be sent as an empty array: %v", seen.body)
	}
}
