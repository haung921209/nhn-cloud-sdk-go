package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const debugSecret = "s3cr3t-value-do-not-print"

// captureStdio runs fn with os.Stdout and os.Stderr redirected to pipes.
func captureStdio(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	outR, outW, _ := os.Pipe()
	errR, errW, _ := os.Pipe()
	os.Stdout, os.Stderr = outW, errW
	outC, errC := make(chan string), make(chan string)
	go func() { b, _ := io.ReadAll(outR); outC <- string(b) }()
	go func() { b, _ := io.ReadAll(errR); errC <- string(b) }()
	defer func() { os.Stdout, os.Stderr = oldOut, oldErr }()
	fn()
	outW.Close()
	errW.Close()
	return <-outC, <-errC
}

func TestDebugOutputRedactedAndOnStderr(t *testing.T) {
	t.Setenv("NHN_SDK_CAPTURE_DIR", "")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"header":{"isSuccessful":true},"dbUserPassword":"` + debugSecret + `","dbInstanceName":"db-01"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithDebug(true), WithoutRetry(),
		WithAppKeyAuth("example-app-key", "example-access-id", debugSecret))
	stdout, stderr := captureStdio(t, func() {
		body := map[string]string{"dbInstanceName": "db-01", "dbPassword": debugSecret}
		if err := c.POST(context.Background(), "/v3.0/db-instances", body, nil); err != nil {
			t.Errorf("POST: %v", err)
		}
	})

	if stdout != "" {
		t.Errorf("debug output must not go to stdout, got:\n%s", stdout)
	}
	if strings.Contains(stderr, debugSecret) {
		t.Errorf("secret leaked to stderr:\n%s", stderr)
	}
	for _, want := range []string{"POST", "/v3.0/db-instances", "200 OK", "X-Tc-Authentication-Secret: ***", `"db-01"`} {
		if !strings.Contains(stderr, want) {
			t.Errorf("debug output missing %q:\n%s", want, stderr)
		}
	}
}
