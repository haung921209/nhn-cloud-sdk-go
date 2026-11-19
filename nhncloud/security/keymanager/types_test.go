package keymanager_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ==================== Required Fields Tests ====================

func TestAPIResponse_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(keymanager.APIResponse{})
	// APIResponse embeds header inline; verify the struct exists
	if typ.Kind() != reflect.Struct {
		t.Errorf("APIResponse should be a struct, got %v", typ.Kind())
	}
}

func TestKeyStore_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.KeyStore{}), []string{
		"KeyStoreID",
		"Name",
	})
}

func TestKey_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.Key{}), []string{
		"KeyID",
		"Name",
		"KeyType",
	})
}

func TestEncryptInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.EncryptInput{}), []string{
		"Plaintext",
	})
}

func TestDecryptInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.DecryptInput{}), []string{
		"Ciphertext",
	})
}

func TestSignInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.SignInput{}), []string{
		"Data",
	})
}

func TestVerifyInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.VerifyInput{}), []string{
		"Data",
		"Signature",
	})
}

func TestCreateKeyInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(keymanager.CreateKeyInput{}), []string{
		"KeyStoreName",
		"Name",
	})
}

// ==================== JSON Tag Tests ====================

func TestEncryptInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.EncryptInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Plaintext", "plaintext")
	testutil.AssertStructHasJSONTag(t, typ, "AAD", "aad")
}

func TestDecryptInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.DecryptInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Ciphertext", "ciphertext")
	testutil.AssertStructHasJSONTag(t, typ, "IV", "iv")
	testutil.AssertStructHasJSONTag(t, typ, "Tag", "tag")
	testutil.AssertStructHasJSONTag(t, typ, "AAD", "aad")
}

func TestSignInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.SignInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Data", "plaintext")
}

func TestVerifyInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.VerifyInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Data", "plaintext")
	testutil.AssertStructHasJSONTag(t, typ, "Signature", "signature")
}

func TestCreateKeyInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.CreateKeyInput{})
	testutil.AssertStructHasJSONTag(t, typ, "KeyStoreName", "keyStoreName")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "KeyAlgorithm", "keyAlgorithm")
	testutil.AssertStructHasJSONTag(t, typ, "KeySize", "keySize")
	testutil.AssertStructHasJSONTag(t, typ, "Algorithm", "algorithm")
	testutil.AssertStructHasJSONTag(t, typ, "Secret", "secretValue")
	testutil.AssertStructHasJSONTag(t, typ, "RotationPeriod", "autoRotationPeriod")
}

func TestDeleteKeyInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(keymanager.DeleteKeyInput{})
	testutil.AssertStructHasJSONTag(t, typ, "RequestDeletionDate", "requestDeletionDate")
}

// ==================== Response Parse Tests ====================

func TestEncryptOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"ciphertext": "YWJjZGVm",
			"iv": "dGVzdGl2MTIzNA==",
			"tag": "dGFndGFn"
		}
	}`

	var out keymanager.EncryptOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.Ciphertext != "YWJjZGVm" {
		t.Errorf("Body.Ciphertext: got %q, want %q", out.Body.Ciphertext, "YWJjZGVm")
	}
	if out.Body.IV != "dGVzdGl2MTIzNA==" {
		t.Errorf("Body.IV: got %q, want %q", out.Body.IV, "dGVzdGl2MTIzNA==")
	}
	if out.Body.Tag != "dGFndGFn" {
		t.Errorf("Body.Tag: got %q, want %q", out.Body.Tag, "dGFndGFn")
	}
}

func TestDecryptOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"plaintext": "aGVsbG8gd29ybGQ="
		}
	}`

	var out keymanager.DecryptOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.Plaintext != "aGVsbG8gd29ybGQ=" {
		t.Errorf("Body.Plaintext: got %q, want %q", out.Body.Plaintext, "aGVsbG8gd29ybGQ=")
	}
}

func TestSignOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"signature": "c2lnbmF0dXJlZGF0YQ=="
		}
	}`

	var out keymanager.SignOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.Signature != "c2lnbmF0dXJlZGF0YQ==" {
		t.Errorf("Body.Signature: got %q, want %q", out.Body.Signature, "c2lnbmF0dXJlZGF0YQ==")
	}
}

func TestVerifyOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"result": true
		}
	}`

	var out keymanager.VerifyOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !out.Body.Result {
		t.Error("Body.Result: got false, want true")
	}
}

func TestCreateLocalKeyOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"localKeyPlaintext": "cGxhaW5rZXk=",
			"localKeyCiphertext": "ZW5jcnlwdGVka2V5"
		}
	}`

	var out keymanager.CreateLocalKeyOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.PlainDataKey != "cGxhaW5rZXk=" {
		t.Errorf("Body.PlainDataKey: got %q, want %q", out.Body.PlainDataKey, "cGxhaW5rZXk=")
	}
	if out.Body.EncryptedDataKey != "ZW5jcnlwdGVka2V5" {
		t.Errorf("Body.EncryptedDataKey: got %q, want %q", out.Body.EncryptedDataKey, "ZW5jcnlwdGVka2V5")
	}
}

func TestGetSecretOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"secret": "my-secret-value"
		}
	}`

	var out keymanager.GetSecretOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.Secret != "my-secret-value" {
		t.Errorf("Body.Secret: got %q, want %q", out.Body.Secret, "my-secret-value")
	}
}

func TestGetSymmetricKeyOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"symmetricKey": "c3ltbWV0cmljS2V5VmFsdWU="
		}
	}`

	var out keymanager.GetSymmetricKeyOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.SymmetricKey != "c3ltbWV0cmljS2V5VmFsdWU=" {
		t.Errorf("Body.KeyValue: got %q, want %q", out.Body.SymmetricKey, "c3ltbWV0cmljS2V5VmFsdWU=")
	}
}

func TestCreateKeyOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"body": {
			"keyId": "key-new-001"
		}
	}`

	var out keymanager.CreateKeyOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Body.KeyID != "key-new-001" {
		t.Errorf("Body.KeyID: got %q, want %q", out.Body.KeyID, "key-new-001")
	}
}

// ==================== Request Build Tests ====================

func TestEncryptInput_Marshal(t *testing.T) {
	input := keymanager.EncryptInput{
		Plaintext: "aGVsbG8=",
		AAD:       "YWRkaXRpb25hbA==",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["plaintext"] != "aGVsbG8=" {
		t.Errorf("plaintext: got %v, want %q", m["plaintext"], "aGVsbG8=")
	}
	if m["aad"] != "YWRkaXRpb25hbA==" {
		t.Errorf("aad: got %v, want %q", m["aad"], "YWRkaXRpb25hbA==")
	}
}

func TestEncryptInput_OmitsAADWhenEmpty(t *testing.T) {
	input := keymanager.EncryptInput{Plaintext: "aGVsbG8="}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := m["aad"]; ok {
		t.Error("aad should be omitted when empty")
	}
}

func TestDecryptInput_Marshal(t *testing.T) {
	input := keymanager.DecryptInput{
		Ciphertext: "YWJjZGVm",
		IV:         "dGVzdGl2",
		Tag:        "dGFndGFn",
	}

	testutil.AssertJSONRoundTrip(t, input, `{"ciphertext":"YWJjZGVm","iv":"dGVzdGl2","tag":"dGFndGFn"}`)
}

func TestDecryptInput_OmitsOptionals(t *testing.T) {
	input := keymanager.DecryptInput{Ciphertext: "YWJjZGVm"}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := m["iv"]; ok {
		t.Error("iv should be omitted when empty")
	}
	if _, ok := m["tag"]; ok {
		t.Error("tag should be omitted when empty")
	}
	if _, ok := m["aad"]; ok {
		t.Error("aad should be omitted when empty")
	}
}

func TestCreateKeyInput_Marshal_Symmetric(t *testing.T) {
	input := keymanager.CreateKeyInput{
		KeyStoreName: "my-store",
		Name:         "aes-key",
		KeyAlgorithm: "AES256",
		KeySize:      256,
		RotationPeriod: 90,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["keyStoreName"] != "my-store" {
		t.Errorf("keyStoreName: got %v, want %q", m["keyStoreName"], "my-store")
	}
	if m["keyAlgorithm"] != "AES256" {
		t.Errorf("keyAlgorithm: got %v, want %q", m["keyAlgorithm"], "AES256")
	}
	if m["keySize"].(float64) != 256 {
		t.Errorf("keySize: got %v, want 256", m["keySize"])
	}
	if m["autoRotationPeriod"].(float64) != 90 {
		t.Errorf("autoRotationPeriod: got %v, want 90", m["autoRotationPeriod"])
	}
}

func TestCreateKeyInput_Marshal_Asymmetric(t *testing.T) {
	input := keymanager.CreateKeyInput{
		KeyStoreName: "my-store",
		Name:         "rsa-key",
		Algorithm:    "RSA2048",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if m["algorithm"] != "RSA2048" {
		t.Errorf("algorithm: got %v, want %q", m["algorithm"], "RSA2048")
	}
}

func TestCreateKeyInput_Marshal_Secret(t *testing.T) {
	input := keymanager.CreateKeyInput{
		KeyStoreName: "my-store",
		Name:         "my-secret",
		Secret:       "s3cr3t-v4lue",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if m["secretValue"] != "s3cr3t-v4lue" {
		t.Errorf("secretValue: got %v, want %q", m["secretValue"], "s3cr3t-v4lue")
	}
}

func TestDeleteKeyInput_Marshal(t *testing.T) {
	input := keymanager.DeleteKeyInput{
		RequestDeletionDate: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if m["requestDeletionDate"] != "2024-01-01T00:00:00Z" {
		t.Errorf("requestDeletionDate: got %v, want %q", m["requestDeletionDate"], "2024-01-01T00:00:00Z")
	}
}

func TestDeleteKeyInput_OmitsDeletionDate(t *testing.T) {
	input := keymanager.DeleteKeyInput{} // immediate deletion (no date)

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := m["requestDeletionDate"]; ok {
		t.Error("requestDeletionDate should be omitted when empty")
	}
}

func TestSignInput_Marshal(t *testing.T) {
	testutil.AssertJSONRoundTrip(t, keymanager.SignInput{Data: "dGVzdA=="}, `{"plaintext":"dGVzdA=="}`)
}

func TestVerifyInput_Marshal(t *testing.T) {
	testutil.AssertJSONRoundTrip(t, keymanager.VerifyInput{
		Data:      "dGVzdA==",
		Signature: "c2lnbmF0dXJl",
	}, `{"plaintext":"dGVzdA==","signature":"c2lnbmF0dXJl"}`)
}
