package nas_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/nas"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ==================== Required Fields Tests ====================

func TestHeader_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.Header{}), []string{
		"IsSuccessful",
		"ResultCode",
		"ResultMessage",
	})
}

func TestPaging_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.Paging{}), []string{
		"Limit",
		"Page",
		"TotalCount",
	})
}

func TestVolume_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.Volume{}), []string{
		"ID",
		"Name",
		"Status",
		"SizeGB",
		"ProjectID",
		"TenantID",
		"ACL",
		"Encryption",
		"MountProtocol",
		"SnapshotPolicy",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestInterface_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.Interface{}), []string{
		"ID",
		"Path",
		"Status",
		"SubnetID",
		"TenantID",
	})
}

func TestSnapshot_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.Snapshot{}), []string{
		"ID",
		"Name",
		"Preserved",
		"Size",
		"CreatedAt",
	})
}

func TestCreateVolumeInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.CreateVolumeInput{}), []string{
		"Name",
		"SizeGB",
	})
}

func TestCreateSnapshotInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.CreateSnapshotInput{}), []string{
		"Name",
	})
}

func TestCreateInterfaceInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(nas.CreateInterfaceInput{}), []string{
		"SubnetID",
	})
}

// ==================== JSON Tag Tests ====================

func TestHeader_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Header{})
	testutil.AssertStructHasJSONTag(t, typ, "IsSuccessful", "isSuccessful")
	testutil.AssertStructHasJSONTag(t, typ, "ResultCode", "resultCode")
	testutil.AssertStructHasJSONTag(t, typ, "ResultMessage", "resultMessage")
}

func TestPaging_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Paging{})
	testutil.AssertStructHasJSONTag(t, typ, "Limit", "limit")
	testutil.AssertStructHasJSONTag(t, typ, "Page", "page")
	testutil.AssertStructHasJSONTag(t, typ, "TotalCount", "totalCount")
}

func TestVolume_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Volume{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"Status", "status"},
		{"Description", "description"},
		{"SizeGB", "sizeGb"},
		{"ProjectID", "projectId"},
		{"TenantID", "tenantId"},
		{"ACL", "acl"},
		{"Encryption", "encryption"},
		{"Interfaces", "interfaces"},
		{"Mirrors", "mirrors"},
		{"MountProtocol", "mountProtocol"},
		{"SnapshotPolicy", "snapshotPolicy"},
		{"CreatedAt", "createdAt"},
		{"UpdatedAt", "updatedAt"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(t *testing.T) {
			testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.tag)
		})
	}
}

func TestInterface_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Interface{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Path", "path")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "SubnetID", "subnetId")
	testutil.AssertStructHasJSONTag(t, typ, "TenantID", "tenantId")
}

func TestEncryption_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Encryption{})
	testutil.AssertStructHasJSONTag(t, typ, "Enabled", "enabled")
	testutil.AssertStructHasJSONTag(t, typ, "Keys", "keys")
}

func TestEncryptionKey_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.EncryptionKey{})
	testutil.AssertStructHasJSONTag(t, typ, "KeyID", "keyId")
	testutil.AssertStructHasJSONTag(t, typ, "Version", "version")
}

func TestMountProtocol_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.MountProtocol{})
	testutil.AssertStructHasJSONTag(t, typ, "CIFSAuthIDs", "cifsAuthIds")
	testutil.AssertStructHasJSONTag(t, typ, "Protocol", "protocol")
}

func TestSnapshotPolicy_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.SnapshotPolicy{})
	testutil.AssertStructHasJSONTag(t, typ, "MaxScheduledCount", "maxScheduledCount")
	testutil.AssertStructHasJSONTag(t, typ, "ReservePercent", "reservePercent")
	testutil.AssertStructHasJSONTag(t, typ, "Schedule", "schedule")
}

func TestSnapshotSchedule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.SnapshotSchedule{})
	testutil.AssertStructHasJSONTag(t, typ, "Time", "time")
	testutil.AssertStructHasJSONTag(t, typ, "TimeOffset", "timeOffset")
}

func TestSnapshot_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.Snapshot{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Preserved", "preserved")
	testutil.AssertStructHasJSONTag(t, typ, "Size", "size")
	testutil.AssertStructHasJSONTag(t, typ, "ReclaimableSpace", "reclaimableSpace")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
}

func TestCreateVolumeInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.CreateVolumeInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "SizeGB", "sizeGb")
	testutil.AssertStructHasJSONTag(t, typ, "Encryption", "encryption")
	testutil.AssertStructHasJSONTag(t, typ, "Interfaces", "interfaces")
	testutil.AssertStructHasJSONTag(t, typ, "MountProtocol", "mountProtocol")
	testutil.AssertStructHasJSONTag(t, typ, "SnapshotPolicy", "snapshotPolicy")
}

func TestCreateVolumeMirrorInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.CreateVolumeMirrorInput{})
	testutil.AssertStructHasJSONTag(t, typ, "DstRegion", "dstRegion")
	testutil.AssertStructHasJSONTag(t, typ, "DstTenantID", "dstTenantId")
	testutil.AssertStructHasJSONTag(t, typ, "DstVolume", "dstVolume")
}

func TestVolumeMirror_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(nas.VolumeMirror{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "DstProjectID", "dstProjectId")
	testutil.AssertStructHasJSONTag(t, typ, "DstRegion", "dstRegion")
	testutil.AssertStructHasJSONTag(t, typ, "DstTenantID", "dstTenantId")
	testutil.AssertStructHasJSONTag(t, typ, "DstVolumeID", "dstVolumeId")
	testutil.AssertStructHasJSONTag(t, typ, "SrcVolumeID", "srcVolumeId")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
}

// ==================== Response Parse Tests ====================

func TestGetVolumeOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {
			"isSuccessful": true,
			"resultCode": 0,
			"resultMessage": "SUCCESS"
		},
		"volume": {
			"id": "vol-nas-001",
			"name": "my-nas",
			"status": "available",
			"description": null,
			"sizeGb": 300,
			"projectId": "proj-abc",
			"tenantId": "tenant-xyz",
			"acl": ["0.0.0.0/0"],
			"encryption": {
				"enabled": false,
				"keys": []
			},
			"interfaces": [],
			"mirrors": [],
			"mountProtocol": {
				"cifsAuthIds": [],
				"protocol": "NFS"
			},
			"snapshotPolicy": {
				"maxScheduledCount": 10,
				"reservePercent": 20,
				"schedule": {
					"time": "02:00",
					"timeOffset": "+09:00"
				}
			},
			"createdAt": "2023-01-01T00:00:00Z",
			"updatedAt": "2023-01-02T00:00:00Z"
		}
	}`

	var out nas.GetVolumeOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !out.Header.IsSuccessful {
		t.Error("Header.IsSuccessful: got false, want true")
	}
	if out.Header.ResultCode != 0 {
		t.Errorf("Header.ResultCode: got %d, want 0", out.Header.ResultCode)
	}
	if out.Volume.ID != "vol-nas-001" {
		t.Errorf("Volume.ID: got %q, want %q", out.Volume.ID, "vol-nas-001")
	}
	if out.Volume.SizeGB != 300 {
		t.Errorf("Volume.SizeGB: got %d, want 300", out.Volume.SizeGB)
	}
	if out.Volume.MountProtocol.Protocol != "NFS" {
		t.Errorf("Volume.MountProtocol.Protocol: got %q, want %q", out.Volume.MountProtocol.Protocol, "NFS")
	}
	if out.Volume.SnapshotPolicy.MaxScheduledCount != 10 {
		t.Errorf("SnapshotPolicy.MaxScheduledCount: got %d, want 10", out.Volume.SnapshotPolicy.MaxScheduledCount)
	}
	if out.Volume.SnapshotPolicy.Schedule.TimeOffset != "+09:00" {
		t.Errorf("SnapshotPolicy.Schedule.TimeOffset: got %q, want %q", out.Volume.SnapshotPolicy.Schedule.TimeOffset, "+09:00")
	}
	if len(out.Volume.ACL) != 1 || out.Volume.ACL[0] != "0.0.0.0/0" {
		t.Errorf("Volume.ACL: got %v, want [0.0.0.0/0]", out.Volume.ACL)
	}
}

func TestListVolumesOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 0, "resultMessage": "SUCCESS"},
		"paging": {"limit": 10, "page": 1, "totalCount": 2},
		"volumes": [
			{
				"id": "vol-001",
				"name": "nas-a",
				"status": "available",
				"description": null,
				"sizeGb": 100,
				"projectId": "proj-1",
				"tenantId": "tenant-1",
				"acl": [],
				"encryption": {"enabled": false, "keys": []},
				"interfaces": [],
				"mirrors": [],
				"mountProtocol": {"cifsAuthIds": [], "protocol": "NFS"},
				"snapshotPolicy": {"maxScheduledCount": 5, "reservePercent": 10, "schedule": {"time": "03:00", "timeOffset": "+09:00"}},
				"createdAt": "2023-01-01T00:00:00Z",
				"updatedAt": "2023-01-01T00:00:00Z"
			}
		]
	}`

	var out nas.ListVolumesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Paging.TotalCount != 2 {
		t.Errorf("Paging.TotalCount: got %d, want 2", out.Paging.TotalCount)
	}
	if len(out.Volumes) != 1 {
		t.Fatalf("Volumes count: got %d, want 1", len(out.Volumes))
	}
	if out.Volumes[0].ID != "vol-001" {
		t.Errorf("Volumes[0].ID: got %q, want %q", out.Volumes[0].ID, "vol-001")
	}
}

func TestListSnapshotsOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 0, "resultMessage": "SUCCESS"},
		"snapshots": [
			{
				"id": "snap-001",
				"name": "daily-snap",
				"preserved": false,
				"size": 100,
				"createdAt": "2023-06-01T00:00:00Z"
			}
		]
	}`

	var out nas.ListSnapshotsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.Snapshots) != 1 {
		t.Fatalf("Snapshots count: got %d, want 1", len(out.Snapshots))
	}
	if out.Snapshots[0].ID != "snap-001" {
		t.Errorf("Snapshots[0].ID: got %q, want %q", out.Snapshots[0].ID, "snap-001")
	}
	if out.Snapshots[0].Preserved != false {
		t.Error("Snapshots[0].Preserved: got true, want false")
	}
}

func TestCreateInterfaceOutput_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 0, "resultMessage": "SUCCESS"},
		"interface": {
			"id": "iface-001",
			"path": "nas.kr-pub-a.nhncloudservice.com:/vol-001",
			"status": "available",
			"subnetId": "subnet-abc",
			"tenantId": "tenant-xyz"
		}
	}`

	var out nas.CreateInterfaceOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Interface.ID != "iface-001" {
		t.Errorf("Interface.ID: got %q, want %q", out.Interface.ID, "iface-001")
	}
	if out.Interface.SubnetID != "subnet-abc" {
		t.Errorf("Interface.SubnetID: got %q, want %q", out.Interface.SubnetID, "subnet-abc")
	}
}

func TestVolumeMirrorStat_ParseFromJSON(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 0, "resultMessage": "SUCCESS"},
		"volumeMirrorStat": {
			"lastSuccessTransferBytes": 1048576,
			"lastSuccessTransferEndTime": "2023-06-01T04:00:00Z",
			"lastTransferBytes": 2097152,
			"lastTransferEndTime": "2023-06-01T04:05:00Z",
			"lastTransferStatus": "success"
		}
	}`

	var out nas.GetVolumeMirrorStatOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.VolumeMirrorStat.LastTransferStatus != "success" {
		t.Errorf("LastTransferStatus: got %q, want %q", out.VolumeMirrorStat.LastTransferStatus, "success")
	}
	if out.VolumeMirrorStat.LastSuccessTransferBytes != 1048576 {
		t.Errorf("LastSuccessTransferBytes: got %d, want 1048576", out.VolumeMirrorStat.LastSuccessTransferBytes)
	}
}

// ==================== Request Build Tests ====================

func TestCreateVolumeInput_Marshal(t *testing.T) {
	desc := "test nas volume"
	input := nas.CreateVolumeInput{
		Name:        "new-nas",
		Description: &desc,
		SizeGB:      200,
		MountProtocol: &nas.MountProtocol{
			Protocol: "NFS",
		},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["name"] != "new-nas" {
		t.Errorf("name: got %v, want %q", m["name"], "new-nas")
	}
	if m["sizeGb"].(float64) != 200 {
		t.Errorf("sizeGb: got %v, want 200", m["sizeGb"])
	}
	if m["description"] != desc {
		t.Errorf("description: got %v, want %q", m["description"], desc)
	}
}

func TestCreateVolumeInput_OmitsOptionals(t *testing.T) {
	input := nas.CreateVolumeInput{
		Name:   "minimal",
		SizeGB: 100,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := m["description"]; ok {
		t.Error("description should be omitted when nil")
	}
	if _, ok := m["encryption"]; ok {
		t.Error("encryption should be omitted when nil")
	}
	if _, ok := m["snapshotPolicy"]; ok {
		t.Error("snapshotPolicy should be omitted when nil")
	}
}

func TestCreateVolumeMirrorInput_Marshal(t *testing.T) {
	input := nas.CreateVolumeMirrorInput{
		DstRegion:   "jp-east-2",
		DstTenantID: "tenant-remote",
		DstVolume: nas.CreateVolumeMirrorDstInput{
			Name: "mirror-vol",
		},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal map failed: %v", err)
	}

	if m["dstRegion"] != "jp-east-2" {
		t.Errorf("dstRegion: got %v, want %q", m["dstRegion"], "jp-east-2")
	}
	if m["dstTenantId"] != "tenant-remote" {
		t.Errorf("dstTenantId: got %v, want %q", m["dstTenantId"], "tenant-remote")
	}
	dv, ok := m["dstVolume"].(map[string]interface{})
	if !ok {
		t.Fatal("dstVolume should be an object")
	}
	if dv["name"] != "mirror-vol" {
		t.Errorf("dstVolume.name: got %v, want %q", dv["name"], "mirror-vol")
	}
}

func TestListVolumesInput_Marshal(t *testing.T) {
	sizeGB := 100
	input := nas.ListVolumesInput{
		SizeGB:   &sizeGB,
		Name:     "my-vol",
		SubnetID: "subnet-001",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if m["name"] != "my-vol" {
		t.Errorf("name: got %v, want %q", m["name"], "my-vol")
	}
	if m["subnetId"] != "subnet-001" {
		t.Errorf("subnetId: got %v, want %q", m["subnetId"], "subnet-001")
	}
}

// ==================== Constants Tests ====================

func TestVolumeStatusConstants(t *testing.T) {
	if nas.VolumeStatusAvailable != "available" {
		t.Errorf("VolumeStatusAvailable: got %q, want %q", nas.VolumeStatusAvailable, "available")
	}
	if nas.VolumeStatusCreating != "creating" {
		t.Errorf("VolumeStatusCreating: got %q, want %q", nas.VolumeStatusCreating, "creating")
	}
	if nas.VolumeStatusDeleting != "deleting" {
		t.Errorf("VolumeStatusDeleting: got %q, want %q", nas.VolumeStatusDeleting, "deleting")
	}
	if nas.VolumeStatusError != "error" {
		t.Errorf("VolumeStatusError: got %q, want %q", nas.VolumeStatusError, "error")
	}
}

func TestProtocolConstants(t *testing.T) {
	if nas.ProtocolNFS != "NFS" {
		t.Errorf("ProtocolNFS: got %q, want %q", nas.ProtocolNFS, "NFS")
	}
	if nas.ProtocolCIFS != "CIFS" {
		t.Errorf("ProtocolCIFS: got %q, want %q", nas.ProtocolCIFS, "CIFS")
	}
}

// ==================== Type assertion tests ====================

func TestVolume_TimestampTypes(t *testing.T) {
	// Verify time.Time is used for CreatedAt / UpdatedAt
	field, ok := reflect.TypeOf(nas.Volume{}).FieldByName("CreatedAt")
	if !ok {
		t.Fatal("Volume has no field CreatedAt")
	}
	if field.Type != reflect.TypeOf(time.Time{}) {
		t.Errorf("CreatedAt type: got %v, want time.Time", field.Type)
	}

	field2, ok2 := reflect.TypeOf(nas.Volume{}).FieldByName("UpdatedAt")
	if !ok2 {
		t.Fatal("Volume has no field UpdatedAt")
	}
	if field2.Type != reflect.TypeOf(time.Time{}) {
		t.Errorf("UpdatedAt type: got %v, want time.Time", field2.Type)
	}
}
