package floatingip_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/floatingip"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Required Fields Tests
// ---------------------------------------------------------------------------

func TestFloatingIP_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.FloatingIP{}), []string{
		"ID",
		"FloatingNetworkID",
		"FloatingIPAddress",
		"TenantID",
		"Status",
	})
}

func TestCreateFloatingIPInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.CreateFloatingIPInput{}), []string{
		"FloatingNetworkID",
	})
}

func TestCreateFloatingIPRequest_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.CreateFloatingIPRequest{}), []string{
		"FloatingIP",
	})
}

func TestListFloatingIPsOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.ListFloatingIPsOutput{}), []string{
		"FloatingIPs",
	})
}

func TestGetFloatingIPOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.GetFloatingIPOutput{}), []string{
		"FloatingIP",
	})
}

func TestUpdateFloatingIPInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(floatingip.UpdateFloatingIPInput{}), []string{
		"PortID",
	})
}

// ---------------------------------------------------------------------------
// JSON Tag Tests
// ---------------------------------------------------------------------------

func TestFloatingIP_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(floatingip.FloatingIP{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"FloatingNetworkID", "floating_network_id"},
		{"FloatingIPAddress", "floating_ip_address"},
		{"FixedIPAddress", "fixed_ip_address"},
		{"PortID", "port_id"},
		{"TenantID", "tenant_id"},
		{"Status", "status"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateFloatingIPInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(floatingip.CreateFloatingIPInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"FloatingNetworkID", "floating_network_id"},
		{"PortID", "port_id"},
		{"FixedIPAddress", "fixed_ip_address"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateFloatingIPRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(floatingip.CreateFloatingIPRequest{}), "FloatingIP", "floatingip")
}

func TestUpdateFloatingIPInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(floatingip.UpdateFloatingIPInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"PortID", "port_id"},
		{"FixedIPAddress", "fixed_ip_address"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestUpdateFloatingIPRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(floatingip.UpdateFloatingIPRequest{}), "FloatingIP", "floatingip")
}

func TestListFloatingIPsOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(floatingip.ListFloatingIPsOutput{}), "FloatingIPs", "floatingips")
}

func TestGetFloatingIPOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(floatingip.GetFloatingIPOutput{}), "FloatingIP", "floatingip")
}

// ---------------------------------------------------------------------------
// Response Parse Tests
// ---------------------------------------------------------------------------

func TestListFloatingIPsOutput_Parse(t *testing.T) {
	raw := `{
		"floatingips": [
			{
				"id": "fip-001",
				"floating_network_id": "net-ext-001",
				"floating_ip_address": "203.0.113.10",
				"fixed_ip_address": "192.168.1.5",
				"port_id": "port-abc",
				"tenant_id": "tenant-xyz",
				"status": "ACTIVE"
			}
		]
	}`
	var out floatingip.ListFloatingIPsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.FloatingIPs) != 1 {
		t.Fatalf("expected 1 floatingip, got %d", len(out.FloatingIPs))
	}
	fip := out.FloatingIPs[0]
	if fip.ID != "fip-001" {
		t.Errorf("ID: got %q, want %q", fip.ID, "fip-001")
	}
	if fip.FloatingNetworkID != "net-ext-001" {
		t.Errorf("FloatingNetworkID: got %q, want %q", fip.FloatingNetworkID, "net-ext-001")
	}
	if fip.FloatingIPAddress != "203.0.113.10" {
		t.Errorf("FloatingIPAddress: got %q, want %q", fip.FloatingIPAddress, "203.0.113.10")
	}
	if fip.FixedIPAddress != "192.168.1.5" {
		t.Errorf("FixedIPAddress: got %q, want %q", fip.FixedIPAddress, "192.168.1.5")
	}
	if fip.TenantID != "tenant-xyz" {
		t.Errorf("TenantID: got %q, want %q", fip.TenantID, "tenant-xyz")
	}
	if fip.Status != "ACTIVE" {
		t.Errorf("Status: got %q, want %q", fip.Status, "ACTIVE")
	}
}

func TestGetFloatingIPOutput_Parse(t *testing.T) {
	raw := `{
		"floatingip": {
			"id": "fip-002",
			"floating_network_id": "net-ext-002",
			"floating_ip_address": "203.0.113.20",
			"fixed_ip_address": "",
			"tenant_id": "tenant-abc",
			"status": "DOWN"
		}
	}`
	var out floatingip.GetFloatingIPOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if out.FloatingIP.ID != "fip-002" {
		t.Errorf("ID: got %q, want %q", out.FloatingIP.ID, "fip-002")
	}
	if out.FloatingIP.Status != "DOWN" {
		t.Errorf("Status: got %q, want %q", out.FloatingIP.Status, "DOWN")
	}
}

// ---------------------------------------------------------------------------
// Request Build Tests
// ---------------------------------------------------------------------------

func TestCreateFloatingIPRequest_Build(t *testing.T) {
	req := floatingip.CreateFloatingIPRequest{
		FloatingIP: floatingip.CreateFloatingIPInput{
			FloatingNetworkID: "net-ext-001",
			PortID:            "port-abc",
		},
	}
	testutil.AssertJSONRoundTrip(t, req, `{
		"floatingip": {
			"floating_network_id": "net-ext-001",
			"port_id": "port-abc"
		}
	}`)
}

func TestCreateFloatingIPRequest_MinimalBuild(t *testing.T) {
	req := floatingip.CreateFloatingIPRequest{
		FloatingIP: floatingip.CreateFloatingIPInput{
			FloatingNetworkID: "net-ext-001",
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	fipMap, ok := m["floatingip"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'floatingip' key in output")
	}
	if fipMap["floating_network_id"] != "net-ext-001" {
		t.Errorf("floating_network_id: got %v, want %q", fipMap["floating_network_id"], "net-ext-001")
	}
	// port_id should be omitted when empty
	if _, exists := fipMap["port_id"]; exists {
		t.Error("port_id should be omitted when empty")
	}
}

func TestUpdateFloatingIPRequest_Build_Detach(t *testing.T) {
	// nil port_id signals detach
	req := floatingip.UpdateFloatingIPRequest{
		FloatingIP: floatingip.UpdateFloatingIPInput{
			PortID: nil,
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	fipMap, ok := m["floatingip"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'floatingip' key in output")
	}
	// port_id must be present and null for detach
	if _, exists := fipMap["port_id"]; !exists {
		t.Error("port_id must be present (as null) for detach operation")
	}
}

func TestAssociateFloatingIPInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(floatingip.AssociateFloatingIPInput{})
	testutil.AssertStructHasJSONTag(t, typ, "PortID", "port_id")
	testutil.AssertStructHasJSONTag(t, typ, "FixedIPAddress", "fixed_ip_address")
}

func TestDisassociateFloatingIPInput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(floatingip.DisassociateFloatingIPInput{}), "PortID", "port_id")
}
