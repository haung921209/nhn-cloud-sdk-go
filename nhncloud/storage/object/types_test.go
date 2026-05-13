package object_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ==================== Required Fields Tests ====================

func TestAccountInfo_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.AccountInfo{}), []string{
		"ContainerCount",
		"ObjectCount",
		"BytesUsed",
	})
}

func TestContainer_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.Container{}), []string{
		"Name",
		"Count",
		"Bytes",
	})
}

func TestObject_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.Object{}), []string{
		"Name",
		"Hash",
		"Bytes",
		"ContentType",
		"LastModified",
	})
}

func TestSLOSegment_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.SLOSegment{}), []string{
		"Path",
		"ETag",
		"SizeBytes",
	})
}

func TestPutObjectInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.PutObjectInput{}), []string{
		"Container",
		"ObjectName",
		"Body",
	})
}

func TestCopyObjectInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.CopyObjectInput{}), []string{
		"SourceContainer",
		"SourceObjectName",
		"DestinationContainer",
		"DestinationObjectName",
	})
}

func TestCreateContainerInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.CreateContainerInput{}), []string{
		"Name",
	})
}

// ==================== JSON Tag Tests ====================

func TestContainer_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(object.Container{})
	cases := []struct {
		field string
		tag   string
	}{
		{"Name", "name"},
		{"Count", "count"},
		{"Bytes", "bytes"},
		{"LastModified", "last_modified"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestObject_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(object.Object{})
	cases := []struct {
		field string
		tag   string
	}{
		{"Name", "name"},
		{"Hash", "hash"},
		{"Bytes", "bytes"},
		{"ContentType", "content_type"},
		{"LastModified", "last_modified"},
		{"Subdir", "subdir"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestSLOSegment_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(object.SLOSegment{})
	testutil.AssertStructHasJSONTag(t, typ, "Path", "path")
	testutil.AssertStructHasJSONTag(t, typ, "ETag", "etag")
	testutil.AssertStructHasJSONTag(t, typ, "SizeBytes", "size_bytes")
}

// ==================== Response Parse Tests ====================

func TestContainer_ParseFromJSON(t *testing.T) {
	raw := `{
		"name": "my-container",
		"count": 42,
		"bytes": 10240,
		"last_modified": "2023-01-01T00:00:00.000000"
	}`

	var c object.Container
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if c.Name != "my-container" {
		t.Errorf("Name: got %q, want %q", c.Name, "my-container")
	}
	if c.Count != 42 {
		t.Errorf("Count: got %d, want 42", c.Count)
	}
	if c.Bytes != 10240 {
		t.Errorf("Bytes: got %d, want 10240", c.Bytes)
	}
	if c.LastModified != "2023-01-01T00:00:00.000000" {
		t.Errorf("LastModified: got %q, want %q", c.LastModified, "2023-01-01T00:00:00.000000")
	}
}

func TestListContainersOutput_ParseFromJSON(t *testing.T) {
	// Object Storage returns a JSON array at the account level
	raw := `[
		{"name": "container-a", "count": 10, "bytes": 1024},
		{"name": "container-b", "count": 5,  "bytes": 512}
	]`

	var containers []object.Container
	if err := json.Unmarshal([]byte(raw), &containers); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(containers) != 2 {
		t.Fatalf("containers count: got %d, want 2", len(containers))
	}
	if containers[0].Name != "container-a" {
		t.Errorf("containers[0].Name: got %q, want %q", containers[0].Name, "container-a")
	}
	if containers[1].Count != 5 {
		t.Errorf("containers[1].Count: got %d, want 5", containers[1].Count)
	}
}

func TestObject_ParseFromJSON(t *testing.T) {
	raw := `{
		"name": "images/photo.jpg",
		"hash": "d41d8cd98f00b204e9800998ecf8427e",
		"bytes": 204800,
		"content_type": "image/jpeg",
		"last_modified": "2023-06-15T12:00:00.000000"
	}`

	var o object.Object
	if err := json.Unmarshal([]byte(raw), &o); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if o.Name != "images/photo.jpg" {
		t.Errorf("Name: got %q, want %q", o.Name, "images/photo.jpg")
	}
	if o.Hash != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("Hash: got %q, want %q", o.Hash, "d41d8cd98f00b204e9800998ecf8427e")
	}
	if o.Bytes != 204800 {
		t.Errorf("Bytes: got %d, want 204800", o.Bytes)
	}
	if o.ContentType != "image/jpeg" {
		t.Errorf("ContentType: got %q, want %q", o.ContentType, "image/jpeg")
	}
}

func TestObject_SubdirParseFromJSON(t *testing.T) {
	// When listing with delimiter, subdirectory entries have a "subdir" field
	raw := `{"subdir": "images/"}`

	var o object.Object
	if err := json.Unmarshal([]byte(raw), &o); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if o.Subdir != "images/" {
		t.Errorf("Subdir: got %q, want %q", o.Subdir, "images/")
	}
}

func TestSLOSegment_ParseFromJSON(t *testing.T) {
	raw := `{
		"path": "/my-container/object-segment/001",
		"etag": "b026324c6904b2a9cb4b88d6d61c81d1",
		"size_bytes": 5242880
	}`

	var seg object.SLOSegment
	if err := json.Unmarshal([]byte(raw), &seg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if seg.Path != "/my-container/object-segment/001" {
		t.Errorf("Path: got %q, want %q", seg.Path, "/my-container/object-segment/001")
	}
	if seg.ETag != "b026324c6904b2a9cb4b88d6d61c81d1" {
		t.Errorf("ETag: got %q, want %q", seg.ETag, "b026324c6904b2a9cb4b88d6d61c81d1")
	}
	if seg.SizeBytes != 5242880 {
		t.Errorf("SizeBytes: got %d, want 5242880", seg.SizeBytes)
	}
}

func TestGetSLOManifestOutput_ParseSegments(t *testing.T) {
	raw := `[
		{
			"path": "/container/obj/001",
			"etag": "abc123",
			"size_bytes": 1048576
		},
		{
			"path": "/container/obj/002",
			"etag": "def456",
			"size_bytes": 524288
		}
	]`

	var segs []object.SLOSegment
	if err := json.Unmarshal([]byte(raw), &segs); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(segs) != 2 {
		t.Fatalf("segments count: got %d, want 2", len(segs))
	}
	if segs[0].ETag != "abc123" {
		t.Errorf("segs[0].ETag: got %q, want %q", segs[0].ETag, "abc123")
	}
	if segs[1].SizeBytes != 524288 {
		t.Errorf("segs[1].SizeBytes: got %d, want 524288", segs[1].SizeBytes)
	}
}

// ==================== Request Build Tests ====================

func TestSLOSegment_Marshal(t *testing.T) {
	seg := object.SLOSegment{
		Path:      "/my-container/obj/001",
		ETag:      "abc123",
		SizeBytes: 5242880,
	}

	testutil.AssertJSONRoundTrip(t, seg, `{"path":"/my-container/obj/001","etag":"abc123","size_bytes":5242880}`)
}

func TestSLOSegments_MarshalArray(t *testing.T) {
	segs := []object.SLOSegment{
		{Path: "/c/o/001", ETag: "etag1", SizeBytes: 1048576},
		{Path: "/c/o/002", ETag: "etag2", SizeBytes: 1048576},
	}

	data, err := json.Marshal(segs)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var back []object.SLOSegment
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(back) != 2 {
		t.Fatalf("round-tripped count: got %d, want 2", len(back))
	}
	if back[0].Path != "/c/o/001" {
		t.Errorf("back[0].Path: got %q, want %q", back[0].Path, "/c/o/001")
	}
}

// ==================== Storage Class Constants Tests ====================

func TestStorageClassConstants(t *testing.T) {
	if object.StorageClassStandard != "Standard" {
		t.Errorf("StorageClassStandard: got %q, want %q", object.StorageClassStandard, "Standard")
	}
	if object.StorageClassEconomy != "Economy" {
		t.Errorf("StorageClassEconomy: got %q, want %q", object.StorageClassEconomy, "Economy")
	}
}

// ==================== Struct Type Tests ====================

func TestAccountInfo_FieldTypes(t *testing.T) {
	typ := reflect.TypeOf(object.AccountInfo{})

	for _, fieldName := range []string{"ContainerCount", "ObjectCount", "BytesUsed"} {
		f, ok := typ.FieldByName(fieldName)
		if !ok {
			t.Errorf("AccountInfo has no field %s", fieldName)
			continue
		}
		if f.Type.Kind() != reflect.Int64 {
			t.Errorf("AccountInfo.%s type: got %v, want int64", fieldName, f.Type)
		}
	}
}

func TestContainer_FieldTypes(t *testing.T) {
	typ := reflect.TypeOf(object.Container{})

	for _, fieldName := range []string{"Count", "Bytes"} {
		f, ok := typ.FieldByName(fieldName)
		if !ok {
			t.Errorf("Container has no field %s", fieldName)
			continue
		}
		if f.Type.Kind() != reflect.Int64 {
			t.Errorf("Container.%s type: got %v, want int64", fieldName, f.Type)
		}
	}
}

func TestObject_FieldTypes(t *testing.T) {
	f, ok := reflect.TypeOf(object.Object{}).FieldByName("Bytes")
	if !ok {
		t.Fatal("Object has no field Bytes")
	}
	if f.Type.Kind() != reflect.Int64 {
		t.Errorf("Object.Bytes type: got %v, want int64", f.Type)
	}
}

func TestSLOSegment_SizeBytesIsInt64(t *testing.T) {
	f, ok := reflect.TypeOf(object.SLOSegment{}).FieldByName("SizeBytes")
	if !ok {
		t.Fatal("SLOSegment has no field SizeBytes")
	}
	if f.Type.Kind() != reflect.Int64 {
		t.Errorf("SLOSegment.SizeBytes type: got %v, want int64", f.Type)
	}
}

func TestPutObjectOutput_HasETag(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.PutObjectOutput{}), []string{"ETag"})
}

func TestGetObjectOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(object.GetObjectOutput{}), []string{
		"Body",
		"ContentType",
		"ContentLength",
		"ETag",
	})
}
