package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// A token cached for one User Access Key must not be used with another key.
func TestTokenCacheIsPerAccessKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if TokenCacheFile(home, "KEY-A") == TokenCacheFile(home, "KEY-B") {
		t.Fatal("two access keys share one cache file")
	}
	data, _ := json.Marshal(TokenCache{Token: "token-of-A", ExpiresAt: time.Now().Add(time.Hour)})
	os.MkdirAll(filepath.Join(home, ".nhncloud"), 0700)
	if err := os.WriteFile(TokenCacheFile(home, "KEY-A"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if a := NewBearerAuthWithAutoRefresh("app", "KEY-B", "secret"); a.token != "" {
		t.Fatalf("key B loaded key A's token %q", a.token)
	}
	if a := NewBearerAuthWithAutoRefresh("app", "KEY-A", "secret"); a.token != "token-of-A" {
		t.Fatalf("key A did not load its own token: %q", a.token)
	}
}
