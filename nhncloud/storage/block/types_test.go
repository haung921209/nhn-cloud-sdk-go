package block_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/block"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ==================== Required Fields Tests ====================

func TestVolume_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.Volume{}), []string{
		"ID",
		"Name",
		"Status",
		"Size",
		"VolumeType",
		"Bootable",
		"Encrypted",
		"AvailabilityZone",
		"CreatedAt",
	})
}

func TestVolumeAttachment_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.VolumeAttachment{}), []string{
		"ID",
		"VolumeID",
		"ServerID",
		"Device",
		"AttachedAt",
	})
}

func TestSnapshot_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.Snapshot{}), []string{
		"ID",
		"Name",
		"Status",
		"Size",
		"VolumeID",
		"CreatedAt",
	})
}

func TestVolumeType_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.VolumeType{}), []string{
		"ID",
		"Name",
	})
}

func TestCreateVolumeInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.CreateVolumeInput{}), []string{
		"Size",
	})
}

func TestCreateSnapshotInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(block.CreateSnapshotInput{}), []string{
		"VolumeID",
	})
}

// ==================== JSON Tag Tests ====================

func TestVolume_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.Volume{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"Status", "status"},
		{"Size", "size"},
		{"VolumeType", "volume_type"},
		{"Bootable", "bootable"},
		{"Encrypted", "encrypted"},
		{"AvailabilityZone", "availability_zone"},
		{"SnapshotID", "snapshot_id"},
		{"SourceVolID", "source_volid"},
		{"Description", "description"},
		{"Attachments", "attachments"},
		{"Metadata", "metadata"},
		{"CreatedAt", "created_at"},
		{"UpdatedAt", "updated_at"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestVolumeAttachment_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.VolumeAttachment{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"VolumeID", "volume_id"},
		{"ServerID", "server_id"},
		{"Device", "device"},
		{"AttachedAt", "attached_at"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestSnapshot_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.Snapshot{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"Status", "status"},
		{"Size", "size"},
		{"VolumeID", "volume_id"},
		{"Description", "description"},
		{"CreatedAt", "created_at"},
		{"UpdatedAt", "updated_at"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestCreateVolumeInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.CreateVolumeInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"Name", "name"},
		{"Size", "size"},
		{"VolumeType", "volume_type"},
		{"AvailabilityZone", "availability_zone"},
		{"SnapshotID", "snapshot_id"},
		{"SourceVolID", "source_volid"},
		{"Description", "description"},
		{"Metadata", "metadata"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestAttachVolumeInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.AttachVolumeInput{})
	// API spec uses instance_uuid and mountpoint
	testutil.AssertStructHasJSONTag(t, typ, "ServerID", "instance_uuid")
	testutil.AssertStructHasJSONTag(t, typ, "Device", "mountpoint")
}

func TestExtendVolumeInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(block.ExtendVolumeInput{})
	testutil.AssertStructHasJSONTag(t, typ, "NewSize", "new_size")
}

// ==================== Response Parse Tests ====================

func TestVolume_ParseFromJSON(t *testing.T) {
	raw := `{
		"id": "vol-001",
		"name": "my-volume",
		"status": "available",
		"size": 50,
		"volume_type": "General HDD",
		"bootable": "false",
		"encrypted": false,
		"availability_zone": "kr-pub-a",
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-02T00:00:00Z",
		"description": "test volume",
		"attachments": [],
		"metadata": {}
	}`

	var v block.Volume
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if v.ID != "vol-001" {
		t.Errorf("ID: got %q, want %q", v.ID, "vol-001")
	}
	if v.Name != "my-volume" {
		t.Errorf("Name: got %q, want %q", v.Name, "my-volume")
	}
	if v.Status != "available" {
		t.Errorf("Status: got %q, want %q", v.Status, "available")
	}
	if v.Size != 50 {
		t.Errorf("Size: got %d, want %d", v.Size, 50)
	}
	if v.VolumeType != "General HDD" {
		t.Errorf("VolumeType: got %q, want %q", v.VolumeType, "General HDD")
	}
	if v.AvailabilityZone != "kr-pub-a" {
		t.Errorf("AvailabilityZone: got %q, want %q", v.AvailabilityZone, "kr-pub-a")
	}
}

func TestListVolumesOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"volumes": [
			{
				"id": "vol-001",
				"name": "vol-a",
				"status": "available",
				"size": 20,
				"volume_type": "General SSD",
				"bootable": "true",
				"encrypted": false,
				"availability_zone": "kr-pub-a",
				"created_at": "2023-01-01T00:00:00Z"
			},
			{
				"id": "vol-002",
				"name": "vol-b",
				"status": "in-use",
				"size": 100,
				"volume_type": "General HDD",
				"bootable": "false",
				"encrypted": true,
				"availability_zone": "kr-pub-b",
				"created_at": "2023-02-01T00:00:00Z"
			}
		]
	}`

	var out block.ListVolumesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.Volumes) != 2 {
		t.Fatalf("Volumes count: got %d, want 2", len(out.Volumes))
	}
	if out.Volumes[0].ID != "vol-001" {
		t.Errorf("Volumes[0].ID: got %q, want %q", out.Volumes[0].ID, "vol-001")
	}
	if out.Volumes[1].Encrypted != true {
		t.Errorf("Volumes[1].Encrypted: got %v, want true", out.Volumes[1].Encrypted)
	}
}

func TestSnapshot_ParseFromJSON(t *testing.T) {
	raw := `{
		"id": "snap-001",
		"name": "my-snapshot",
		"status": "available",
		"size": 50,
		"volume_id": "vol-001",
		"description": "daily snapshot",
		"created_at": "2023-01-01T00:00:00Z"
	}`

	var s block.Snapshot
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if s.ID != "snap-001" {
		t.Errorf("ID: got %q, want %q", s.ID, "snap-001")
	}
	if s.VolumeID != "vol-001" {
		t.Errorf("VolumeID: got %q, want %q", s.VolumeID, "vol-001")
	}
	if s.Size != 50 {
		t.Errorf("Size: got %d, want 50", s.Size)
	}
}

func TestVolumeWithAttachments_ParseFromJSON(t *testing.T) {
	raw := `{
		"volume": {
			"id": "vol-001",
			"name": "boot-vol",
			"status": "in-use",
			"size": 50,
			"volume_type": "General SSD",
			"bootable": "true",
			"encrypted": false,
			"availability_zone": "kr-pub-a",
			"created_at": "2023-01-01T00:00:00Z",
			"attachments": [
				{
					"id": "attach-001",
					"volume_id": "vol-001",
					"server_id": "server-abc",
					"device": "/dev/vda",
					"attached_at": "2023-01-01T01:00:00Z"
				}
			]
		}
	}`

	var out block.GetVolumeOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Volume.ID != "vol-001" {
		t.Errorf("Volume.ID: got %q, want %q", out.Volume.ID, "vol-001")
	}
	if len(out.Volume.Attachments) != 1 {
		t.Fatalf("Attachments count: got %d, want 1", len(out.Volume.Attachments))
	}
	att := out.Volume.Attachments[0]
	if att.ServerID != "server-abc" {
		t.Errorf("Attachment.ServerID: got %q, want %q", att.ServerID, "server-abc")
	}
	if att.Device != "/dev/vda" {
		t.Errorf("Attachment.Device: got %q, want %q", att.Device, "/dev/vda")
	}
}

func TestListVolumeTypesOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"volume_types": [
			{"id": "type-001", "name": "General HDD"},
			{"id": "type-002", "name": "General SSD", "description": "SSD storage"}
		]
	}`

	var out block.ListVolumeTypesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.VolumeTypes) != 2 {
		t.Fatalf("VolumeTypes count: got %d, want 2", len(out.VolumeTypes))
	}
	if out.VolumeTypes[0].Name != "General HDD" {
		t.Errorf("VolumeTypes[0].Name: got %q, want %q", out.VolumeTypes[0].Name, "General HDD")
	}
}

// ==================== Request Build Tests ====================

func TestCreateVolumeInput_Marshal(t *testing.T) {
	input := block.CreateVolumeInput{
		Name:             "new-vol",
		Size:             100,
		VolumeType:       "General SSD",
		AvailabilityZone: "kr-pub-a",
		Description:      "test",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["name"] != "new-vol" {
		t.Errorf("name: got %v, want %q", m["name"], "new-vol")
	}
	if m["size"].(float64) != 100 {
		t.Errorf("size: got %v, want 100", m["size"])
	}
	if m["volume_type"] != "General SSD" {
		t.Errorf("volume_type: got %v, want %q", m["volume_type"], "General SSD")
	}
	if m["availability_zone"] != "kr-pub-a" {
		t.Errorf("availability_zone: got %v, want %q", m["availability_zone"], "kr-pub-a")
	}
}

func TestCreateVolumeInput_OmitsEmptyOptionals(t *testing.T) {
	// Only required field: Size
	input := block.CreateVolumeInput{Size: 20}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := m["name"]; ok {
		t.Error("name should be omitted when empty")
	}
	if _, ok := m["snapshot_id"]; ok {
		t.Error("snapshot_id should be omitted when empty")
	}
	if _, ok := m["description"]; ok {
		t.Error("description should be omitted when empty")
	}
}

func TestCreateSnapshotInput_Marshal(t *testing.T) {
	input := block.CreateSnapshotInput{
		Name:        "snap-daily",
		VolumeID:    "vol-001",
		Description: "daily backup",
		Force:       true,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["volume_id"] != "vol-001" {
		t.Errorf("volume_id: got %v, want %q", m["volume_id"], "vol-001")
	}
	if m["name"] != "snap-daily" {
		t.Errorf("name: got %v, want %q", m["name"], "snap-daily")
	}
}

func TestUpdateVolumeInput_Marshal(t *testing.T) {
	input := block.UpdateVolumeInput{
		Name:        "updated-vol",
		Description: "updated description",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if m["name"] != "updated-vol" {
		t.Errorf("name: got %v, want %q", m["name"], "updated-vol")
	}
	if m["description"] != "updated description" {
		t.Errorf("description: got %v, want %q", m["description"], "updated description")
	}
}

func TestExtendVolumeInput_Marshal(t *testing.T) {
	testutil.AssertJSONRoundTrip(t, block.ExtendVolumeInput{NewSize: 200}, `{"new_size":200}`)
}
