package testutil_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---- sample structs used across tests ----

type sampleResource struct {
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	Status       string `json:"status,omitempty"`
}

type missingTagField struct {
	Name string // deliberately no json tag
}

// ---- AssertStructHasJSONTag ----

func TestAssertStructHasJSONTag_Match(t *testing.T) {
	// Should pass silently when the tag matches.
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(sampleResource{}), "ResourceID", "resourceId")
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(sampleResource{}), "ResourceName", "resourceName")
}

func TestAssertStructHasJSONTag_OmitemptyStripped(t *testing.T) {
	// Options like ",omitempty" must be stripped before comparison.
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(sampleResource{}), "Status", "status")
}

func TestAssertStructHasJSONTag_PointerType(t *testing.T) {
	// Should work with pointer types as well.
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(&sampleResource{}), "ResourceID", "resourceId")
}

func TestAssertStructHasJSONTag_WrongTag(t *testing.T) {
	// Verify that a mismatch is reported via a fake *testing.T so the parent
	// test is not failed by the intentional mismatch.
	fake := &testing.T{}
	testutil.AssertStructHasJSONTag(fake, reflect.TypeOf(sampleResource{}), "ResourceID", "wrongTag")
	if !fake.Failed() {
		t.Error("expected fake test to fail on wrong tag, but it did not")
	}
}

func TestAssertStructHasJSONTag_MissingField(t *testing.T) {
	fake := &testing.T{}
	testutil.AssertStructHasJSONTag(fake, reflect.TypeOf(sampleResource{}), "NonExistentField", "nope")
	if !fake.Failed() {
		t.Error("expected fake test to fail on missing field, but it did not")
	}
}

func TestAssertStructHasJSONTag_NonStruct(t *testing.T) {
	fake := &testing.T{}
	testutil.AssertStructHasJSONTag(fake, reflect.TypeOf(""), "Field", "tag")
	if !fake.Failed() {
		t.Error("expected fake test to fail on non-struct type, but it did not")
	}
}

// ---- AssertAllRequiredFields ----

func TestAssertAllRequiredFields_AllPresent(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(sampleResource{}), []string{
		"ResourceID", "ResourceName", "Status",
	})
}

func TestAssertAllRequiredFields_MissingField(t *testing.T) {
	fake := &testing.T{}
	testutil.AssertAllRequiredFields(fake, reflect.TypeOf(sampleResource{}), []string{
		"ResourceID", "DoesNotExist",
	})
	if !fake.Failed() {
		t.Error("expected fake test to fail on missing field, but it did not")
	}
}

func TestAssertAllRequiredFields_EmptyList(t *testing.T) {
	// Empty required list should always pass.
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(sampleResource{}), nil)
}

// ---- AssertJSONRoundTrip ----

func TestAssertJSONRoundTrip_Pass(t *testing.T) {
	obj := sampleResource{
		ResourceID:   "res-001",
		ResourceName: "my-resource",
	}
	// Status is omitempty so it should be absent when empty.
	testutil.AssertJSONRoundTrip(t, obj, `{"resourceId":"res-001","resourceName":"my-resource"}`)
}

func TestAssertJSONRoundTrip_WithOptionalField(t *testing.T) {
	obj := sampleResource{
		ResourceID:   "res-002",
		ResourceName: "another",
		Status:       "ACTIVE",
	}
	testutil.AssertJSONRoundTrip(t, obj, `{"resourceId":"res-002","resourceName":"another","status":"ACTIVE"}`)
}

func TestAssertJSONRoundTrip_WhitespaceNormalized(t *testing.T) {
	obj := sampleResource{ResourceID: "r1", ResourceName: "n1"}
	// Extra whitespace in expectedJSON should not cause failure.
	testutil.AssertJSONRoundTrip(t, obj, `{
		"resourceId":   "r1",
		"resourceName": "n1"
	}`)
}

func TestAssertJSONRoundTrip_Mismatch(t *testing.T) {
	fake := &testing.T{}
	obj := sampleResource{ResourceID: "correct-id", ResourceName: "name"}
	testutil.AssertJSONRoundTrip(fake, obj, `{"resourceId":"wrong-id","resourceName":"name"}`)
	if !fake.Failed() {
		t.Error("expected fake test to fail on JSON mismatch, but it did not")
	}
}

// ---- NewTestHTTPServer ----

func TestNewTestHTTPServer_DefaultHandler(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/items") //nolint:noctx
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
	if recorded.Path != "/api/v1/items" {
		t.Errorf("path: got %q, want %q", recorded.Path, "/api/v1/items")
	}
	if recorded.Method != http.MethodGet {
		t.Errorf("method: got %q, want %q", recorded.Method, http.MethodGet)
	}
}

func TestNewTestHTTPServer_CustomHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"42"}`))
	})

	server, _ := testutil.NewTestHTTPServer(handler)
	defer server.Close()

	resp, err := http.Post(server.URL+"/resources", "application/json", strings.NewReader(`{"name":"x"}`)) //nolint:noctx
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status: got %d, want 201", resp.StatusCode)
	}
}

func TestNewTestHTTPServer_RecordsBody(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	payload := `{"key":"value"}`
	resp, err := http.Post(server.URL+"/submit", "application/json", strings.NewReader(payload)) //nolint:noctx
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if string(recorded.Body) != payload {
		t.Errorf("body: got %q, want %q", string(recorded.Body), payload)
	}
}

// ---- AssertRequestMethod ----

func TestAssertRequestMethod_Pass(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	http.Get(server.URL + "/x") //nolint:errcheck,noctx
	testutil.AssertRequestMethod(t, recorded, "GET")
}

func TestAssertRequestMethod_Fail(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	http.Get(server.URL + "/x") //nolint:errcheck,noctx

	fake := &testing.T{}
	testutil.AssertRequestMethod(fake, recorded, "POST")
	if !fake.Failed() {
		t.Error("expected fake test to fail on method mismatch, but it did not")
	}
}

// ---- AssertRequestPath ----

func TestAssertRequestPath_Pass(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	http.Get(server.URL + "/v2/instances") //nolint:errcheck,noctx
	testutil.AssertRequestPath(t, recorded, "/v2/instances")
}

func TestAssertRequestPath_Fail(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	http.Get(server.URL + "/v2/instances") //nolint:errcheck,noctx

	fake := &testing.T{}
	testutil.AssertRequestPath(fake, recorded, "/v1/instances")
	if !fake.Failed() {
		t.Error("expected fake test to fail on path mismatch, but it did not")
	}
}

// ---- AssertRequestHasHeader ----

func TestAssertRequestHasHeader_Pass(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/x", nil) //nolint:noctx
	req.Header.Set("X-Auth-Token", "secret-token")
	http.DefaultClient.Do(req) //nolint:errcheck

	testutil.AssertRequestHasHeader(t, recorded, "X-Auth-Token", "secret-token")
}

func TestAssertRequestHasHeader_Fail(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/x", nil) //nolint:noctx
	req.Header.Set("X-Auth-Token", "secret-token")
	http.DefaultClient.Do(req) //nolint:errcheck

	fake := &testing.T{}
	testutil.AssertRequestHasHeader(fake, recorded, "X-Auth-Token", "wrong-value")
	if !fake.Failed() {
		t.Error("expected fake test to fail on header mismatch, but it did not")
	}
}

// ---- RecordedRequest ----

func TestRecordedRequest_CapturesAllFields(t *testing.T) {
	server, recorded := testutil.NewTestHTTPServer(nil)
	defer server.Close()

	body := `{"hello":"world"}`
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/resource?page=2", strings.NewReader(body)) //nolint:noctx
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer tok")
	http.DefaultClient.Do(req) //nolint:errcheck

	if recorded.Method != http.MethodPost {
		t.Errorf("Method: got %q, want %q", recorded.Method, http.MethodPost)
	}
	if recorded.Path != "/api/resource" {
		t.Errorf("Path: got %q, want %q", recorded.Path, "/api/resource")
	}
	if recorded.Query != "page=2" {
		t.Errorf("Query: got %q, want %q", recorded.Query, "page=2")
	}
	if string(recorded.Body) != body {
		t.Errorf("Body: got %q, want %q", string(recorded.Body), body)
	}
	if recorded.Headers.Get("Authorization") != "Bearer tok" {
		t.Errorf("Authorization header not captured")
	}

	// Verify body is valid JSON (sanity check).
	var parsed map[string]string
	if err := json.Unmarshal(recorded.Body, &parsed); err != nil {
		t.Errorf("captured body is not valid JSON: %v", err)
	}
}
