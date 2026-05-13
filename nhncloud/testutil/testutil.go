// Package testutil provides test utilities for API spec compliance testing.
//
// Provides helpers for:
//  1. Checking struct fields match API documentation
//  2. Verifying JSON tags on struct fields
//  3. Testing HTTP request building (method, URL, headers, body)
//  4. Testing response parsing from sample JSON
package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// RecordedRequest captures an HTTP request for assertion in tests.
type RecordedRequest struct {
	Method  string
	Path    string
	Headers http.Header
	Body    []byte
	Query   string
}

// NewTestHTTPServer creates a test HTTP server that records the most recent
// request for later assertion. The handler parameter is optional; when nil
// the server responds with 200 OK and an empty JSON object "{}".
// The caller is responsible for closing the server with server.Close().
//
// Usage:
//
//	var recorded *testutil.RecordedRequest
//	server := testutil.NewTestHTTPServer(nil)
//	defer server.Close()
//	// ... make request to server.URL ...
//	// recorded is updated after each request
func NewTestHTTPServer(handler http.HandlerFunc) (server *httptest.Server, recorded *RecordedRequest) {
	rec := &RecordedRequest{}

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request.
		body, _ := io.ReadAll(r.Body)
		rec.Method = r.Method
		rec.Path = r.URL.Path
		rec.Headers = r.Header.Clone()
		rec.Body = body
		rec.Query = r.URL.RawQuery

		if handler != nil {
			// Restore body so the user handler can also read it.
			r.Body = io.NopCloser(strings.NewReader(string(body)))
			handler(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	})

	srv := httptest.NewServer(mux)
	return srv, rec
}

// AssertRequestMethod asserts that the recorded HTTP request used the expected
// method (e.g. "GET", "POST").
func AssertRequestMethod(t *testing.T, recorded *RecordedRequest, expected string) {
	t.Helper()
	if recorded.Method != expected {
		t.Errorf("HTTP method: got %q, want %q", recorded.Method, expected)
	}
}

// AssertRequestPath asserts that the recorded HTTP request targeted the
// expected URL path (e.g. "/v2.0/appkeys/key/db-instances").
func AssertRequestPath(t *testing.T, recorded *RecordedRequest, expected string) {
	t.Helper()
	if recorded.Path != expected {
		t.Errorf("request path: got %q, want %q", recorded.Path, expected)
	}
}

// AssertRequestHasHeader asserts that the recorded HTTP request contained a
// header with the given key and value.
func AssertRequestHasHeader(t *testing.T, recorded *RecordedRequest, key, value string) {
	t.Helper()
	got := recorded.Headers.Get(key)
	if got != value {
		t.Errorf("header %q: got %q, want %q", key, got, value)
	}
}

// AssertStructHasJSONTag verifies that a field on structType carries the
// expected JSON struct tag. structType must be a reflect.Type of a struct.
//
// Usage:
//
//	AssertStructHasJSONTag(t, reflect.TypeOf(mysql.DBInstance{}), "DBInstanceID", "dbInstanceId")
func AssertStructHasJSONTag(t *testing.T, structType reflect.Type, fieldName, expectedTag string) {
	t.Helper()

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		t.Errorf("AssertStructHasJSONTag: %v is not a struct type", structType)
		return
	}

	field, ok := structType.FieldByName(fieldName)
	if !ok {
		t.Errorf("struct %v has no field %q", structType.Name(), fieldName)
		return
	}

	tag := field.Tag.Get("json")
	// Strip options like ",omitempty".
	tagName := strings.Split(tag, ",")[0]

	if tagName != expectedTag {
		t.Errorf("field %v.%s json tag: got %q, want %q", structType.Name(), fieldName, tagName, expectedTag)
	}
}

// AssertAllRequiredFields verifies that every field name in requiredFields
// exists on structType. This ensures API-required fields are not accidentally
// omitted from the Go struct definition.
//
// Usage:
//
//	AssertAllRequiredFields(t, reflect.TypeOf(mysql.DBInstance{}), []string{
//	    "DBInstanceID", "DBInstanceName", "DBInstanceStatus",
//	})
func AssertAllRequiredFields(t *testing.T, structType reflect.Type, requiredFields []string) {
	t.Helper()

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		t.Errorf("AssertAllRequiredFields: %v is not a struct type", structType)
		return
	}

	for _, fieldName := range requiredFields {
		if _, ok := structType.FieldByName(fieldName); !ok {
			t.Errorf("struct %v is missing required field %q", structType.Name(), fieldName)
		}
	}
}

// AssertJSONRoundTrip verifies that obj marshals to expectedJSON and that
// unmarshaling expectedJSON back into the same type produces an equal value.
//
// expectedJSON may contain extra whitespace; both sides are compared after
// re-marshaling to canonical form so field ordering does not matter.
//
// Usage:
//
//	instance := mysql.DBInstance{DBInstanceID: "id-1", DBInstanceName: "my-db"}
//	AssertJSONRoundTrip(t, instance, `{"dbInstanceId":"id-1","dbInstanceName":"my-db"}`)
func AssertJSONRoundTrip(t *testing.T, obj interface{}, expectedJSON string) {
	t.Helper()

	// Marshal obj → JSON.
	got, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Normalize both sides to canonical JSON for comparison.
	var gotNorm, wantNorm interface{}
	if err := json.Unmarshal(got, &gotNorm); err != nil {
		t.Fatalf("failed to normalize marshaled JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedJSON), &wantNorm); err != nil {
		t.Fatalf("expectedJSON is not valid JSON: %v", err)
	}

	gotCanon, _ := json.Marshal(gotNorm)
	wantCanon, _ := json.Marshal(wantNorm)

	if string(gotCanon) != string(wantCanon) {
		t.Errorf("JSON mismatch:\n  got:  %s\n  want: %s", gotCanon, wantCanon)
	}

	// Unmarshal expectedJSON back into the same type and verify equality.
	typ := reflect.TypeOf(obj)
	ptr := reflect.New(typ)
	if err := json.Unmarshal([]byte(expectedJSON), ptr.Interface()); err != nil {
		t.Fatalf("json.Unmarshal into %v failed: %v", typ, err)
	}
	roundTripped := ptr.Elem().Interface()

	if !reflect.DeepEqual(obj, roundTripped) {
		t.Errorf("round-trip value mismatch:\n  original:    %+v\n  round-tripped: %+v", obj, roundTripped)
	}
}
