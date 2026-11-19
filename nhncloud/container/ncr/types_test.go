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
// Response shapes from the API guide
// ---------------------------------------------------------------------------

// The registry ID is project_id; describe-registries printed 0 because it
// read a field "id" the API does not send.
func TestRegistry_ParseGuideShape(t *testing.T) {
	raw := `{"header":{"isSuccessful":true,"resultCode":200,"resultMessage":"SUCCESS"},
		"registries":[{"project_id":2,"name":"hy","creation_time":"2023-03-14T02:41:10.525Z",
		"repo_count":2,"metadata":{"public":"false"},"registry_id":0,"usage":57464288,
		"uri":"example.invalid/hy"}]}`
	var out ListRegistriesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	r := out.Registries[0]
	if r.ID != 2 || r.Name != "hy" || r.URI != "example.invalid/hy" || r.RepoCount != 2 || r.Metadata["public"] != "false" {
		t.Errorf("registry %+v", r)
	}
	var one GetRegistryOutput
	if err := json.Unmarshal([]byte(`{"registry":{"project_id":3,"name":"x"}}`), &one); err != nil {
		t.Fatal(err)
	}
	if one.Registry.ID != 3 || one.Registry.Name != "x" {
		t.Errorf("get registry %+v", one.Registry)
	}
}

// images[].id and the counters are numbers.
func TestImage_ParseGuideShape(t *testing.T) {
	raw := `{"images":[{"id":1,"name":"hy/busybox","project_id":2,"description":"",
		"pull_count":8,"artifact_count":1,"creation_time":"2023-03-27T01:55:13.176Z"}]}`
	var out ListImagesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	img := out.Images[0]
	if img.ID != 1 || img.ProjectID != 2 || img.PullCount != 8 || img.ArtifactCount != 1 {
		t.Errorf("image %+v", img)
	}
	var one GetImageOutput
	if err := json.Unmarshal([]byte(`{"image":{"id":5,"name":"hy/busybox"}}`), &one); err != nil {
		t.Fatal(err)
	}
	if one.Image.ID != 5 {
		t.Errorf("get image %+v", one.Image)
	}
}

func TestArtifact_ParseGuideShape(t *testing.T) {
	raw := `{"artifacts":[{"digest":"sha256:ba70","id":1,"size":2003144,"type":"IMAGE",
		"push_time":"2023-03-27T01:55:14.082Z","tags":[{"artifact_id":1,"id":1,"name":"latest","immutable":false}]}]}`
	var out ListArtifactsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	a := out.Artifacts[0]
	if a.Digest != "sha256:ba70" || a.Size != 2003144 || len(a.Tags) != 1 || a.Tags[0].Name != "latest" {
		t.Errorf("artifact %+v", a)
	}
}

func TestWebhook_GuideShape(t *testing.T) {
	raw := `{"policies":[{"creation_time":"2023-03-27T22:53:16.719Z","enabled":true,
		"event_types":["PUSH_ARTIFACT"],"id":1,"name":"test","project_id":2,
		"targets":[{"address":"https://example.invalid","skip_cert_verify":true,"type":"http"}]}]}`
	var out ListWebhooksOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	wh := out.Webhooks[0]
	if wh.ID != 1 || wh.Name != "test" || len(wh.EventTypes) != 1 || wh.Targets[0].Address != "https://example.invalid" {
		t.Errorf("webhook %+v", wh)
	}
	b, _ := json.Marshal(CreateWebhookInput{Name: "h", EventTypes: []string{"PUSH_ARTIFACT"},
		Targets: []WebhookTarget{{Address: "https://example.invalid", Type: "http"}}})
	want := `{"name":"h","enabled":false,"event_types":["PUSH_ARTIFACT"],"targets":[{"address":"https://example.invalid","type":"http"}]}`
	if string(b) != want {
		t.Errorf("create body %s, want %s", b, want)
	}
}
