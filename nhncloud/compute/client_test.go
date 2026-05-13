package compute

// client_test.go lives in package compute (white-box) so it can call the
// unexported newClientWithHTTPClient helper defined below (also in this
// package) without needing to export anything from the production code.
//
// Strategy: each test spins up two httptest servers:
//   1. identitySrv – returns a fake OpenStack token + service catalog entry
//      pointing at the second server.
//   2. computeSrv  – the actual compute endpoint; requests are captured via
//      testutil.NewTestHTTPServer for assertion.
//
// The compute.Client is constructed via newTestClient() which bypasses the
// real identity service entirely.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	iclient "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// staticTokenProvider satisfies client.TokenProvider without any network I/O.
type staticTokenProvider struct {
	token string
}

func (s *staticTokenProvider) GetToken(_ context.Context) (string, error) {
	return s.token, nil
}

func (s *staticTokenProvider) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("X-Auth-Token", token)
}

// newTestClient returns a compute.Client whose httpClient is pre-wired to
// baseURL, bypassing the real identity flow.  The token is fixed to
// "test-token".
func newTestClient(baseURL string) *Client {
	tp := &staticTokenProvider{token: "test-token"}
	c := &Client{}
	c.httpClient = iclient.NewClient(baseURL, tp)
	return c
}

// fakeIdentityServer returns an httptest.Server that mimics the OpenStack
// identity v2.0 token endpoint.  It issues a token that expires in 1 hour
// and places computeEndpoint into the service catalog under type "compute".
func fakeIdentityServer(t *testing.T, computeEndpoint string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2.0/tokens" {
			http.NotFound(w, r)
			return
		}
		resp := map[string]interface{}{
			"access": map[string]interface{}{
				"token": map[string]interface{}{
					"id":      "test-token",
					"expires": time.Now().Add(time.Hour).Format(time.RFC3339),
				},
				"serviceCatalog": []map[string]interface{}{
					{
						"name": "nova",
						"type": "compute",
						"endpoints": []map[string]string{
							{"publicURL": computeEndpoint, "region": "kr1"},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
}

// ---------------------------------------------------------------------------
// ListServers
// ---------------------------------------------------------------------------

func TestListServers_RequestMethod(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListServersOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, err := c.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers returned error: %v", err)
	}

	testutil.AssertRequestMethod(t, recorded, http.MethodGet)
}

func TestListServers_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListServersOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.ListServers(context.Background())

	// The SDK calls /servers/detail (detailed list endpoint per API spec).
	testutil.AssertRequestPath(t, recorded, "/servers/detail")
}

func TestListServers_AuthHeader(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListServersOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.ListServers(context.Background())

	testutil.AssertRequestHasHeader(t, recorded, "X-Auth-Token", "test-token")
}

func TestListServers_ParseResponse(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"servers":[{"id":"s1","name":"web","status":"ACTIVE","tenant_id":"t","user_id":"u","created":"2026-01-01T00:00:00Z"}]}`)) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	out, err := c.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers: %v", err)
	}

	if len(out.Servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(out.Servers))
	}
	if out.Servers[0].ID != "s1" {
		t.Errorf("Server.ID: got %q, want %q", out.Servers[0].ID, "s1")
	}
	if out.Servers[0].Status != "ACTIVE" {
		t.Errorf("Server.Status: got %q, want %q", out.Servers[0].Status, "ACTIVE")
	}
}

// ---------------------------------------------------------------------------
// GetServer
// ---------------------------------------------------------------------------

func TestGetServer_RequestMethod(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.GetServer(context.Background(), "srv-abc")

	testutil.AssertRequestMethod(t, recorded, http.MethodGet)
}

func TestGetServer_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.GetServer(context.Background(), "srv-abc")

	testutil.AssertRequestPath(t, recorded, "/servers/srv-abc")
}

// ---------------------------------------------------------------------------
// CreateServer
// ---------------------------------------------------------------------------

func TestCreateServer_RequestMethod(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.CreateServer(context.Background(), &CreateServerInput{
		Name:      "web-01",
		FlavorRef: "m2.c2m4",
	})

	testutil.AssertRequestMethod(t, recorded, http.MethodPost)
}

func TestCreateServer_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.CreateServer(context.Background(), &CreateServerInput{Name: "web-01", FlavorRef: "m2.c2m4"})

	testutil.AssertRequestPath(t, recorded, "/servers")
}

func TestCreateServer_RequestBodyWrapped(t *testing.T) {
	// testutil.NewTestHTTPServer captures the raw request body in recorded.Body
	// before the handler runs, so we can inspect it directly.
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	input := &CreateServerInput{
		Name:      "web-01",
		FlavorRef: "m2.c2m4",
		ImageRef:  "img-001",
		KeyName:   "my-key",
		Networks:  []ServerNetwork{{UUID: "net-abc"}},
		BlockDeviceMapping: []BlockDeviceMapping{
			{BootIndex: 0, SourceType: "image", DestinationType: "volume", VolumeSize: 50, DeleteOnTermination: true},
		},
	}
	_, _ = c.CreateServer(context.Background(), input)

	// The request body must have a top-level "server" key.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(recorded.Body, &top); err != nil {
		t.Fatalf("request body is not valid JSON: %v (body=%s)", err, recorded.Body)
	}
	if _, ok := top["server"]; !ok {
		t.Errorf("expected top-level 'server' key in request body; got keys: %v", jsonKeys(top))
	}

	// Inner fields must match API spec field names.
	var inner map[string]json.RawMessage
	if err := json.Unmarshal(top["server"], &inner); err != nil {
		t.Fatalf("inner server body is not valid JSON: %v", err)
	}
	for _, key := range []string{"name", "flavorRef", "imageRef", "key_name", "networks", "block_device_mapping_v2"} {
		if _, ok := inner[key]; !ok {
			t.Errorf("expected field %q in server body; got keys: %v", key, jsonKeys(inner))
		}
	}
}

func TestCreateServer_ContentTypeHeader(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateServerOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.CreateServer(context.Background(), &CreateServerInput{Name: "x", FlavorRef: "f"})

	ct := recorded.Headers.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: got %q, want application/json", ct)
	}
}

// ---------------------------------------------------------------------------
// DeleteServer
// ---------------------------------------------------------------------------

func TestDeleteServer_RequestMethod(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.DeleteServer(context.Background(), "srv-del")

	testutil.AssertRequestMethod(t, recorded, http.MethodDelete)
}

func TestDeleteServer_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.DeleteServer(context.Background(), "srv-del")

	testutil.AssertRequestPath(t, recorded, "/servers/srv-del")
}

func TestDeleteServer_ReturnsNilOn204(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	if err := c.DeleteServer(context.Background(), "srv-del"); err != nil {
		t.Errorf("expected nil error on 204, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// StartServer / StopServer actions
// ---------------------------------------------------------------------------

func TestStartServer_ActionBody(t *testing.T) {
	var capturedBody []byte
	computeSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		capturedBody = buf[:n]
		w.WriteHeader(http.StatusAccepted)
	}))
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.StartServer(context.Background(), "srv-xyz")

	var body map[string]interface{}
	if err := json.Unmarshal(capturedBody, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if _, ok := body["os-start"]; !ok {
		t.Errorf("expected 'os-start' key in action body; got: %v", body)
	}
}

func TestStopServer_ActionBody(t *testing.T) {
	var capturedBody []byte
	computeSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		capturedBody = buf[:n]
		w.WriteHeader(http.StatusAccepted)
	}))
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.StopServer(context.Background(), "srv-xyz")

	var body map[string]interface{}
	if err := json.Unmarshal(capturedBody, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if _, ok := body["os-stop"]; !ok {
		t.Errorf("expected 'os-stop' key in action body; got: %v", body)
	}
}

func TestStartServer_ActionPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.StartServer(context.Background(), "srv-xyz")

	testutil.AssertRequestPath(t, recorded, "/servers/srv-xyz/action")
	testutil.AssertRequestMethod(t, recorded, http.MethodPost)
}

func TestStopServer_ActionPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.StopServer(context.Background(), "srv-xyz")

	testutil.AssertRequestPath(t, recorded, "/servers/srv-xyz/action")
	testutil.AssertRequestMethod(t, recorded, http.MethodPost)
}

// ---------------------------------------------------------------------------
// RebootServer
// ---------------------------------------------------------------------------

func TestRebootServer_SoftBody(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.RebootServer(context.Background(), "srv-r", false)

	var body map[string]map[string]string
	if err := json.Unmarshal(recorded.Body, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if body["reboot"]["type"] != "SOFT" {
		t.Errorf("reboot.type: got %q, want SOFT", body["reboot"]["type"])
	}
}

func TestRebootServer_HardBody(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.RebootServer(context.Background(), "srv-r", true)

	var body map[string]map[string]string
	if err := json.Unmarshal(recorded.Body, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if body["reboot"]["type"] != "HARD" {
		t.Errorf("reboot.type: got %q, want HARD", body["reboot"]["type"])
	}
}

// ---------------------------------------------------------------------------
// ResizeServer
// ---------------------------------------------------------------------------

func TestResizeServer_Body(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.ResizeServer(context.Background(), "srv-r", "m2.c4m8")

	var body map[string]map[string]string
	if err := json.Unmarshal(recorded.Body, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if body["resize"]["flavorRef"] != "m2.c4m8" {
		t.Errorf("resize.flavorRef: got %q, want m2.c4m8", body["resize"]["flavorRef"])
	}
}

func TestResizeServer_Path(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.ResizeServer(context.Background(), "srv-r", "m2.c4m8")

	testutil.AssertRequestPath(t, recorded, "/servers/srv-r/action")
}

// ---------------------------------------------------------------------------
// ConfirmResize
// ---------------------------------------------------------------------------

func TestConfirmResize_Body(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.ConfirmResize(context.Background(), "srv-r")

	var body map[string]interface{}
	if err := json.Unmarshal(recorded.Body, &body); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if _, ok := body["confirmResize"]; !ok {
		t.Errorf("expected 'confirmResize' key; got: %v", body)
	}
}

// ---------------------------------------------------------------------------
// ListFlavors
// ---------------------------------------------------------------------------

func TestListFlavors_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListFlavorsOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.ListFlavors(context.Background())

	testutil.AssertRequestPath(t, recorded, "/flavors/detail")
	testutil.AssertRequestMethod(t, recorded, http.MethodGet)
}

func TestListFlavors_ParseResponse(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"flavors":[{"id":"m2.c1m2","name":"m2.c1m2","ram":2048,"vcpus":1,"disk":20}]}`)) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	out, err := c.ListFlavors(context.Background())
	if err != nil {
		t.Fatalf("ListFlavors: %v", err)
	}
	if len(out.Flavors) != 1 {
		t.Fatalf("expected 1 flavor, got %d", len(out.Flavors))
	}
	if out.Flavors[0].RAM != 2048 {
		t.Errorf("RAM: got %d, want 2048", out.Flavors[0].RAM)
	}
}

// ---------------------------------------------------------------------------
// ListImages
// ---------------------------------------------------------------------------

func TestListImages_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListImagesOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.ListImages(context.Background())

	testutil.AssertRequestPath(t, recorded, "/images/detail")
	testutil.AssertRequestMethod(t, recorded, http.MethodGet)
}

// ---------------------------------------------------------------------------
// ListKeyPairs
// ---------------------------------------------------------------------------

func TestListKeyPairs_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListKeyPairsOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.ListKeyPairs(context.Background())

	testutil.AssertRequestPath(t, recorded, "/os-keypairs")
	testutil.AssertRequestMethod(t, recorded, http.MethodGet)
}

func TestListKeyPairs_ParseResponse(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keypairs":[{"keypair":{"name":"k1","public_key":"ssh-rsa AAA","fingerprint":"aa:bb"}}]}`)) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	out, err := c.ListKeyPairs(context.Background())
	if err != nil {
		t.Fatalf("ListKeyPairs: %v", err)
	}
	if len(out.KeyPairs) != 1 {
		t.Fatalf("expected 1 keypair, got %d", len(out.KeyPairs))
	}
	if out.KeyPairs[0].KeyPair.Name != "k1" {
		t.Errorf("KeyPair.Name: got %q, want k1", out.KeyPairs[0].KeyPair.Name)
	}
}

// ---------------------------------------------------------------------------
// CreateKeyPair
// ---------------------------------------------------------------------------

func TestCreateKeyPair_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateKeyPairOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.CreateKeyPair(context.Background(), &CreateKeyPairInput{Name: "new-key"})

	testutil.AssertRequestPath(t, recorded, "/os-keypairs")
	testutil.AssertRequestMethod(t, recorded, http.MethodPost)
}

func TestCreateKeyPair_RequestBodyWrapped(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateKeyPairOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, _ = c.CreateKeyPair(context.Background(), &CreateKeyPairInput{
		Name:      "new-key",
		PublicKey: "ssh-rsa AAAA...",
	})

	var top map[string]json.RawMessage
	if err := json.Unmarshal(recorded.Body, &top); err != nil {
		t.Fatalf("body not valid JSON: %v (body=%s)", err, recorded.Body)
	}
	if _, ok := top["keypair"]; !ok {
		t.Errorf("expected top-level 'keypair' key; got: %v", jsonKeys(top))
	}

	var inner map[string]string
	if err := json.Unmarshal(top["keypair"], &inner); err != nil {
		t.Fatalf("inner keypair unmarshal failed: %v", err)
	}
	if inner["name"] != "new-key" {
		t.Errorf("keypair.name: got %q, want new-key", inner["name"])
	}
	if inner["public_key"] != "ssh-rsa AAAA..." {
		t.Errorf("keypair.public_key: got %q, want ssh-rsa AAAA...", inner["public_key"])
	}
}

// ---------------------------------------------------------------------------
// DeleteKeyPair
// ---------------------------------------------------------------------------

func TestDeleteKeyPair_RequestPath(t *testing.T) {
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_ = c.DeleteKeyPair(context.Background(), "my-key")

	testutil.AssertRequestPath(t, recorded, "/os-keypairs/my-key")
	testutil.AssertRequestMethod(t, recorded, http.MethodDelete)
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

func TestListServers_APIError(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized","error_code":"AUTH_FAILED"}`)) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, err := c.ListServers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 in error; got: %v", err)
	}
}

func TestGetServer_APIError404(t *testing.T) {
	computeSrv, _ := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Instance not found"}`)) //nolint:errcheck
	})
	defer computeSrv.Close()

	c := newTestClient(computeSrv.URL)
	_, err := c.GetServer(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected 404 in error; got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ErrNoCredentials
// ---------------------------------------------------------------------------

func TestListServers_ErrNoCredentials(t *testing.T) {
	// A Client with no tokenProvider and no httpClient returns ErrNoCredentials.
	c := &Client{}
	_, err := c.ListServers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrNoCredentials {
		t.Errorf("expected ErrNoCredentials, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Bootstrap via fake identity server (integration of NewClient + ensureClient)
// ---------------------------------------------------------------------------

func TestNewClient_BootstrapViaIdentity(t *testing.T) {
	// Spin up the compute endpoint first so we can embed its URL in the catalog.
	computeSrv, recorded := testutil.NewTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListServersOutput{}) //nolint:errcheck
	})
	defer computeSrv.Close()

	identitySrv := fakeIdentityServer(t, computeSrv.URL)
	defer identitySrv.Close()

	// Build a real IdentityTokenProvider pointed at the fake identity server.
	tp := iclient.NewIdentityTokenProviderWithURL(
		identitySrv.URL, "ten-abc", "user", "pass",
	)

	c := &Client{
		region:        "kr1",
		tokenProvider: tp,
	}

	_, err := c.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers via fake identity: %v", err)
	}

	// Verify that the request landed on the compute server.
	testutil.AssertRequestPath(t, recorded, "/servers/detail")
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func jsonKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
