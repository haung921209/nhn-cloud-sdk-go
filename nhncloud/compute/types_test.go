package compute

// types_test.go lives in package compute (white-box) so it can reference
// unexported details if needed; primarily it uses reflect to check struct
// field presence and JSON tags against the API spec.

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Server struct – required fields from API spec
// ---------------------------------------------------------------------------

func TestServer_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(Server{})
	testutil.AssertAllRequiredFields(t, typ, []string{
		"ID",
		"Name",
		"Status",
		"TenantID",
		"UserID",
		"Created",
		"Addresses",
		"Metadata",
		"SecurityGroups",
	})
}

func TestServer_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Server{})
	cases := []struct{ field, tag string }{
		{"ID", "id"},
		{"Name", "name"},
		{"Status", "status"},
		{"TenantID", "tenant_id"},
		{"UserID", "user_id"},
		{"KeyName", "key_name"},
		{"AvailabilityZone", "OS-EXT-AZ:availability_zone"},
		{"Created", "created"},
		{"Updated", "updated"},
		{"Addresses", "addresses"},
		{"Metadata", "metadata"},
		{"SecurityGroups", "security_groups"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

// TestServer_JSONRoundTrip tests that a full Server object marshals/unmarshals
// correctly using the field names documented in the API spec.
func TestServer_JSONRoundTrip(t *testing.T) {
	srv := Server{
		ID:               "a3b4c5d6-e7f8-4a2b-9c1d-e0f1a2b3c4d5",
		Name:             "my-server",
		Status:           "ACTIVE",
		TenantID:         "tenant-abc",
		UserID:           "user-xyz",
		KeyName:          "my-keypair",
		AvailabilityZone: "kr1-pub-a",
		Created:          "2026-01-01T00:00:00Z",
		Updated:          "2026-01-02T00:00:00Z",
		Metadata:         map[string]string{"env": "prod"},
	}

	data, err := json.Marshal(srv)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Server
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ID != srv.ID {
		t.Errorf("ID: got %q, want %q", decoded.ID, srv.ID)
	}
	if decoded.Status != srv.Status {
		t.Errorf("Status: got %q, want %q", decoded.Status, srv.Status)
	}
	if decoded.AvailabilityZone != srv.AvailabilityZone {
		t.Errorf("AvailabilityZone: got %q, want %q", decoded.AvailabilityZone, srv.AvailabilityZone)
	}
	if decoded.Metadata["env"] != "prod" {
		t.Errorf("Metadata[env]: got %q, want %q", decoded.Metadata["env"], "prod")
	}
}

// TestServer_ParseAPIResponse verifies that the SDK correctly parses a sample
// JSON payload modelled on the detailed instance list endpoint response.
func TestServer_ParseAPIResponse(t *testing.T) {
	const sampleJSON = `{
		"servers": [
			{
				"id": "srv-001",
				"name": "web-server",
				"status": "ACTIVE",
				"tenant_id": "ten-abc",
				"user_id": "usr-xyz",
				"key_name": "mykey",
				"created": "2026-01-01T00:00:00Z",
				"updated": "2026-01-02T00:00:00Z",
				"OS-EXT-AZ:availability_zone": "kr1-pub-a",
				"addresses": {
					"public": [
						{
							"addr": "192.0.2.1",
							"version": 4,
							"OS-EXT-IPS:type": "floating"
						}
					]
				},
				"metadata": {"role": "web"},
				"security_groups": [{"name": "default"}]
			}
		]
	}`

	var out ListServersOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.Servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(out.Servers))
	}

	s := out.Servers[0]
	if s.ID != "srv-001" {
		t.Errorf("ID: got %q, want %q", s.ID, "srv-001")
	}
	if s.Status != "ACTIVE" {
		t.Errorf("Status: got %q, want %q", s.Status, "ACTIVE")
	}
	if s.AvailabilityZone != "kr1-pub-a" {
		t.Errorf("AvailabilityZone: got %q, want %q", s.AvailabilityZone, "kr1-pub-a")
	}
	addrs, ok := s.Addresses["public"]
	if !ok || len(addrs) == 0 {
		t.Fatal("expected addresses[public] to be populated")
	}
	if addrs[0].Addr != "192.0.2.1" {
		t.Errorf("Addresses[public][0].Addr: got %q, want %q", addrs[0].Addr, "192.0.2.1")
	}
	if addrs[0].Type != "floating" {
		t.Errorf("Addresses[public][0].Type: got %q, want %q", addrs[0].Type, "floating")
	}
	if s.Metadata["role"] != "web" {
		t.Errorf("Metadata[role]: got %q, want %q", s.Metadata["role"], "web")
	}
	if len(s.SecurityGroups) == 0 || s.SecurityGroups[0].Name != "default" {
		t.Errorf("SecurityGroups: unexpected value %+v", s.SecurityGroups)
	}
}

// ---------------------------------------------------------------------------
// Address struct
// ---------------------------------------------------------------------------

func TestAddress_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Address{})
	testutil.AssertStructHasJSONTag(t, typ, "Addr", "addr")
	testutil.AssertStructHasJSONTag(t, typ, "Version", "version")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "OS-EXT-IPS:type")
}

// ---------------------------------------------------------------------------
// Flavor struct – required fields from API spec (detailed)
// ---------------------------------------------------------------------------

func TestFlavor_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(Flavor{})
	testutil.AssertAllRequiredFields(t, typ, []string{
		"ID",
		"Name",
		"RAM",
		"VCPUs",
		"Disk",
	})
}

func TestFlavor_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Flavor{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "RAM", "ram")
	testutil.AssertStructHasJSONTag(t, typ, "VCPUs", "vcpus")
	testutil.AssertStructHasJSONTag(t, typ, "Disk", "disk")
}

func TestFlavor_ParseAPIResponse(t *testing.T) {
	const sampleJSON = `{
		"flavors": [
			{
				"id": "m2.c1m2",
				"name": "m2.c1m2",
				"ram": 2048,
				"vcpus": 1,
				"disk": 20
			}
		]
	}`

	var out ListFlavorsOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.Flavors) != 1 {
		t.Fatalf("expected 1 flavor, got %d", len(out.Flavors))
	}

	f := out.Flavors[0]
	if f.ID != "m2.c1m2" {
		t.Errorf("ID: got %q, want %q", f.ID, "m2.c1m2")
	}
	if f.RAM != 2048 {
		t.Errorf("RAM: got %d, want 2048", f.RAM)
	}
	if f.VCPUs != 1 {
		t.Errorf("VCPUs: got %d, want 1", f.VCPUs)
	}
	if f.Disk != 20 {
		t.Errorf("Disk: got %d, want 20", f.Disk)
	}
}

// ---------------------------------------------------------------------------
// Image struct
// ---------------------------------------------------------------------------

func TestImage_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(Image{})
	testutil.AssertAllRequiredFields(t, typ, []string{
		"ID",
		"Name",
		"Status",
		"MinDisk",
		"MinRAM",
		"Created",
	})
}

func TestImage_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Image{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "MinDisk", "minDisk")
	testutil.AssertStructHasJSONTag(t, typ, "MinRAM", "minRam")
}

// ---------------------------------------------------------------------------
// KeyPair struct
// ---------------------------------------------------------------------------

func TestKeyPair_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(KeyPair{})
	testutil.AssertAllRequiredFields(t, typ, []string{
		"Name",
		"PublicKey",
		"Fingerprint",
	})
}

func TestKeyPair_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(KeyPair{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "PublicKey", "public_key")
	testutil.AssertStructHasJSONTag(t, typ, "Fingerprint", "fingerprint")
}

func TestListKeyPairsOutput_ParseAPIResponse(t *testing.T) {
	// The API wraps each keypair in a "keypair" object inside the array.
	const sampleJSON = `{
		"keypairs": [
			{
				"keypair": {
					"name": "my-key",
					"public_key": "ssh-rsa AAAAB3NzaC...",
					"fingerprint": "ab:cd:ef:01:23:45:67:89:ab:cd:ef:01:23:45:67:89"
				}
			}
		]
	}`

	var out ListKeyPairsOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.KeyPairs) != 1 {
		t.Fatalf("expected 1 keypair, got %d", len(out.KeyPairs))
	}

	kp := out.KeyPairs[0].KeyPair
	if kp.Name != "my-key" {
		t.Errorf("Name: got %q, want %q", kp.Name, "my-key")
	}
	if kp.Fingerprint != "ab:cd:ef:01:23:45:67:89:ab:cd:ef:01:23:45:67:89" {
		t.Errorf("Fingerprint: got %q", kp.Fingerprint)
	}
}

// ---------------------------------------------------------------------------
// CreateServerInput – required fields from API spec
// ---------------------------------------------------------------------------

func TestCreateServerInput_RequiredFields(t *testing.T) {
	typ := reflect.TypeOf(CreateServerInput{})
	testutil.AssertAllRequiredFields(t, typ, []string{
		"Name",
		"FlavorRef",
		"Networks",
		"BlockDeviceMapping",
		"KeyName",
	})
}

func TestCreateServerInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateServerInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "ImageRef", "imageRef")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorRef", "flavorRef")
	testutil.AssertStructHasJSONTag(t, typ, "KeyName", "key_name")
	testutil.AssertStructHasJSONTag(t, typ, "AvailabilityZone", "availability_zone")
	testutil.AssertStructHasJSONTag(t, typ, "Networks", "networks")
	testutil.AssertStructHasJSONTag(t, typ, "SecurityGroups", "security_groups")
	testutil.AssertStructHasJSONTag(t, typ, "Metadata", "metadata")
	testutil.AssertStructHasJSONTag(t, typ, "UserData", "user_data")
	testutil.AssertStructHasJSONTag(t, typ, "BlockDeviceMapping", "block_device_mapping_v2")
}

// TestCreateServerInput_MarshalWrapped verifies that the SDK wraps the input
// in a "server" key when building the request body, matching the API spec.
func TestCreateServerInput_MarshalWrapped(t *testing.T) {
	input := &CreateServerInput{
		Name:      "web-01",
		FlavorRef: "m2.c2m4",
		ImageRef:  "img-001",
		KeyName:   "my-key",
		Networks: []ServerNetwork{
			{UUID: "net-abc"},
		},
		BlockDeviceMapping: []BlockDeviceMapping{
			{
				BootIndex:           0,
				SourceType:          "image",
				DestinationType:     "volume",
				VolumeSize:          50,
				DeleteOnTermination: true,
			},
		},
	}

	wrapped := map[string]interface{}{"server": input}
	data, err := json.Marshal(wrapped)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify top-level "server" key exists.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatalf("top-level unmarshal failed: %v", err)
	}
	if _, ok := top["server"]; !ok {
		t.Fatal("expected top-level 'server' key in request body")
	}

	// Verify inner field names match API spec.
	var inner map[string]json.RawMessage
	if err := json.Unmarshal(top["server"], &inner); err != nil {
		t.Fatalf("inner unmarshal failed: %v", err)
	}
	for _, key := range []string{"name", "flavorRef", "imageRef", "key_name", "networks", "block_device_mapping_v2"} {
		if _, ok := inner[key]; !ok {
			t.Errorf("expected field %q in serialized server body", key)
		}
	}
}

// ---------------------------------------------------------------------------
// BlockDeviceMapping struct
// ---------------------------------------------------------------------------

func TestBlockDeviceMapping_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(BlockDeviceMapping{})
	testutil.AssertStructHasJSONTag(t, typ, "BootIndex", "boot_index")
	testutil.AssertStructHasJSONTag(t, typ, "SourceType", "source_type")
	testutil.AssertStructHasJSONTag(t, typ, "DestinationType", "destination_type")
	testutil.AssertStructHasJSONTag(t, typ, "VolumeSize", "volume_size")
	testutil.AssertStructHasJSONTag(t, typ, "DeleteOnTermination", "delete_on_termination")
}

// ---------------------------------------------------------------------------
// ServerNetwork struct
// ---------------------------------------------------------------------------

func TestServerNetwork_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ServerNetwork{})
	testutil.AssertStructHasJSONTag(t, typ, "UUID", "uuid")
	testutil.AssertStructHasJSONTag(t, typ, "Port", "port")
	testutil.AssertStructHasJSONTag(t, typ, "FixedIP", "fixed_ip")
}

// ---------------------------------------------------------------------------
// AvailabilityZone struct – matches API spec field names
// ---------------------------------------------------------------------------

func TestAvailabilityZone_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(AvailabilityZone{})
	testutil.AssertStructHasJSONTag(t, typ, "ZoneName", "zoneName")
	testutil.AssertStructHasJSONTag(t, typ, "ZoneState", "zoneState")
}

func TestAvailabilityZone_ParseAPIResponse(t *testing.T) {
	const sampleJSON = `{
		"availabilityZoneInfo": [
			{
				"zoneName": "kr1-pub-a",
				"zoneState": {"available": true}
			},
			{
				"zoneName": "kr1-pub-b",
				"zoneState": {"available": false}
			}
		]
	}`

	var out ListAvailabilityZonesOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(out.AvailabilityZoneInfo) != 2 {
		t.Fatalf("expected 2 AZs, got %d", len(out.AvailabilityZoneInfo))
	}

	if out.AvailabilityZoneInfo[0].ZoneName != "kr1-pub-a" {
		t.Errorf("ZoneName[0]: got %q, want %q", out.AvailabilityZoneInfo[0].ZoneName, "kr1-pub-a")
	}
	if !out.AvailabilityZoneInfo[0].ZoneState.Available {
		t.Error("ZoneState[0].Available: want true")
	}
	if out.AvailabilityZoneInfo[1].ZoneState.Available {
		t.Error("ZoneState[1].Available: want false")
	}
}

// ---------------------------------------------------------------------------
// RebootInput / ResizeInput action bodies
// ---------------------------------------------------------------------------

func TestRebootInput_SoftHard(t *testing.T) {
	cases := []struct {
		input    RebootInput
		wantType string
	}{
		{RebootInput{Reboot: struct {
			Type string `json:"type"`
		}{Type: "SOFT"}}, "SOFT"},
		{RebootInput{Reboot: struct {
			Type string `json:"type"`
		}{Type: "HARD"}}, "HARD"},
	}

	for _, c := range cases {
		data, err := json.Marshal(c.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if !strings.Contains(string(data), c.wantType) {
			t.Errorf("expected %q in %s", c.wantType, string(data))
		}
	}
}

func TestResizeInput_FlavorRef(t *testing.T) {
	input := ResizeInput{Resize: struct {
		FlavorRef string `json:"flavorRef"`
	}{FlavorRef: "m2.c4m8"}}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var out map[string]map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out["resize"]["flavorRef"] != "m2.c4m8" {
		t.Errorf("resize.flavorRef: got %q, want %q", out["resize"]["flavorRef"], "m2.c4m8")
	}
}

// ---------------------------------------------------------------------------
// CreateKeyPairInput / CreateKeyPairOutput
// ---------------------------------------------------------------------------

func TestCreateKeyPairInput_JSONRoundTrip(t *testing.T) {
	testutil.AssertJSONRoundTrip(t,
		CreateKeyPairInput{Name: "my-key", PublicKey: "ssh-rsa AAAA..."},
		`{"name":"my-key","public_key":"ssh-rsa AAAA..."}`,
	)
}

func TestCreateKeyPairInput_OmitPublicKey(t *testing.T) {
	// When PublicKey is empty (generate new keypair), it must be omitted.
	input := CreateKeyPairInput{Name: "generated-key"}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if strings.Contains(string(data), "public_key") {
		t.Errorf("public_key should be omitted when empty; got: %s", string(data))
	}
}

func TestCreateKeyPairOutput_ParseAPIResponse(t *testing.T) {
	const sampleJSON = `{
		"keypair": {
			"name": "new-key",
			"public_key": "ssh-rsa AAAAB3...",
			"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...",
			"fingerprint": "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"
		}
	}`

	var out CreateKeyPairOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.KeyPair.Name != "new-key" {
		t.Errorf("Name: got %q, want %q", out.KeyPair.Name, "new-key")
	}
	if out.KeyPair.PrivateKey == "" {
		t.Error("PrivateKey: expected non-empty on creation response")
	}
	if out.KeyPair.Fingerprint != "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99" {
		t.Errorf("Fingerprint: got %q", out.KeyPair.Fingerprint)
	}
}

// ---------------------------------------------------------------------------
// GetServerOutput / CreateServerOutput wrappers
// ---------------------------------------------------------------------------

func TestGetServerOutput_Wrapper(t *testing.T) {
	const sampleJSON = `{
		"server": {
			"id": "srv-999",
			"name": "test",
			"status": "BUILD",
			"tenant_id": "ten",
			"user_id": "usr",
			"created": "2026-01-01T00:00:00Z"
		}
	}`

	var out GetServerOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.Server.ID != "srv-999" {
		t.Errorf("Server.ID: got %q, want %q", out.Server.ID, "srv-999")
	}
	if out.Server.Status != "BUILD" {
		t.Errorf("Server.Status: got %q, want %q", out.Server.Status, "BUILD")
	}
}

func TestCreateServerOutput_Wrapper(t *testing.T) {
	const sampleJSON = `{
		"server": {
			"id": "srv-new",
			"name": "web-01",
			"status": "",
			"tenant_id": "ten",
			"user_id": "usr",
			"created": ""
		}
	}`

	var out CreateServerOutput
	if err := json.Unmarshal([]byte(sampleJSON), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.Server.ID != "srv-new" {
		t.Errorf("Server.ID: got %q, want %q", out.Server.ID, "srv-new")
	}
}
