package securitygroup_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Required Fields Tests
// ---------------------------------------------------------------------------

func TestSecurityGroup_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(securitygroup.SecurityGroup{}), []string{
		"ID",
		"Name",
		"TenantID",
		"Rules",
	})
}

func TestSecurityRule_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(securitygroup.SecurityRule{}), []string{
		"ID",
		"TenantID",
		"SecurityGroupID",
		"Direction",
		"EtherType",
	})
}

func TestListSecurityGroupsOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(securitygroup.ListSecurityGroupsOutput{}), []string{
		"SecurityGroups",
	})
}

func TestGetSecurityGroupOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(securitygroup.GetSecurityGroupOutput{}), []string{
		"SecurityGroup",
	})
}

func TestCreateRuleInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(securitygroup.CreateRuleInput{}), []string{
		"SecurityGroupID",
		"Direction",
	})
}

// ---------------------------------------------------------------------------
// JSON Tag Tests
// ---------------------------------------------------------------------------

func TestSecurityGroup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(securitygroup.SecurityGroup{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"Description", "description"},
		{"Rules", "security_group_rules"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestSecurityRule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(securitygroup.SecurityRule{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"TenantID", "tenant_id"},
		{"SecurityGroupID", "security_group_id"},
		{"Direction", "direction"},
		{"EtherType", "ethertype"},
		{"Protocol", "protocol"},
		{"PortRangeMin", "port_range_min"},
		{"PortRangeMax", "port_range_max"},
		{"RemoteIPPrefix", "remote_ip_prefix"},
		{"RemoteGroupID", "remote_group_id"},
		{"Description", "description"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestListSecurityGroupsOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(securitygroup.ListSecurityGroupsOutput{}), "SecurityGroups", "security_groups")
}

func TestGetSecurityGroupOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(securitygroup.GetSecurityGroupOutput{}), "SecurityGroup", "security_group")
}

func TestCreateSecurityGroupInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(securitygroup.CreateSecurityGroupInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
}

func TestCreateRuleInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(securitygroup.CreateRuleInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"SecurityGroupID", "security_group_id"},
		{"Direction", "direction"},
		{"EtherType", "ethertype"},
		{"Protocol", "protocol"},
		{"PortRangeMin", "port_range_min"},
		{"PortRangeMax", "port_range_max"},
		{"RemoteIPPrefix", "remote_ip_prefix"},
		{"RemoteGroupID", "remote_group_id"},
		{"Description", "description"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateRuleOutput_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(securitygroup.CreateRuleOutput{}), "SecurityGroupRule", "security_group_rule")
}

// ---------------------------------------------------------------------------
// Response Parse Tests
// ---------------------------------------------------------------------------

func TestListSecurityGroupsOutput_Parse(t *testing.T) {
	raw := `{
		"security_groups": [
			{
				"id": "sg-001",
				"name": "default",
				"tenant_id": "tenant-001",
				"description": "Default security group",
				"security_group_rules": [
					{
						"id": "rule-001",
						"tenant_id": "tenant-001",
						"security_group_id": "sg-001",
						"direction": "egress",
						"ethertype": "IPv4",
						"protocol": null,
						"port_range_min": null,
						"port_range_max": null
					}
				]
			}
		]
	}`
	var out securitygroup.ListSecurityGroupsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.SecurityGroups) != 1 {
		t.Fatalf("expected 1 security group, got %d", len(out.SecurityGroups))
	}
	sg := out.SecurityGroups[0]
	if sg.ID != "sg-001" {
		t.Errorf("ID: got %q, want %q", sg.ID, "sg-001")
	}
	if sg.Name != "default" {
		t.Errorf("Name: got %q, want %q", sg.Name, "default")
	}
	if sg.TenantID != "tenant-001" {
		t.Errorf("TenantID: got %q, want %q", sg.TenantID, "tenant-001")
	}
	if len(sg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(sg.Rules))
	}
	rule := sg.Rules[0]
	if rule.Direction != "egress" {
		t.Errorf("rule.Direction: got %q, want %q", rule.Direction, "egress")
	}
	if rule.EtherType != "IPv4" {
		t.Errorf("rule.EtherType: got %q, want %q", rule.EtherType, "IPv4")
	}
	if rule.Protocol != nil {
		t.Errorf("rule.Protocol: got %v, want nil", rule.Protocol)
	}
	if rule.PortRangeMin != nil {
		t.Errorf("rule.PortRangeMin: got %v, want nil", rule.PortRangeMin)
	}
}

func TestGetSecurityGroupOutput_Parse(t *testing.T) {
	raw := `{
		"security_group": {
			"id": "sg-002",
			"name": "web-sg",
			"tenant_id": "tenant-002",
			"description": "Web security group",
			"security_group_rules": []
		}
	}`
	var out securitygroup.GetSecurityGroupOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if out.SecurityGroup.ID != "sg-002" {
		t.Errorf("ID: got %q, want %q", out.SecurityGroup.ID, "sg-002")
	}
	if out.SecurityGroup.Name != "web-sg" {
		t.Errorf("Name: got %q, want %q", out.SecurityGroup.Name, "web-sg")
	}
}

func TestCreateRuleOutput_Parse(t *testing.T) {
	portMin := 80
	portMax := 80
	raw := `{
		"security_group_rule": {
			"id": "rule-002",
			"tenant_id": "tenant-001",
			"security_group_id": "sg-001",
			"direction": "ingress",
			"ethertype": "IPv4",
			"protocol": "tcp",
			"port_range_min": 80,
			"port_range_max": 80,
			"remote_ip_prefix": "0.0.0.0/0"
		}
	}`
	var out securitygroup.CreateRuleOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	rule := out.SecurityGroupRule
	if rule.ID != "rule-002" {
		t.Errorf("ID: got %q, want %q", rule.ID, "rule-002")
	}
	if rule.Direction != "ingress" {
		t.Errorf("Direction: got %q, want %q", rule.Direction, "ingress")
	}
	if rule.Protocol == nil || *rule.Protocol != "tcp" {
		t.Errorf("Protocol: got %v, want %q", rule.Protocol, "tcp")
	}
	if rule.PortRangeMin == nil || *rule.PortRangeMin != portMin {
		t.Errorf("PortRangeMin: got %v, want %d", rule.PortRangeMin, portMin)
	}
	if rule.PortRangeMax == nil || *rule.PortRangeMax != portMax {
		t.Errorf("PortRangeMax: got %v, want %d", rule.PortRangeMax, portMax)
	}
	if rule.RemoteIPPrefix != "0.0.0.0/0" {
		t.Errorf("RemoteIPPrefix: got %q, want %q", rule.RemoteIPPrefix, "0.0.0.0/0")
	}
}

// ---------------------------------------------------------------------------
// Request Build Tests
// ---------------------------------------------------------------------------

func TestCreateSecurityGroupInput_Build(t *testing.T) {
	input := securitygroup.CreateSecurityGroupInput{
		Name:        "my-sg",
		Description: "test security group",
	}
	testutil.AssertJSONRoundTrip(t, input, `{"name":"my-sg","description":"test security group"}`)
}

func TestCreateSecurityGroupInput_OmitEmptyDescription(t *testing.T) {
	input := securitygroup.CreateSecurityGroupInput{
		Name: "my-sg",
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if _, exists := m["description"]; exists {
		t.Error("description should be omitted when empty")
	}
}

func TestCreateRuleInput_Build_TCP(t *testing.T) {
	portMin := 443
	portMax := 443
	input := securitygroup.CreateRuleInput{
		SecurityGroupID: "sg-001",
		Direction:       "ingress",
		EtherType:       "IPv4",
		Protocol:        "tcp",
		PortRangeMin:    &portMin,
		PortRangeMax:    &portMax,
		RemoteIPPrefix:  "0.0.0.0/0",
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if m["security_group_id"] != "sg-001" {
		t.Errorf("security_group_id: got %v, want %q", m["security_group_id"], "sg-001")
	}
	if m["direction"] != "ingress" {
		t.Errorf("direction: got %v, want %q", m["direction"], "ingress")
	}
	if m["protocol"] != "tcp" {
		t.Errorf("protocol: got %v, want %q", m["protocol"], "tcp")
	}
	if m["port_range_min"] != float64(443) {
		t.Errorf("port_range_min: got %v, want 443", m["port_range_min"])
	}
	if m["port_range_max"] != float64(443) {
		t.Errorf("port_range_max: got %v, want 443", m["port_range_max"])
	}
}

func TestCreateRuleInput_OmitEmptyFields(t *testing.T) {
	input := securitygroup.CreateRuleInput{
		SecurityGroupID: "sg-001",
		Direction:       "egress",
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	for _, field := range []string{"ethertype", "protocol", "port_range_min", "port_range_max", "remote_ip_prefix", "remote_group_id", "description"} {
		if _, exists := m[field]; exists {
			t.Errorf("field %q should be omitted when empty", field)
		}
	}
}

func TestUpdateSecurityGroupInput_Build(t *testing.T) {
	input := securitygroup.UpdateSecurityGroupInput{
		Name:        "updated-sg",
		Description: "updated description",
	}
	testutil.AssertJSONRoundTrip(t, input, `{"name":"updated-sg","description":"updated description"}`)
}
