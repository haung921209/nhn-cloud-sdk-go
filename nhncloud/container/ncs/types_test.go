package ncs

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
	// NCS common response format: all HTTP 200, check resultCode for actual status
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
}

// ---------------------------------------------------------------------------
// Container
// ---------------------------------------------------------------------------

func TestContainer_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Container{}), []string{
		"Name",
		"Image",
	})
}

func TestContainer_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Container{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Image", "image")
	testutil.AssertStructHasJSONTag(t, typ, "Command", "command")
	testutil.AssertStructHasJSONTag(t, typ, "Args", "args")
	testutil.AssertStructHasJSONTag(t, typ, "Ports", "ports")
	testutil.AssertStructHasJSONTag(t, typ, "Env", "env")
	testutil.AssertStructHasJSONTag(t, typ, "Resources", "resources")
	testutil.AssertStructHasJSONTag(t, typ, "VolumeMounts", "volumeMounts")
}

func TestContainer_Parse(t *testing.T) {
	raw := `{
		"name": "web",
		"image": "nginx:latest",
		"command": ["/bin/sh"],
		"args": ["-c","echo hi"],
		"ports": [{"containerPort": 80, "protocol": "TCP"}],
		"env": [{"name":"ENV_VAR","value":"hello"}],
		"resources": {
			"limits": {"cpu":"500m","memory":"256Mi"},
			"requests": {"cpu":"100m","memory":"128Mi"}
		},
		"volumeMounts": [{"name":"data","mountPath":"/data"}]
	}`
	var c Container
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Name != "web" {
		t.Errorf("Name: got %q, want web", c.Name)
	}
	if c.Image != "nginx:latest" {
		t.Errorf("Image: got %q, want nginx:latest", c.Image)
	}
	if len(c.Ports) != 1 || c.Ports[0].ContainerPort != 80 {
		t.Errorf("Ports: got %v", c.Ports)
	}
	if len(c.Env) != 1 || c.Env[0].Name != "ENV_VAR" {
		t.Errorf("Env: got %v", c.Env)
	}
	if c.Resources == nil {
		t.Fatal("Resources should not be nil")
	}
	if c.Resources.Limits.CPU != "500m" {
		t.Errorf("Resources.Limits.CPU: got %q, want 500m", c.Resources.Limits.CPU)
	}
}

// ---------------------------------------------------------------------------
// ContainerPort
// ---------------------------------------------------------------------------

func TestContainerPort_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ContainerPort{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "ContainerPort", "containerPort")
	testutil.AssertStructHasJSONTag(t, typ, "Protocol", "protocol")
}

// ---------------------------------------------------------------------------
// EnvVar
// ---------------------------------------------------------------------------

func TestEnvVar_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(EnvVar{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Value", "value")
}

// ---------------------------------------------------------------------------
// ResourceRequirements / ResourceList
// ---------------------------------------------------------------------------

func TestResourceList_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ResourceList{})
	testutil.AssertStructHasJSONTag(t, typ, "CPU", "cpu")
	testutil.AssertStructHasJSONTag(t, typ, "Memory", "memory")
	testutil.AssertStructHasJSONTag(t, typ, "GPU", "nvidia.com/gpu")
}

func TestResourceList_GPUTag(t *testing.T) {
	// The GPU field must use the exact Kubernetes extended-resource key
	r := ResourceList{GPU: "1"}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["nvidia.com/gpu"] != "1" {
		t.Errorf("nvidia.com/gpu: got %v", m["nvidia.com/gpu"])
	}
}

// ---------------------------------------------------------------------------
// VolumeMount
// ---------------------------------------------------------------------------

func TestVolumeMount_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(VolumeMount{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "MountPath", "mountPath")
	testutil.AssertStructHasJSONTag(t, typ, "ReadOnly", "readOnly")
}

// ---------------------------------------------------------------------------
// Event
// ---------------------------------------------------------------------------

func TestEvent_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Event{})
	testutil.AssertStructHasJSONTag(t, typ, "EventID", "eventId")
	testutil.AssertStructHasJSONTag(t, typ, "EventType", "eventType")
	testutil.AssertStructHasJSONTag(t, typ, "Reason", "reason")
	testutil.AssertStructHasJSONTag(t, typ, "Message", "message")
	testutil.AssertStructHasJSONTag(t, typ, "Count", "count")
	testutil.AssertStructHasJSONTag(t, typ, "FirstTime", "firstTimestamp")
	testutil.AssertStructHasJSONTag(t, typ, "LastTime", "lastTimestamp")
}

// ---------------------------------------------------------------------------
// NetworkPolicy / NetworkPolicyRule
// ---------------------------------------------------------------------------

func TestNetworkPolicyRule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(NetworkPolicyRule{})
	testutil.AssertStructHasJSONTag(t, typ, "Direction", "direction")
	testutil.AssertStructHasJSONTag(t, typ, "Protocol", "protocol")
	testutil.AssertStructHasJSONTag(t, typ, "Port", "port")
	testutil.AssertStructHasJSONTag(t, typ, "FromCIDR", "fromCidr")
	testutil.AssertStructHasJSONTag(t, typ, "ToCIDR", "toCidr")
	testutil.AssertStructHasJSONTag(t, typ, "Action", "action")
}

func TestNetworkPolicy_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(NetworkPolicy{})
	testutil.AssertStructHasJSONTag(t, typ, "PolicyID", "policyId")
	testutil.AssertStructHasJSONTag(t, typ, "PolicyName", "policyName")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "Rules", "rules")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestCreateNetworkPolicyInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateNetworkPolicyInput{})
	testutil.AssertStructHasJSONTag(t, typ, "PolicyName", "policyName")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "Rules", "rules")
}

func TestNetworkPolicy_Parse(t *testing.T) {
	raw := `{
		"policyId": "np-001",
		"policyName": "allow-web",
		"description": "Allow web traffic",
		"rules": [
			{"direction":"Ingress","protocol":"TCP","port":80,"fromCidr":["0.0.0.0/0"],"action":"Allow"}
		],
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-02T00:00:00Z"
	}`
	var np NetworkPolicy
	if err := json.Unmarshal([]byte(raw), &np); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if np.PolicyID != "np-001" {
		t.Errorf("PolicyID: got %q, want np-001", np.PolicyID)
	}
	if len(np.Rules) != 1 {
		t.Fatalf("Rules: got %d, want 1", len(np.Rules))
	}
	r := np.Rules[0]
	if r.Direction != "Ingress" {
		t.Errorf("Direction: got %q, want Ingress", r.Direction)
	}
	if r.Action != "Allow" {
		t.Errorf("Action: got %q, want Allow", r.Action)
	}
	if len(r.FromCIDR) != 1 || r.FromCIDR[0] != "0.0.0.0/0" {
		t.Errorf("FromCIDR: got %v", r.FromCIDR)
	}
}
