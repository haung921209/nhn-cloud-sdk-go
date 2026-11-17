package keymanager

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

func TestDebugOutputHidesSecretPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"header":{"isSuccessful":true},"body":{"secret":"` + debugSecret + `","payload":"` + debugSecret + `"}}`))
	}))
	defer srv.Close()

	c := NewClient("KR1", "example-app-key", "example-access-id", debugSecret, true)
	c.baseURL = srv.URL
	stdout, stderr := captureStdio(t, func() {
		if err := c.request(context.Background(), "GET", "/keymanager/v1.3/appkey/example-app-key/secrets/abc", nil, nil); err != nil {
			t.Errorf("request: %v", err)
		}
	})

	if stdout != "" {
		t.Errorf("debug output must not go to stdout, got:\n%s", stdout)
	}
	if strings.Contains(stderr, debugSecret) {
		t.Errorf("secret leaked to stderr:\n%s", stderr)
	}
	if !strings.Contains(stderr, "200 OK") || !strings.Contains(stderr, "isSuccessful") {
		t.Errorf("debug output lost status or structure:\n%s", stderr)
	}
}
