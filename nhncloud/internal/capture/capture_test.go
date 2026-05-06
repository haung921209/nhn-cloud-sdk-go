package capture

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransportSavesResponseWhenEnvSet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NHN_SDK_CAPTURE_DIR", dir)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := &http.Client{Transport: NewTransport(http.DefaultTransport)}
	resp, err := c.Get(srv.URL + "/v4.0/db-instances")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != `{"ok":true}` {
		t.Fatalf("body wasn't passed through; got %q", body)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "GET_v4.0_db-instances_*.json"))
	if len(matches) != 1 {
		t.Fatalf("expected one capture file, got %v", matches)
	}
	data, _ := os.ReadFile(matches[0])
	if !strings.Contains(string(data), `"ok":true`) {
		t.Fatalf("body missing in capture file: %s", data)
	}
}

func TestTransportNoopWhenEnvUnset(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NHN_SDK_CAPTURE_DIR", "")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`hello`))
	}))
	defer srv.Close()

	c := &http.Client{Transport: NewTransport(http.DefaultTransport)}
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("expected no capture, got %d files", len(entries))
	}
}

func TestWrapTransportNilInner(t *testing.T) {
	// Wrapping nil should yield a Transport that uses http.DefaultTransport.
	rt := WrapTransport(nil)
	if rt == nil {
		t.Fatal("WrapTransport(nil) returned nil")
	}
}

func TestTransportPassthroughWithNilBody(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NHN_SDK_CAPTURE_DIR", dir)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := &http.Client{Transport: NewTransport(http.DefaultTransport)}
	resp, err := c.Get(srv.URL + "/empty")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	resp.Body.Close()
}
