package keymanager

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Response bodies from the API v1.2 guide examples must decode.
func TestOutputsDecodeGuideExamples(t *testing.T) {
	var ks ListKeyStoresOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyStoreList":[{"keyStoreId":1,"name":"n","ip4AuthUse":"Y","macAuthUse":"N"}]}}`), &ks)
	if len(ks.Body.KeyStores) != 1 || ks.Body.KeyStores[0].KeyStoreID != 1 || ks.Body.KeyStores[0].Name != "n" {
		t.Errorf("key stores %+v", ks.Body)
	}
	var one GetKeyStoreOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyStoreId":2,"name":"n","macAuthUse":"Y"}}`), &one)
	if one.Body.KeyStoreID != 2 || one.Body.MACAuthUse != "Y" {
		t.Errorf("key store %+v", one.Body)
	}
	var keys ListKeysOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyList":[{"keyId":"k","name":"n","keyType":"SYMMETRIC_KEY","currentKeyValueVersion":2}]}}`), &keys)
	if len(keys.Body.Keys) != 1 || keys.Body.Keys[0].KeyID != "k" || keys.Body.Keys[0].CurrentKeyValueVersion != 2 {
		t.Errorf("keys %+v", keys.Body)
	}
	var key GetKeyOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyId":"k","keyType":"SECRET","deletionDatetime":null}}`), &key)
	if key.Body.KeyID != "k" || key.Body.KeyType != "SECRET" {
		t.Errorf("key %+v", key.Body)
	}
	var pub GetPublicKeyOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyType":"PublicKey","key":"0x30","encodedKey":"MIIB","keyVersion":0}}`), &pub)
	if pub.Body.EncodedKey != "MIIB" || pub.Body.KeyType != "PublicKey" {
		t.Errorf("public key %+v", pub.Body)
	}
	var del DeleteKeyOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"keyId":"k","deletionDateTime":"2023-11-20T22:00:00.00"}}`), &del)
	if del.Body.DeletionDatetime == "" {
		t.Errorf("delete %+v", del.Body)
	}
	var info GetClientInfoOutput
	json.Unmarshal([]byte(`{"header":{"isSuccessful":true},"body":{"clientIp":"192.0.2.1","clientMacHeader":"02:00:00:00:00:01","clientSentCertificate":false}}`), &info)
	if info.Body.ClientIP != "192.0.2.1" || info.Body.ClientMACHeader != "02:00:00:00:00:01" {
		t.Errorf("client info %+v", info.Body)
	}
}

// Requests use API v1.2 (the version that takes the X-TC-AUTHENTICATION
// pair), list key stores with details, and send the client MAC when set.
func TestRequestShape(t *testing.T) {
	var got *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		w.Write([]byte(`{"header":{"isSuccessful":true}}`))
	}))
	defer srv.Close()
	c := NewClient("KR1", "ak", "id", "secret", false)
	c.baseURL = srv.URL
	c.ClientMAC = "02:00:00:00:00:01"
	if _, err := c.ListKeyStores(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got.URL.Path != "/keymanager/v1.2/appkey/ak/keystores" || got.URL.RawQuery != "detail=true" ||
		got.Header.Get("X-TOAST-CLIENT-MAC-ADDR") != "02:00:00:00:00:01" || got.Header.Get("X-TC-AUTHENTICATION-ID") != "id" {
		t.Errorf("request %s?%s headers %v", got.URL.Path, got.URL.RawQuery, got.Header)
	}
	c.ListKeys(context.Background(), "1")
	if got.URL.Path != "/keymanager/v1.2/appkey/ak/keystores/1/keys" || got.URL.RawQuery != "detail=true" {
		t.Errorf("list keys request %s?%s", got.URL.Path, got.URL.RawQuery)
	}
}
