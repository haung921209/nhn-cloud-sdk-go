package ncr

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// ResponseHeader
// ---------------------------------------------------------------------------

func TestResponseHeader_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(ResponseHeader{}), []string{
		"IsSuccessful",
		"ResultCode",
		"ResultMessage",
	})
}

func TestResponseHeader_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ResponseHeader{})
	testutil.AssertStructHasJSONTag(t, typ, "IsSuccessful", "isSuccessful")
	testutil.AssertStructHasJSONTag(t, typ, "ResultCode", "resultCode")
	testutil.AssertStructHasJSONTag(t, typ, "ResultMessage", "resultMessage")
}

func TestResponseHeader_Parse(t *testing.T) {
	raw := `{"isSuccessful":true,"resultCode":200,"resultMessage":"SUCCESS"}`
	var h ResponseHeader
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !h.IsSuccessful {
		t.Error("IsSuccessful should be true")
	}
	if h.ResultCode != 200 {
		t.Errorf("ResultCode: got %d, want 200", h.ResultCode)
	}
	if h.ResultMessage != "SUCCESS" {
		t.Errorf("ResultMessage: got %q, want SUCCESS", h.ResultMessage)
	}
}

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

func TestRegistry_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Registry{}), []string{
		"ID",
		"Name",
		"URI",
		"IsPublic",
		"Status",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestRegistry_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Registry{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "URI", "registry_url")
	testutil.AssertStructHasJSONTag(t, typ, "IsPublic", "isPublic")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "creation_time")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "update_time")
}

func TestRegistry_Parse(t *testing.T) {
	raw := `{
		"id": 42,
		"name": "my-registry",
		"registry_url": "kr1-ncr.api.nhncloudservice.com/my-registry",
		"isPublic": false,
		"status": "ACTIVE",
		"creation_time": "2024-01-01T00:00:00Z",
		"update_time": "2024-06-01T00:00:00Z"
	}`
	var r Registry
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.ID != 42 {
		t.Errorf("ID: got %d, want 42", r.ID)
	}
	if r.Name != "my-registry" {
		t.Errorf("Name: got %q, want my-registry", r.Name)
	}
	if r.Status != "ACTIVE" {
		t.Errorf("Status: got %q, want ACTIVE", r.Status)
	}
	if r.CreatedAt != "2024-01-01T00:00:00Z" {
		t.Errorf("CreatedAt mismatch: %q", r.CreatedAt)
	}
}

// ---------------------------------------------------------------------------
// ListRegistriesOutput
// ---------------------------------------------------------------------------

func TestListRegistriesOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(ListRegistriesOutput{}), []string{
		"Header",
		"Registries",
	})
}

func TestListRegistriesOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ListRegistriesOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "Header", "header")
	testutil.AssertStructHasJSONTag(t, typ, "Registries", "registries")
}

func TestListRegistriesOutput_Parse(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 200, "resultMessage": "SUCCESS"},
		"registries": [
			{"id": 1, "name": "reg-a", "registry_url": "uri-a", "isPublic": true, "status": "ACTIVE",
			 "creation_time": "2024-01-01T00:00:00Z", "update_time": "2024-01-02T00:00:00Z"}
		]
	}`
	var out ListRegistriesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Header == nil || !out.Header.IsSuccessful {
		t.Error("header not parsed correctly")
	}
	if len(out.Registries) != 1 {
		t.Fatalf("Registries: got %d, want 1", len(out.Registries))
	}
	if out.Registries[0].Name != "reg-a" {
		t.Errorf("Registries[0].Name: got %q, want reg-a", out.Registries[0].Name)
	}
}

// ---------------------------------------------------------------------------
// CreateRegistryInput
// ---------------------------------------------------------------------------

func TestCreateRegistryInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(CreateRegistryInput{}), []string{
		"Name",
	})
}

func TestCreateRegistryInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateRegistryInput{})
	// API spec requires project_name for the registry name field
	testutil.AssertStructHasJSONTag(t, typ, "Name", "project_name")
	testutil.AssertStructHasJSONTag(t, typ, "IsPublic", "isPublic")
}

func TestCreateRegistryInput_Marshal(t *testing.T) {
	input := CreateRegistryInput{
		Name:     "new-registry",
		IsPublic: true,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["project_name"] != "new-registry" {
		t.Errorf("project_name: got %v", m["project_name"])
	}
	// isPublic is omitempty bool — true should appear
	if v, ok := m["isPublic"]; !ok || v != true {
		t.Errorf("isPublic: got %v (ok=%v)", v, ok)
	}
}

// ---------------------------------------------------------------------------
// UpdateRegistryInput
// ---------------------------------------------------------------------------

func TestUpdateRegistryInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(UpdateRegistryInput{})
	testutil.AssertStructHasJSONTag(t, typ, "IsPublic", "isPublic")
}

func TestUpdateRegistryInput_OmitemptyOnNilBool(t *testing.T) {
	// When IsPublic is nil, the field must be omitted so a partial update does
	// not accidentally overwrite the existing value on the server.
	input := UpdateRegistryInput{}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["isPublic"]; ok {
		t.Error("isPublic should be omitted when nil")
	}
}

// ---------------------------------------------------------------------------
// Image
// ---------------------------------------------------------------------------

func TestImage_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Image{}), []string{
		"ID",
		"Name",
		"RegistryID",
		"PullCount",
		"Size",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestImage_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Image{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "RegistryID", "registryId")
	testutil.AssertStructHasJSONTag(t, typ, "PullCount", "pullCount")
	testutil.AssertStructHasJSONTag(t, typ, "Digest", "digest")
	testutil.AssertStructHasJSONTag(t, typ, "Size", "size")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestImage_Parse(t *testing.T) {
	raw := `{
		"id": "sha256:abc123",
		"name": "nginx",
		"registryId": "reg-001",
		"pullCount": 999,
		"tags": ["latest","1.25"],
		"digest": "sha256:abc123",
		"size": 65536000,
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-06-01T00:00:00Z"
	}`
	var img Image
	if err := json.Unmarshal([]byte(raw), &img); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if img.Name != "nginx" {
		t.Errorf("Name: got %q, want nginx", img.Name)
	}
	if img.PullCount != 999 {
		t.Errorf("PullCount: got %d, want 999", img.PullCount)
	}
	if len(img.Tags) != 2 {
		t.Errorf("Tags: got %d elements, want 2", len(img.Tags))
	}
}

// ---------------------------------------------------------------------------
// Tag
// ---------------------------------------------------------------------------

func TestTag_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Tag{}), []string{
		"Name",
		"Digest",
		"Size",
		"CreatedAt",
	})
}

func TestTag_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Tag{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Digest", "digest")
	testutil.AssertStructHasJSONTag(t, typ, "Size", "size")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "LastPulledAt", "lastPulledAt")
}

func TestTag_RoundTrip(t *testing.T) {
	tag := Tag{
		Name:      "v1.0",
		Digest:    "sha256:deadbeef",
		Size:      12345,
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	testutil.AssertJSONRoundTrip(t, tag,
		`{"name":"v1.0","digest":"sha256:deadbeef","size":12345,"createdAt":"2024-01-01T00:00:00Z"}`)
}

// ---------------------------------------------------------------------------
// Vulnerability / VulnerabilitySummary / ImageScanResult
// ---------------------------------------------------------------------------

func TestVulnerability_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Vulnerability{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Package", "package")
	testutil.AssertStructHasJSONTag(t, typ, "Version", "version")
	testutil.AssertStructHasJSONTag(t, typ, "FixVersion", "fixVersion")
	testutil.AssertStructHasJSONTag(t, typ, "Severity", "severity")
}

func TestVulnerabilitySummary_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(VulnerabilitySummary{})
	testutil.AssertStructHasJSONTag(t, typ, "Critical", "critical")
	testutil.AssertStructHasJSONTag(t, typ, "High", "high")
	testutil.AssertStructHasJSONTag(t, typ, "Medium", "medium")
	testutil.AssertStructHasJSONTag(t, typ, "Low", "low")
	testutil.AssertStructHasJSONTag(t, typ, "Unknown", "unknown")
	testutil.AssertStructHasJSONTag(t, typ, "Total", "total")
}

func TestImageScanResult_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(ImageScanResult{}), []string{
		"ID",
		"ImageID",
		"Tag",
		"Digest",
		"Status",
		"ScanStartedAt",
	})
}

func TestImageScanResult_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ImageScanResult{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "ImageID", "imageId")
	testutil.AssertStructHasJSONTag(t, typ, "Tag", "tag")
	testutil.AssertStructHasJSONTag(t, typ, "Digest", "digest")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "ScanStartedAt", "scanStartedAt")
	testutil.AssertStructHasJSONTag(t, typ, "ScanCompletedAt", "scanCompletedAt")
	testutil.AssertStructHasJSONTag(t, typ, "Vulnerabilities", "vulnerabilities")
	testutil.AssertStructHasJSONTag(t, typ, "Summary", "summary")
}

func TestImageScanResult_Parse(t *testing.T) {
	raw := `{
		"id": "scan-001",
		"imageId": "img-001",
		"tag": "latest",
		"digest": "sha256:abc",
		"status": "COMPLETE",
		"scanStartedAt": "2024-01-01T00:00:00Z",
		"scanCompletedAt": "2024-01-01T00:01:00Z",
		"summary": {"critical":1,"high":2,"medium":3,"low":4,"unknown":0,"total":10},
		"vulnerabilities": [
			{"id":"CVE-2024-0001","package":"openssl","version":"1.1.1","severity":"HIGH"}
		]
	}`
	var result ImageScanResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Status != "COMPLETE" {
		t.Errorf("Status: got %q, want COMPLETE", result.Status)
	}
	if result.Summary == nil {
		t.Fatal("Summary should not be nil")
	}
	if result.Summary.Critical != 1 {
		t.Errorf("Summary.Critical: got %d, want 1", result.Summary.Critical)
	}
	if len(result.Vulnerabilities) != 1 {
		t.Fatalf("Vulnerabilities: got %d, want 1", len(result.Vulnerabilities))
	}
	if result.Vulnerabilities[0].Severity != "HIGH" {
		t.Errorf("Vulnerability.Severity: got %q, want HIGH", result.Vulnerabilities[0].Severity)
	}
}

// ---------------------------------------------------------------------------
// Webhook / CreateWebhookInput
// ---------------------------------------------------------------------------

func TestWebhook_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Webhook{}), []string{
		"ID",
		"Name",
		"RegistryID",
		"TargetURL",
		"Enabled",
		"Events",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestWebhook_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Webhook{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "RegistryID", "registryId")
	testutil.AssertStructHasJSONTag(t, typ, "TargetURL", "targetUrl")
	testutil.AssertStructHasJSONTag(t, typ, "Enabled", "enabled")
	testutil.AssertStructHasJSONTag(t, typ, "Events", "events")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestCreateWebhookInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateWebhookInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "TargetURL", "targetUrl")
	testutil.AssertStructHasJSONTag(t, typ, "Enabled", "enabled")
	testutil.AssertStructHasJSONTag(t, typ, "Events", "events")
}

func TestCreateWebhookInput_Marshal(t *testing.T) {
	input := CreateWebhookInput{
		Name:      "push-hook",
		TargetURL: "https://example.com/hook",
		Events:    []string{"PUSH_ARTIFACT", "DELETE_ARTIFACT"},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["targetUrl"] != "https://example.com/hook" {
		t.Errorf("targetUrl: got %v", m["targetUrl"])
	}
	events, ok := m["events"].([]interface{})
	if !ok || len(events) != 2 {
		t.Errorf("events: got %v", m["events"])
	}
}

func TestListWebhooksOutput_Parse(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 200, "resultMessage": "SUCCESS"},
		"webhooks": [
			{
				"id": "wh-001",
				"name": "my-hook",
				"registryId": "reg-001",
				"targetUrl": "https://example.com/hook",
				"enabled": true,
				"events": ["PUSH_ARTIFACT"],
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-02T00:00:00Z"
			}
		]
	}`
	var out ListWebhooksOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Webhooks) != 1 {
		t.Fatalf("Webhooks: got %d, want 1", len(out.Webhooks))
	}
	wh := out.Webhooks[0]
	if wh.TargetURL != "https://example.com/hook" {
		t.Errorf("TargetURL: got %q", wh.TargetURL)
	}
	if len(wh.Events) != 1 || wh.Events[0] != "PUSH_ARTIFACT" {
		t.Errorf("Events: got %v", wh.Events)
	}
}
