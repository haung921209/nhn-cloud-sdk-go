package loadbalancer_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Required Fields Tests
// ---------------------------------------------------------------------------

func TestLoadBalancer_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.LoadBalancer{}), []string{
		"ID",
		"Name",
		"TenantID",
		"VIPAddress",
		"VIPPortID",
		"VIPSubnetID",
		"ProvisioningStatus",
		"OperatingStatus",
		"AdminStateUp",
		"Provider",
	})
}

func TestListener_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.Listener{}), []string{
		"ID",
		"Name",
		"TenantID",
		"Protocol",
		"ProtocolPort",
		"AdminStateUp",
	})
}

func TestPool_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.Pool{}), []string{
		"ID",
		"Name",
		"TenantID",
		"Protocol",
		"LBAlgorithm",
		"AdminStateUp",
	})
}

func TestMember_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.Member{}), []string{
		"ID",
		"Name",
		"TenantID",
		"Address",
		"ProtocolPort",
		"Weight",
		"SubnetID",
		"AdminStateUp",
	})
}

func TestHealthMonitor_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.HealthMonitor{}), []string{
		"ID",
		"Name",
		"TenantID",
		"PoolID",
		"Type",
		"Delay",
		"Timeout",
		"MaxRetries",
		"AdminStateUp",
	})
}

func TestCreateLoadBalancerInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.CreateLoadBalancerInput{}), []string{
		"VIPSubnetID",
	})
}

func TestCreateListenerInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.CreateListenerInput{}), []string{
		"LoadBalancerID",
		"Protocol",
		"ProtocolPort",
	})
}

func TestCreateMemberInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.CreateMemberInput{}), []string{
		"Address",
		"ProtocolPort",
	})
}

func TestCreateHealthMonitorInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(loadbalancer.CreateHealthMonitorInput{}), []string{
		"PoolID",
		"Type",
		"Delay",
		"Timeout",
		"MaxRetries",
	})
}

// ---------------------------------------------------------------------------
// JSON Tag Tests
// ---------------------------------------------------------------------------

func TestLoadBalancer_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.LoadBalancer{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"Description", "description"},
		{"TenantID", "tenant_id"},
		{"VIPAddress", "vip_address"},
		{"VIPPortID", "vip_port_id"},
		{"VIPSubnetID", "vip_subnet_id"},
		{"ProvisioningStatus", "provisioning_status"},
		{"OperatingStatus", "operating_status"},
		{"AdminStateUp", "admin_state_up"},
		{"Provider", "provider"},
		{"Listeners", "listeners"},
		{"Pools", "pools"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestListener_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.Listener{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"Protocol", "protocol"},
		{"ProtocolPort", "protocol_port"},
		{"DefaultPoolID", "default_pool_id"},
		{"ConnectionLimit", "connection_limit"},
		{"AdminStateUp", "admin_state_up"},
		{"ProvisioningStatus", "provisioning_status"},
		{"OperatingStatus", "operating_status"},
		{"DefaultTLSContainerRef", "default_tls_container_ref"},
		{"SNIContainerRefs", "sni_container_refs"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestPool_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.Pool{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"Protocol", "protocol"},
		{"LBAlgorithm", "lb_algorithm"},
		{"AdminStateUp", "admin_state_up"},
		{"ProvisioningStatus", "provisioning_status"},
		{"HealthMonitorID", "healthmonitor_id"},
		{"SessionPersistence", "session_persistence"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestMember_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.Member{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"Address", "address"},
		{"ProtocolPort", "protocol_port"},
		{"Weight", "weight"},
		{"SubnetID", "subnet_id"},
		{"AdminStateUp", "admin_state_up"},
		{"ProvisioningStatus", "provisioning_status"},
		{"OperatingStatus", "operating_status"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestHealthMonitor_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.HealthMonitor{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"PoolID", "pool_id"},
		{"Type", "type"},
		{"Delay", "delay"},
		{"Timeout", "timeout"},
		{"MaxRetries", "max_retries"},
		{"MaxRetriesDown", "max_retries_down"},
		{"HTTPMethod", "http_method"},
		{"URLPath", "url_path"},
		{"ExpectedCodes", "expected_codes"},
		{"AdminStateUp", "admin_state_up"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateLoadBalancerInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.CreateLoadBalancerInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"Name", "name"},
		{"Description", "description"},
		{"VIPSubnetID", "vip_subnet_id"},
		{"VIPAddress", "vip_address"},
		{"AdminStateUp", "admin_state_up"},
		{"Provider", "provider"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateLoadBalancerRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(loadbalancer.CreateLoadBalancerRequest{}), "LoadBalancer", "loadbalancer")
}

func TestCreateListenerInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.CreateListenerInput{})
	cases := []struct {
		field string
		tag   string
	}{
		{"LoadBalancerID", "loadbalancer_id"},
		{"Protocol", "protocol"},
		{"ProtocolPort", "protocol_port"},
		{"DefaultPoolID", "default_pool_id"},
		{"ConnectionLimit", "connection_limit"},
		{"AdminStateUp", "admin_state_up"},
		{"DefaultTLSContainerRef", "default_tls_container_ref"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestCreateListenerRequest_JSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(loadbalancer.CreateListenerRequest{}), "Listener", "listener")
}

func TestSessionPersistence_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.SessionPersistence{})
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "CookieName", "cookie_name")
}

func TestL7Policy_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.L7Policy{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"Name", "name"},
		{"TenantID", "tenant_id"},
		{"ListenerID", "listener_id"},
		{"Action", "action"},
		{"Position", "position"},
		{"RedirectPoolID", "redirect_pool_id"},
		{"RedirectURL", "redirect_url"},
		{"AdminStateUp", "admin_state_up"},
		{"ProvisioningStatus", "provisioning_status"},
		{"Rules", "rules"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

func TestL7Rule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(loadbalancer.L7Rule{})
	cases := []struct {
		field string
		tag   string
	}{
		{"ID", "id"},
		{"TenantID", "tenant_id"},
		{"Type", "type"},
		{"CompareType", "compare_type"},
		{"Value", "value"},
		{"Invert", "invert"},
		{"AdminStateUp", "admin_state_up"},
	}
	for _, c := range cases {
		testutil.AssertStructHasJSONTag(t, typ, c.field, c.tag)
	}
}

// ---------------------------------------------------------------------------
// Response Parse Tests
// ---------------------------------------------------------------------------

func TestListLoadBalancersOutput_Parse(t *testing.T) {
	raw := `{
		"loadbalancers": [
			{
				"id": "lb-001",
				"name": "my-lb",
				"description": "test load balancer",
				"tenant_id": "tenant-001",
				"vip_address": "192.168.1.100",
				"vip_port_id": "port-001",
				"vip_subnet_id": "subnet-001",
				"vip_network_id": "net-001",
				"provisioning_status": "ACTIVE",
				"operating_status": "ONLINE",
				"admin_state_up": true,
				"provider": "haproxy",
				"listeners": [{"id": "listener-001"}],
				"pools": [{"id": "pool-001"}]
			}
		]
	}`
	var out loadbalancer.ListLoadBalancersOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.LoadBalancers) != 1 {
		t.Fatalf("expected 1 load balancer, got %d", len(out.LoadBalancers))
	}
	lb := out.LoadBalancers[0]
	if lb.ID != "lb-001" {
		t.Errorf("ID: got %q, want %q", lb.ID, "lb-001")
	}
	if lb.Name != "my-lb" {
		t.Errorf("Name: got %q, want %q", lb.Name, "my-lb")
	}
	if lb.VIPAddress != "192.168.1.100" {
		t.Errorf("VIPAddress: got %q, want %q", lb.VIPAddress, "192.168.1.100")
	}
	if lb.ProvisioningStatus != "ACTIVE" {
		t.Errorf("ProvisioningStatus: got %q, want %q", lb.ProvisioningStatus, "ACTIVE")
	}
	if lb.OperatingStatus != "ONLINE" {
		t.Errorf("OperatingStatus: got %q, want %q", lb.OperatingStatus, "ONLINE")
	}
	if !lb.AdminStateUp {
		t.Error("AdminStateUp: got false, want true")
	}
	if len(lb.Listeners) != 1 || lb.Listeners[0].ID != "listener-001" {
		t.Errorf("Listeners: got %v, want [{listener-001}]", lb.Listeners)
	}
	if len(lb.Pools) != 1 || lb.Pools[0].ID != "pool-001" {
		t.Errorf("Pools: got %v, want [{pool-001}]", lb.Pools)
	}
}

func TestListListenersOutput_Parse(t *testing.T) {
	raw := `{
		"listeners": [
			{
				"id": "listener-001",
				"name": "http-listener",
				"tenant_id": "tenant-001",
				"loadbalancer_id": "lb-001",
				"protocol": "HTTP",
				"protocol_port": 80,
				"default_pool_id": "pool-001",
				"connection_limit": 2000,
				"admin_state_up": true,
				"provisioning_status": "ACTIVE",
				"operating_status": "ONLINE"
			}
		]
	}`
	var out loadbalancer.ListListenersOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.Listeners) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(out.Listeners))
	}
	l := out.Listeners[0]
	if l.ID != "listener-001" {
		t.Errorf("ID: got %q, want %q", l.ID, "listener-001")
	}
	if l.Protocol != "HTTP" {
		t.Errorf("Protocol: got %q, want %q", l.Protocol, "HTTP")
	}
	if l.ProtocolPort != 80 {
		t.Errorf("ProtocolPort: got %d, want 80", l.ProtocolPort)
	}
	if l.ConnectionLimit != 2000 {
		t.Errorf("ConnectionLimit: got %d, want 2000", l.ConnectionLimit)
	}
}

func TestListPoolsOutput_Parse(t *testing.T) {
	raw := `{
		"pools": [
			{
				"id": "pool-001",
				"name": "my-pool",
				"tenant_id": "tenant-001",
				"protocol": "HTTP",
				"lb_algorithm": "ROUND_ROBIN",
				"admin_state_up": true,
				"provisioning_status": "ACTIVE",
				"operating_status": "ONLINE",
				"healthmonitor_id": "hm-001",
				"session_persistence": {
					"type": "HTTP_COOKIE",
					"cookie_name": "SESSION_ID"
				}
			}
		]
	}`
	var out loadbalancer.ListPoolsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.Pools) != 1 {
		t.Fatalf("expected 1 pool, got %d", len(out.Pools))
	}
	p := out.Pools[0]
	if p.ID != "pool-001" {
		t.Errorf("ID: got %q, want %q", p.ID, "pool-001")
	}
	if p.LBAlgorithm != "ROUND_ROBIN" {
		t.Errorf("LBAlgorithm: got %q, want %q", p.LBAlgorithm, "ROUND_ROBIN")
	}
	if p.HealthMonitorID != "hm-001" {
		t.Errorf("HealthMonitorID: got %q, want %q", p.HealthMonitorID, "hm-001")
	}
	if p.SessionPersistence == nil {
		t.Fatal("SessionPersistence: got nil, want non-nil")
	}
	if p.SessionPersistence.Type != "HTTP_COOKIE" {
		t.Errorf("SessionPersistence.Type: got %q, want %q", p.SessionPersistence.Type, "HTTP_COOKIE")
	}
	if p.SessionPersistence.CookieName != "SESSION_ID" {
		t.Errorf("SessionPersistence.CookieName: got %q, want %q", p.SessionPersistence.CookieName, "SESSION_ID")
	}
}

func TestListHealthMonitorsOutput_Parse(t *testing.T) {
	raw := `{
		"healthmonitors": [
			{
				"id": "hm-001",
				"name": "my-hm",
				"tenant_id": "tenant-001",
				"pool_id": "pool-001",
				"type": "HTTP",
				"delay": 5,
				"timeout": 3,
				"max_retries": 3,
				"max_retries_down": 3,
				"http_method": "GET",
				"url_path": "/health",
				"expected_codes": "200",
				"admin_state_up": true,
				"provisioning_status": "ACTIVE",
				"operating_status": "ONLINE"
			}
		]
	}`
	var out loadbalancer.ListHealthMonitorsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if len(out.HealthMonitors) != 1 {
		t.Fatalf("expected 1 health monitor, got %d", len(out.HealthMonitors))
	}
	hm := out.HealthMonitors[0]
	if hm.ID != "hm-001" {
		t.Errorf("ID: got %q, want %q", hm.ID, "hm-001")
	}
	if hm.Type != "HTTP" {
		t.Errorf("Type: got %q, want %q", hm.Type, "HTTP")
	}
	if hm.Delay != 5 {
		t.Errorf("Delay: got %d, want 5", hm.Delay)
	}
	if hm.URLPath != "/health" {
		t.Errorf("URLPath: got %q, want %q", hm.URLPath, "/health")
	}
	if hm.ExpectedCodes != "200" {
		t.Errorf("ExpectedCodes: got %q, want %q", hm.ExpectedCodes, "200")
	}
}

// ---------------------------------------------------------------------------
// Request Build Tests
// ---------------------------------------------------------------------------

func TestCreateLoadBalancerRequest_Build(t *testing.T) {
	req := loadbalancer.CreateLoadBalancerRequest{
		LoadBalancer: loadbalancer.CreateLoadBalancerInput{
			Name:        "my-lb",
			VIPSubnetID: "subnet-001",
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
	lbMap, ok := m["loadbalancer"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'loadbalancer' key in output")
	}
	if lbMap["name"] != "my-lb" {
		t.Errorf("name: got %v, want %q", lbMap["name"], "my-lb")
	}
	if lbMap["vip_subnet_id"] != "subnet-001" {
		t.Errorf("vip_subnet_id: got %v, want %q", lbMap["vip_subnet_id"], "subnet-001")
	}
}

func TestCreateListenerRequest_Build(t *testing.T) {
	req := loadbalancer.CreateListenerRequest{
		Listener: loadbalancer.CreateListenerInput{
			Name:           "http-listener",
			LoadBalancerID: "lb-001",
			Protocol:       "HTTP",
			ProtocolPort:   80,
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
	listenerMap, ok := m["listener"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'listener' key in output")
	}
	if listenerMap["protocol"] != "HTTP" {
		t.Errorf("protocol: got %v, want %q", listenerMap["protocol"], "HTTP")
	}
	if listenerMap["protocol_port"] != float64(80) {
		t.Errorf("protocol_port: got %v, want 80", listenerMap["protocol_port"])
	}
	if listenerMap["loadbalancer_id"] != "lb-001" {
		t.Errorf("loadbalancer_id: got %v, want %q", listenerMap["loadbalancer_id"], "lb-001")
	}
}

func TestCreateHealthMonitorRequest_Build(t *testing.T) {
	req := loadbalancer.CreateHealthMonitorRequest{
		HealthMonitor: loadbalancer.CreateHealthMonitorInput{
			PoolID:     "pool-001",
			Type:       "HTTP",
			Delay:      5,
			Timeout:    3,
			MaxRetries: 3,
			HTTPMethod: "GET",
			URLPath:    "/health",
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
	hmMap, ok := m["healthmonitor"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'healthmonitor' key in output")
	}
	if hmMap["pool_id"] != "pool-001" {
		t.Errorf("pool_id: got %v, want %q", hmMap["pool_id"], "pool-001")
	}
	if hmMap["type"] != "HTTP" {
		t.Errorf("type: got %v, want %q", hmMap["type"], "HTTP")
	}
	if hmMap["delay"] != float64(5) {
		t.Errorf("delay: got %v, want 5", hmMap["delay"])
	}
	if hmMap["url_path"] != "/health" {
		t.Errorf("url_path: got %v, want %q", hmMap["url_path"], "/health")
	}
}

func TestUpdateLoadBalancerRequest_Build(t *testing.T) {
	req := loadbalancer.UpdateLoadBalancerRequest{
		LoadBalancer: loadbalancer.UpdateLoadBalancerInput{
			Name:        "updated-lb",
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
	lbMap, ok := m["loadbalancer"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'loadbalancer' key in output")
	}
	if lbMap["name"] != "updated-lb" {
		t.Errorf("name: got %v, want %q", lbMap["name"], "updated-lb")
	}
}

func TestCreateL7PolicyRequest_Build(t *testing.T) {
	req := loadbalancer.CreateL7PolicyRequest{
		L7Policy: loadbalancer.CreateL7PolicyInput{
			ListenerID:     "listener-001",
			Action:         "REDIRECT_TO_POOL",
			RedirectPoolID: "pool-002",
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
	policyMap, ok := m["l7policy"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'l7policy' key in output")
	}
	if policyMap["listener_id"] != "listener-001" {
		t.Errorf("listener_id: got %v, want %q", policyMap["listener_id"], "listener-001")
	}
	if policyMap["action"] != "REDIRECT_TO_POOL" {
		t.Errorf("action: got %v, want %q", policyMap["action"], "REDIRECT_TO_POOL")
	}
}
