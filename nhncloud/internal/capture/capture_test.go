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

func TestTransportCapturePermissionsAndRedaction(t *testing.T) {
	const secret = "s3cr3t-value-do-not-print"
	dir := filepath.Join(t.TempDir(), "capture")
	t.Setenv("NHN_SDK_CAPTURE_DIR", dir)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pem" {
			_, _ = w.Write([]byte("-----BEGIN PRIVATE KEY-----\n" + secret))
			return
		}
		_, _ = w.Write([]byte(`{"name":"db-01","dbPassword":"` + secret + `"}`))
	}))
	defer srv.Close()

	c := &http.Client{Transport: NewTransport(http.DefaultTransport)}
	for _, p := range []string{"/v1/instances", "/pem"} {
		resp, err := c.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if !strings.Contains(string(body), secret) {
			t.Fatalf("caller must still see the unredacted body, got %q", body)
		}
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("capture dir mode = %o, want 700", perm)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected only the JSON response to be captured, got %d files", len(entries))
	}
	path := filepath.Join(dir, entries[0].Name())
	fi, _ := os.Stat(path)
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("capture file mode = %o, want 600", perm)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), secret) {
		t.Errorf("capture file contains secret: %s", data)
	}
	if !strings.Contains(string(data), `"name":"db-01"`) {
		t.Errorf("capture file lost non-sensitive fields: %s", data)
	}
}

// An existing world-readable capture directory is tightened, and an appkey in
// the request path does not appear in the capture file name.
func TestTransportCaptureExistingDirAndAppKeyPath(t *testing.T) {
	const appkey = "AppKeyInPath123"
	dir := filepath.Join(t.TempDir(), "capture")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("NHN_SDK_CAPTURE_DIR", dir)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	c := &http.Client{Transport: NewTransport(http.DefaultTransport)}
	resp, err := c.Get(srv.URL + "/keymanager/v1.3/appkey/" + appkey + "/secrets/abc")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if info, _ := os.Stat(dir); info.Mode().Perm() != 0o700 {
		t.Errorf("existing capture dir mode = %o, want 700", info.Mode().Perm())
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("want 1 capture file, got %d", len(entries))
	}
	if strings.Contains(entries[0].Name(), appkey) {
		t.Errorf("capture file name leaks the appkey: %s", entries[0].Name())
	}
}
