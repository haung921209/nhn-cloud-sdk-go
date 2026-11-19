package postgresql

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

const okHeader = `"header":{"resultCode":0,"resultMessage":"SUCCESS","isSuccessful":true}`

// newTestClient points a Client at an httptest TLS server (core.Client always
// uses https).
func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	tr := http.DefaultTransport.(*http.Transport)
	old := tr.TLSClientConfig
	tr.TLSClientConfig = srv.Client().Transport.(*http.Transport).TLSClientConfig
	t.Cleanup(func() { tr.TLSClientConfig = old })
	return &Client{core: core.NewClient(srv.Listener.Addr().String(), nil, nil)}
}

func TestListStorageTypes_StringList(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{` + okHeader + `,"storageTypes":["General SSD","General HDD"]}`))
	})
	got, err := c.ListStorageTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.StorageTypes, []string{"General SSD", "General HDD"}) {
		t.Errorf("storageTypes = %v", got.StorageTypes)
	}
}

func TestGetParameterGroup_FlatBody(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{` + okHeader + `,"parameterGroupId":"pg-1","parameterGroupName":"name-1",
			"dbVersion":"POSTGRESQL_V146","parameterGroupStatus":"STABLE",
			"parameters":[{"parameterName":"checkpoint_timeout","value":"300s","allowedValue":"30~86400s","valueType":"NUMERIC_WITH_TIME_UNIT"}],
			"createdYmdt":"2026-01-01T00:00:00+09:00"}`))
	})
	got, err := c.GetParameterGroup(context.Background(), "pg-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ParameterGroupName != "name-1" || got.ParameterGroupStatus != "STABLE" || got.CreatedAt == "" {
		t.Errorf("flat fields not decoded: %+v", got.ParameterGroup)
	}
	if len(got.Parameters) != 1 || got.Parameters[0].AllowedValues != "30~86400s" || got.Parameters[0].DataType == "" {
		t.Errorf("parameters = %+v", got.Parameters)
	}
	b, _ := json.Marshal(got)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m["parameterGroupName"] != "name-1" {
		t.Errorf("JSON output not flat: %s", b)
	}
}

func TestListBackups_SendsPageAndSize(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != "/v1.0/backups" || q.Get("page") != "2" || q.Get("size") != "5" || q.Get("dbInstanceId") != "i-1" {
			w.WriteHeader(400)
			w.Write([]byte(`{"header":{"resultCode":80400,"resultMessage":"bad query ` + r.URL.RawQuery + `","isSuccessful":false}}`))
			return
		}
		w.Write([]byte(`{` + okHeader + `,"totalCounts":1,"backups":[{"backupId":"b-1","createdYmdt":"2026-01-01T00:00:00+09:00","startYmdt":"s","completedYmdt":"c"}]}`))
	})
	got, err := c.ListBackups(context.Background(), "i-1", 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalCounts != 1 || len(got.Backups) != 1 || got.Backups[0].CreatedAt == "" || got.Backups[0].BackupStartedAt == "" || got.Backups[0].BackupCompletedAt == "" {
		t.Errorf("backups = %+v", got)
	}
}

func TestGetMetricStatistics_QueryAndBody(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("from") != "2026-11-17T00:00:00.000+09:00" || q.Get("to") != "2026-11-17T01:00:00.000+09:00" ||
			!reflect.DeepEqual(q["metricNames"], []string{"CPU_USAGE", "MEMORY_USAGE"}) ||
			q.Get("interval") != "5" || q.Get("dbInstanceId") != "i-1" {
			w.WriteHeader(400)
			w.Write([]byte(`{"header":{"resultCode":80400,"resultMessage":"bad query ` + r.URL.RawQuery + `","isSuccessful":false}}`))
			return
		}
		w.Write([]byte(`{` + okHeader + `,"metricStatistics":[{"metricName":"CPU_USAGE","unit":"%","values":[{"timestamp":1679298540,"value":7.5}]}]}`))
	})
	got, err := c.GetMetricStatistics(context.Background(), "i-1", []string{"CPU_USAGE", "MEMORY_USAGE"},
		"2026-11-17T00:00:00.000+09:00", "2026-11-17T01:00:00.000+09:00", 5)
	if err != nil {
		t.Fatal(err)
	}
	s := got.MetricStatistics
	if len(s) != 1 || s[0].MetricName != "CPU_USAGE" || len(s[0].Values) != 1 || s[0].Values[0].Timestamp != 1679298540 || s[0].Values[0].Value != 7.5 {
		t.Errorf("metricStatistics = %+v", s)
	}
}

func TestListMetrics_MetricName(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{` + okHeader + `,"metrics":[{"metricName":"CPU_USAGE","unit":"%"}]}`))
	})
	got, err := c.ListMetrics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Metrics) != 1 || got.Metrics[0].MetricName != "CPU_USAGE" {
		t.Errorf("metrics = %+v", got.Metrics)
	}
}

func TestGetInstance_YmdtFields(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{` + okHeader + `,"dbInstanceId":"i-1","createdYmdt":"2026-01-01T00:00:00+09:00","updatedYmdt":"2026-01-02T00:00:00+09:00"}`))
	})
	got, err := c.GetInstance(context.Background(), "i-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.CreatedAt != "2026-01-01T00:00:00+09:00" || got.UpdatedAt != "2026-01-02T00:00:00+09:00" {
		t.Errorf("created/updated = %q / %q", got.CreatedAt, got.UpdatedAt)
	}
}

// POST /db-security-groups/{id}/rules answers with a jobId, not a rule ID.
func TestCreateSecurityRule_JobID(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{` + okHeader + `,"jobId":"job-1"}`))
	})
	got, err := c.CreateSecurityRule(context.Background(), "sg-1", &CreateSecurityRuleRequest{Direction: "INGRESS", EtherType: "IPV4", CIDR: "10.0.0.0/24"})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(got)
	if !strings.Contains(string(b), `"jobId":"job-1"`) {
		t.Errorf("response JSON %s", b)
	}
}

// Notification groups: notifyEmail/notifySms are booleans, dbInstanceIds and
// userGroupIds are required lists, and GET returns the group at the top level.
func TestNotificationGroup_Shapes(t *testing.T) {
	var body map[string]interface{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			json.NewDecoder(r.Body).Decode(&body)
			w.Write([]byte(`{` + okHeader + `,"notificationGroupId":"ng-1"}`))
			return
		}
		w.Write([]byte(`{` + okHeader + `,"notificationGroupId":"ng-1","notificationGroupName":"n","notifyEmail":false,"notifySms":true,"isEnabled":true}`))
	})
	if _, err := c.CreateNotificationGroup(context.Background(), &CreateNotificationGroupRequest{NotificationGroupName: "n"}); err != nil {
		t.Fatal(err)
	}
	if body["dbInstanceIds"] == nil || body["userGroupIds"] == nil {
		t.Errorf("create body %v", body)
	}
	got, err := c.GetNotificationGroup(context.Background(), "ng-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.NotificationGroupName != "n" || !got.NotifySms || !got.IsEnabled {
		t.Errorf("get %+v", got.NotificationGroup)
	}
}
