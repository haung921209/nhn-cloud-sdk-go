package natgateway_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/natgateway"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Required Fields Tests
// ---------------------------------------------------------------------------

func TestNATGateway_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.NATGateway{}), []string{
		"ID",
		"Name",
		"TenantID",
	})
}

func TestListNATGatewaysOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.ListNATGatewaysOutput{}), []string{
		"NATGateways",
	})
}

func TestGetNATGatewayOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.GetNATGatewayOutput{}), []string{
		"NATGateway",
	})
}

func TestCreateNATGatewayInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.CreateNATGatewayInput{}), []string{
		"Name",
		"VPCID",
		"SubnetID",
	})
}

func TestCreateNATGatewayRequest_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.CreateNATGatewayRequest{}), []string{
		"NATGateway",
	})
}

func TestUpdateNATGatewayRequest_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(natgateway.UpdateNATGatewayRequest{}), []string{
		"NATGateway",
	})
}

// ---------------------------------------------------------------------------
// JSON Tag Tests
// ---------------------------------------------------------------------------

func TestNATGateway_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(natgateway.NATGateway{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"Description", "description"},
		{"TenantID", "tenant_id"},
		{"SubnetID", "subnet_id"},
		{"VPCID", "vpc_id"},
		{"FloatingIPID", "floatingips_id"},
		{"FloatingIPAddress", "floating_ip"},
		{"Status", "status"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestListNATGatewaysOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(natgateway.ListNATGatewaysOutput{}), "NATGateways", "natgateways")
}

func TestGetNATGatewayOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(natgateway.GetNATGatewayOutput{}), "NATGateway", "natgateway")
}

func TestCreateNATGatewayInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(natgateway.CreateNATGatewayInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"Name", "name"},
		{"Description", "description"},
		{"VPCID", "vpc_id"},
		{"SubnetID", "subnet_id"},
		// API spec uses "floatingips_id" for the create request body
		{"FloatingIPID", "floatingips_id"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateNATGatewayRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(natgateway.CreateNATGatewayRequest{}), "NATGateway", "natgateway")
}

func TestUpdateNATGatewayInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(natgateway.UpdateNATGatewayInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
}

func TestUpdateNATGatewayRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(natgateway.UpdateNATGatewayRequest{}), "NATGateway", "natgateway")
}

// ---------------------------------------------------------------------------
// Response Parse Tests
// ---------------------------------------------------------------------------

func TestListNATGatewaysOutput_Parse(t *testing.T) {
	raw := `{
		"natgateways": [
			{
				"id": "nat-001",
				"name": "my-nat",
				"description": "test nat gateway",
				"tenant_id": "tenant-001",
				"vpc_id": "vpc-001",
				"subnet_id": "subnet-001",
				"port_id": "port-001",
				"floatingips_id": "fip-001",
				"floating_ip": "203.0.113.10",
				"create_time": "2026-03-25T00:00:00Z",
				"status": "ACTIVE"
			}
		]
	}`
	var out natgateway.ListNATGatewaysOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.NATGateways) != 1 {
		t.Fatalf("expected 1 nat gateway, got %d", len(out.NATGateways))
	}
	gw := out.NATGateways[0]
	if gw.ID != "nat-001" {
		t.Errorf("ID: got %q, want %q", gw.ID, "nat-001")
	}
	if gw.Name != "my-nat" {
		t.Errorf("Name: got %q, want %q", gw.Name, "my-nat")
	}
	if gw.TenantID != "tenant-001" {
		t.Errorf("TenantID: got %q, want %q", gw.TenantID, "tenant-001")
	}
	if gw.VPCID != "vpc-001" {
		t.Errorf("VPCID: got %q, want %q", gw.VPCID, "vpc-001")
	}
	if gw.SubnetID != "subnet-001" {
		t.Errorf("SubnetID: got %q, want %q", gw.SubnetID, "subnet-001")
	}
	if gw.Status != "ACTIVE" {
		t.Errorf("Status: got %q, want %q", gw.Status, "ACTIVE")
	}
	if gw.FloatingIPID != "fip-001" {
		t.Errorf("FloatingIPID: got %q, want %q", gw.FloatingIPID, "fip-001")
	}
	if gw.FloatingIPAddress != "203.0.113.10" {
		t.Errorf("FloatingIPAddress: got %q, want %q", gw.FloatingIPAddress, "203.0.113.10")
	}
}

func TestGetNATGatewayOutput_Parse(t *testing.T) {
	raw := `{
		"natgateway": {
			"id": "nat-002",
			"name": "gw-2",
			"tenant_id": "tenant-002",
			"vpc_id": "vpc-002",
			"subnet_id": "subnet-002",
			"status": "BUILD"
		}
	}`
	var out natgateway.GetNATGatewayOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if out.NATGateway.ID != "nat-002" {
		t.Errorf("ID: got %q, want %q", out.NATGateway.ID, "nat-002")
	}
	if out.NATGateway.Status != "BUILD" {
		t.Errorf("Status: got %q, want %q", out.NATGateway.Status, "BUILD")
	}
	if out.NATGateway.VPCID != "vpc-002" {
		t.Errorf("VPCID: got %q, want %q", out.NATGateway.VPCID, "vpc-002")
	}
}

func TestNATGateway_ParseOptionalFields(t *testing.T) {
	// Verify that optional fields parse correctly when present
	raw := `{
		"natgateway": {
			"id": "nat-003",
			"name": "gw-3",
			"tenant_id": "tenant-003",
			"vpc_id": "vpc-003",
			"subnet_id": "subnet-003",
			"floatingips_id": "fip-003",
			"floating_ip": "203.0.113.30",
			"status": "ACTIVE",
			"create_time": "2026-03-25T12:00:00Z",
			"description": "optional fields test"
		}
	}`
	var out natgateway.GetNATGatewayOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	gw := out.NATGateway
	if gw.FloatingIPID != "fip-003" {
		t.Errorf("FloatingIPID: got %q, want %q", gw.FloatingIPID, "fip-003")
	}
	if gw.FloatingIPAddress != "203.0.113.30" {
		t.Errorf("FloatingIPAddress: got %q, want %q", gw.FloatingIPAddress, "203.0.113.30")
	}
	if gw.Description != "optional fields test" {
		t.Errorf("Description: got %q, want %q", gw.Description, "optional fields test")
	}
}

// ---------------------------------------------------------------------------
// Request Build Tests
// ---------------------------------------------------------------------------

func TestCreateNATGatewayRequest_Build(t *testing.T) {
	req := natgateway.CreateNATGatewayRequest{
		NATGateway: natgateway.CreateNATGatewayInput{
			Name:         "my-nat",
			VPCID:        "vpc-001",
			SubnetID:     "subnet-001",
			FloatingIPID: "fip-001",
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
	gwMap, ok := m["natgateway"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'natgateway' key in output")
	}
	if gwMap["name"] != "my-nat" {
		t.Errorf("name: got %v, want %q", gwMap["name"], "my-nat")
	}
	if gwMap["vpc_id"] != "vpc-001" {
		t.Errorf("vpc_id: got %v, want %q", gwMap["vpc_id"], "vpc-001")
	}
	if gwMap["subnet_id"] != "subnet-001" {
		t.Errorf("subnet_id: got %v, want %q", gwMap["subnet_id"], "subnet-001")
	}
	// API spec requires "floatingips_id" (plural) in create request
	if gwMap["floatingips_id"] != "fip-001" {
		t.Errorf("floatingips_id: got %v, want %q", gwMap["floatingips_id"], "fip-001")
	}
}

func TestCreateNATGatewayRequest_OmitEmptyDescription(t *testing.T) {
	req := natgateway.CreateNATGatewayRequest{
		NATGateway: natgateway.CreateNATGatewayInput{
			Name:     "my-nat",
			VPCID:    "vpc-001",
			SubnetID: "subnet-001",
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
	gwMap, ok := m["natgateway"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'natgateway' key in output")
	}
	if _, exists := gwMap["description"]; exists {
		t.Error("description should be omitted when empty")
	}
	if _, exists := gwMap["floatingips_id"]; exists {
		t.Error("floatingips_id should be omitted when empty")
	}
}

func TestUpdateNATGatewayRequest_Build(t *testing.T) {
	req := natgateway.UpdateNATGatewayRequest{
		NATGateway: natgateway.UpdateNATGatewayInput{
			Name:        "updated-nat",
			Description: "updated description",
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
	gwMap, ok := m["natgateway"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'natgateway' key in output")
	}
	if gwMap["name"] != "updated-nat" {
		t.Errorf("name: got %v, want %q", gwMap["name"], "updated-nat")
	}
	if gwMap["description"] != "updated description" {
		t.Errorf("description: got %v, want %q", gwMap["description"], "updated description")
	}
}

func TestUpdateNATGatewayRequest_OmitEmptyFields(t *testing.T) {
	req := natgateway.UpdateNATGatewayRequest{
		NATGateway: natgateway.UpdateNATGatewayInput{},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	gwMap, ok := m["natgateway"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'natgateway' key in output")
	}
	for _, field := range []string{"name", "description"} {
		if val, exists := gwMap[field]; exists && val != "" {
			t.Errorf("field %q should be omitted when empty, got %v", field, val)
		}
	}
}
