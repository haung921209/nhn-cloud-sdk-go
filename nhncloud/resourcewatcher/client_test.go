package resourcewatcher

import (
	"context"
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
	c := NewClient("example-app-key", "id", "secret", nil, false)
	c.baseURL = srv.URL
	_, err := c.ListEvents(context.Background())
	if err == nil || !strings.Contains(err.Error(), "4010001") || !strings.Contains(err.Error(), "Invalid appKey.") {
		t.Fatalf("want an error with the result code and message, got %v", err)
	}
}
