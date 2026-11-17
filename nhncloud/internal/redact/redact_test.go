package redact

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const secret = "s3cr3t-value-do-not-print"

func TestHeadersMaskSensitiveNames(t *testing.T) {
	sensitive := []string{
		"X-Auth-Token", "X-TC-AUTHENTICATION-SECRET", "X-TC-AUTHENTICATION-ID",
		"X-TC-APP-KEY", "Authorization", "X-NHN-AUTHORIZATION", "Cookie", "Set-Cookie",
		"X-Custom-Secret", "X-Subject-Token", "X-Db-Password", "X-Api-Key",
	}
	h := http.Header{}
	for _, n := range sensitive {
		h[n] = []string{secret}
	}
	h.Set("Content-Type", "application/json")
	out := Headers(h)
	if strings.Contains(out, secret) {
		t.Fatalf("sensitive header value leaked:\n%s", out)
	}
	for _, n := range sensitive {
		if !strings.Contains(out, n+": "+Mask) {
			t.Errorf("header %s not masked:\n%s", n, out)
		}
	}
	if !strings.Contains(out, "Content-Type: application/json") {
		t.Errorf("non-sensitive header should be kept:\n%s", out)
	}
}

func TestBodyMasksSensitiveKeys(t *testing.T) {
	keys := []string{
		"password", "dbPassword", "passwd", "clientSecret", "token", "accessToken",
		"privateKey", "private_key", "accessKey", "access_key", "secretKey",
		"credentials", "kubeconfig", "payload",
	}
	for _, k := range keys {
		out := Body([]byte(`{"` + k + `":"` + secret + `","name":"db-01"}`))
		if strings.Contains(out, secret) {
			t.Errorf("key %q leaked: %s", k, out)
		}
		if !strings.Contains(out, `"db-01"`) {
			t.Errorf("key %q: non-sensitive value dropped: %s", k, out)
		}
	}
}

func TestBodyMasksNestedJSON(t *testing.T) {
	in := `{"access":{"token":{"id":"` + secret + `"},"user":{"name":"alice"}},` +
		`"items":[{"adminPassword":"` + secret + `","size":100000000000000001}],` +
		`"auth":{"passwordCredentials":{"username":"u","password":"` + secret + `"}}}`
	out := Body([]byte(in))
	if strings.Contains(out, secret) {
		t.Fatalf("nested secret leaked: %s", out)
	}
	for _, want := range []string{`"alice"`, `100000000000000001`, `"items"`} {
		if !strings.Contains(out, want) {
			t.Errorf("structure lost, missing %s in %s", want, out)
		}
	}
}

func TestBodyNonJSONPrintsOnlyLength(t *testing.T) {
	in := "-----BEGIN PRIVATE KEY-----\n" + secret + "\n-----END PRIVATE KEY-----\n"
	out := Body([]byte(in))
	if strings.Contains(out, secret) || strings.Contains(out, "PRIVATE KEY") {
		t.Fatalf("non-JSON body printed: %s", out)
	}
	if !strings.Contains(out, "bytes") {
		t.Errorf("expected length placeholder, got %q", out)
	}
	if Body(nil) != "" {
		t.Errorf("empty body should render empty")
	}
}

func TestValueRedactsRequestBodies(t *testing.T) {
	type req struct {
		Name     string `json:"name"`
		Password string `json:"dbPassword"`
	}
	cases := []interface{}{
		req{Name: "db-01", Password: secret},
		map[string]string{"password": secret},
		url.Values{"client_secret": {secret}},
		`{"password":"` + secret + `"}`,
		"password=" + secret,
		strings.NewReader(secret),
	}
	for i, c := range cases {
		if out := Value(c); strings.Contains(out, secret) {
			t.Errorf("case %d leaked: %s", i, out)
		}
	}
}

func TestURLMasksQueryAndUserinfo(t *testing.T) {
	u, _ := url.Parse("https://user:" + secret + "@api.example.com/v1/items?token=" + secret + "&limit=10")
	out := URL(u)
	if strings.Contains(out, secret) {
		t.Fatalf("URL leaked: %s", out)
	}
	if !strings.Contains(out, "/v1/items") || !strings.Contains(out, "limit=10") {
		t.Errorf("URL lost useful parts: %s", out)
	}
}

// Appkeys and API keys sent as query parameters are masked too.
func TestURLMasksAppKeyQuery(t *testing.T) {
	for _, q := range []string{"appkey", "appKey", "app_key", "api_key", "apiKey"} {
		u, _ := url.Parse("https://api.example.com/v1/items?" + q + "=" + secret + "&limit=10")
		if out := URL(u); strings.Contains(out, secret) || !strings.Contains(out, "limit=10") {
			t.Errorf("%s: %s", q, out)
		}
	}
}
