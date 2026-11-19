package ncr

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	iclient "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type staticToken struct{}

func (staticToken) GetToken(context.Context) (string, error) { return "t", nil }
func (staticToken) SetAuthHeader(*http.Request, string)      {}

// Webhooks live at /webhook/policies and tags under an artifact; /webhooks
// and /images/{image}/tags answer 404.
func TestClientRoutes(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.Path)
		io.WriteString(w, `{"header":{"isSuccessful":true,"resultCode":0,"resultMessage":"SUCCESS"}}`)
	}))
	defer srv.Close()
	c := &Client{httpClient: iclient.NewClient(srv.URL, staticToken{})}
	ctx := context.Background()

	c.ListWebhooks(ctx, "reg")
	c.CreateWebhook(ctx, "reg", &CreateWebhookInput{})
	c.DeleteWebhook(ctx, "reg", "7")
	c.ListArtifacts(ctx, "reg", "img")
	c.DeleteTag(ctx, "reg", "img", "v1")

	want := []string{
		"GET /registries/reg/webhook/policies",
		"POST /registries/reg/webhook/policies",
		"DELETE /registries/reg/webhook/policies/7",
		"GET /registries/reg/images/img/artifacts",
		"DELETE /registries/reg/images/img/artifacts/v1/tags/v1",
	}
	if len(got) != len(want) {
		t.Fatalf("requests %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("request %d: %s, want %s", i, got[i], want[i])
		}
	}
}
