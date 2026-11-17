package client

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

type headerTokenProvider struct{}

func (headerTokenProvider) GetToken(ctx context.Context) (string, error) { return debugSecret, nil }
func (headerTokenProvider) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("X-Auth-Token", token)
}

func TestDebugOutputRedactedAndOnStderr(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Subject-Token", debugSecret)
		_, _ = w.Write([]byte(`{"instance":{"name":"db-01","adminPassword":"` + debugSecret + `"}}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, headerTokenProvider{}, WithDebug(true))
	stdout, stderr := captureStdio(t, func() {
		body := map[string]string{"name": "db-01", "password": debugSecret}
		if err := c.POST(context.Background(), "/v1/instances", body, nil); err != nil {
			t.Errorf("POST: %v", err)
		}
	})

	if stdout != "" {
		t.Errorf("debug output must not go to stdout, got:\n%s", stdout)
	}
	if strings.Contains(stderr, debugSecret) {
		t.Errorf("secret leaked to stderr:\n%s", stderr)
	}
	for _, want := range []string{"POST", "/v1/instances", "200 OK", "X-Auth-Token: ***", `"db-01"`} {
		if !strings.Contains(stderr, want) {
			t.Errorf("debug output missing %q:\n%s", want, stderr)
		}
	}
}
