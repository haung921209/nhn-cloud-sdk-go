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
// Workload
// ---------------------------------------------------------------------------

func TestWorkload_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Workload{}), []string{
		"ID",
		"Name",
		"Namespace",
		"Type",
		"Replicas",
		"AvailableReplicas",
		"Status",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestWorkload_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Workload{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Namespace", "namespace")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "Replicas", "replicas")
	testutil.AssertStructHasJSONTag(t, typ, "AvailableReplicas", "availableReplicas")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "Containers", "containers")
	testutil.AssertStructHasJSONTag(t, typ, "Labels", "labels")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestWorkload_Parse(t *testing.T) {
	raw := `{
		"id": "wl-001",
		"name": "my-workload",
		"namespace": "default",
		"type": "Deployment",
		"replicas": 3,
		"availableReplicas": 2,
		"status": "Running",
		"containers": [{"name":"app","image":"nginx:latest"}],
		"labels": {"app":"web"},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-06-01T00:00:00Z"
	}`
	var wl Workload
	if err := json.Unmarshal([]byte(raw), &wl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if wl.ID != "wl-001" {
		t.Errorf("ID: got %q, want wl-001", wl.ID)
	}
	if wl.Replicas != 3 {
		t.Errorf("Replicas: got %d, want 3", wl.Replicas)
	}
	if wl.AvailableReplicas != 2 {
		t.Errorf("AvailableReplicas: got %d, want 2", wl.AvailableReplicas)
	}
	if wl.Labels["app"] != "web" {
		t.Errorf("Labels[app]: got %q, want web", wl.Labels["app"])
	}
}

// ---------------------------------------------------------------------------
// CreateWorkloadInput
// ---------------------------------------------------------------------------

func TestCreateWorkloadInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(CreateWorkloadInput{}), []string{
		"Name",
		"Type",
		"Containers",
	})
}

func TestCreateWorkloadInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateWorkloadInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Namespace", "namespace")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "Replicas", "replicas")
	testutil.AssertStructHasJSONTag(t, typ, "Containers", "containers")
	testutil.AssertStructHasJSONTag(t, typ, "Labels", "labels")
	testutil.AssertStructHasJSONTag(t, typ, "Volumes", "volumes")
}

func TestCreateWorkloadInput_Marshal(t *testing.T) {
	input := CreateWorkloadInput{
		Name: "my-wl",
		Type: "Deployment",
		Containers: []Container{
			{Name: "app", Image: "nginx:latest"},
		},
		Replicas: 2,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["name"] != "my-wl" {
		t.Errorf("name: got %v", m["name"])
	}
	if m["type"] != "Deployment" {
		t.Errorf("type: got %v", m["type"])
	}
	containers, ok := m["containers"].([]interface{})
	if !ok || len(containers) != 1 {
		t.Errorf("containers: got %v", m["containers"])
	}
}

// ---------------------------------------------------------------------------
// Volume types
// ---------------------------------------------------------------------------

func TestVolume_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Volume{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "EmptyDir", "emptyDir")
	testutil.AssertStructHasJSONTag(t, typ, "ConfigMap", "configMap")
	testutil.AssertStructHasJSONTag(t, typ, "Secret", "secret")
	testutil.AssertStructHasJSONTag(t, typ, "PersistentVolumeClaim", "persistentVolumeClaim")
}

func TestSecretVolumeSource_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(SecretVolumeSource{})
	testutil.AssertStructHasJSONTag(t, typ, "SecretName", "secretName")
}

func TestPVCVolumeSource_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(PVCVolumeSource{})
	testutil.AssertStructHasJSONTag(t, typ, "ClaimName", "claimName")
	testutil.AssertStructHasJSONTag(t, typ, "ReadOnly", "readOnly")
}

// ---------------------------------------------------------------------------
// Template
// ---------------------------------------------------------------------------

func TestTemplate_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Template{}), []string{
		"ID",
		"Name",
		"Version",
		"Type",
		"IsPublic",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestTemplate_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Template{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Description", "description")
	testutil.AssertStructHasJSONTag(t, typ, "Version", "version")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "IsPublic", "isPublic")
	testutil.AssertStructHasJSONTag(t, typ, "Containers", "containers")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestTemplate_Parse(t *testing.T) {
	raw := `{
		"id": "tmpl-001",
		"name": "nginx-template",
		"description": "A basic nginx template",
		"version": "v1",
		"type": "Deployment",
		"isPublic": false,
		"containers": [{"name":"web","image":"nginx:1.25"}],
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-06-01T00:00:00Z"
	}`
	var tmpl Template
	if err := json.Unmarshal([]byte(raw), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tmpl.ID != "tmpl-001" {
		t.Errorf("ID: got %q, want tmpl-001", tmpl.ID)
	}
	if tmpl.Version != "v1" {
		t.Errorf("Version: got %q, want v1", tmpl.Version)
	}
	if len(tmpl.Containers) != 1 {
		t.Fatalf("Containers: got %d, want 1", len(tmpl.Containers))
	}
}

// ---------------------------------------------------------------------------
// Service / ServicePort
// ---------------------------------------------------------------------------

func TestService_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Service{}), []string{
		"ID",
		"Name",
		"Namespace",
		"Type",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestService_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Service{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "id")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Namespace", "namespace")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "ClusterIP", "clusterIP")
	testutil.AssertStructHasJSONTag(t, typ, "ExternalIP", "externalIP")
	testutil.AssertStructHasJSONTag(t, typ, "Ports", "ports")
	testutil.AssertStructHasJSONTag(t, typ, "Selector", "selector")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updatedAt")
}

func TestServicePort_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ServicePort{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Port", "port")
	testutil.AssertStructHasJSONTag(t, typ, "TargetPort", "targetPort")
	testutil.AssertStructHasJSONTag(t, typ, "NodePort", "nodePort")
	testutil.AssertStructHasJSONTag(t, typ, "Protocol", "protocol")
}

func TestCreateServiceInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateServiceInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Namespace", "namespace")
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "Ports", "ports")
	testutil.AssertStructHasJSONTag(t, typ, "Selector", "selector")
}

func TestCreateServiceInput_Marshal(t *testing.T) {
	input := CreateServiceInput{
		Name: "my-svc",
		Type: "ClusterIP",
		Ports: []ServicePort{
			{Port: 80, TargetPort: 8080, Protocol: "TCP"},
		},
		Selector: map[string]string{"app": "web"},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["name"] != "my-svc" {
		t.Errorf("name: got %v", m["name"])
	}
	if m["type"] != "ClusterIP" {
		t.Errorf("type: got %v", m["type"])
	}
	selector, ok := m["selector"].(map[string]interface{})
	if !ok || selector["app"] != "web" {
		t.Errorf("selector: got %v", m["selector"])
	}
}

// ---------------------------------------------------------------------------
// LogEntry
// ---------------------------------------------------------------------------

func TestLogEntry_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(LogEntry{})
	testutil.AssertStructHasJSONTag(t, typ, "Timestamp", "timestamp")
	testutil.AssertStructHasJSONTag(t, typ, "Message", "message")
	testutil.AssertStructHasJSONTag(t, typ, "Stream", "stream")
}

func TestGetWorkloadLogsOutput_Parse(t *testing.T) {
	raw := `{
		"header": {"isSuccessful": true, "resultCode": 200, "resultMessage": "SUCCESS"},
		"logs": [
			{"timestamp":"2024-01-01T00:00:01Z","message":"starting server","stream":"stdout"},
			{"timestamp":"2024-01-01T00:00:02Z","message":"ready","stream":"stdout"}
		]
	}`
	var out GetWorkloadLogsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Logs) != 2 {
		t.Fatalf("Logs: got %d, want 2", len(out.Logs))
	}
	if out.Logs[0].Stream != "stdout" {
		t.Errorf("Logs[0].Stream: got %q, want stdout", out.Logs[0].Stream)
	}
}

// ---------------------------------------------------------------------------
// HealthProbe / HealthCheckConfig
// ---------------------------------------------------------------------------

func TestHealthProbe_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(HealthProbe{})
	testutil.AssertStructHasJSONTag(t, typ, "Type", "type")
	testutil.AssertStructHasJSONTag(t, typ, "HTTPGet", "httpGet")
	testutil.AssertStructHasJSONTag(t, typ, "TCPSocket", "tcpSocket")
	testutil.AssertStructHasJSONTag(t, typ, "Exec", "exec")
	testutil.AssertStructHasJSONTag(t, typ, "InitialDelaySeconds", "initialDelaySeconds")
	testutil.AssertStructHasJSONTag(t, typ, "PeriodSeconds", "periodSeconds")
	testutil.AssertStructHasJSONTag(t, typ, "TimeoutSeconds", "timeoutSeconds")
	testutil.AssertStructHasJSONTag(t, typ, "SuccessThreshold", "successThreshold")
	testutil.AssertStructHasJSONTag(t, typ, "FailureThreshold", "failureThreshold")
}

func TestHTTPGetAction_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(HTTPGetAction{})
	testutil.AssertStructHasJSONTag(t, typ, "Path", "path")
	testutil.AssertStructHasJSONTag(t, typ, "Port", "port")
	testutil.AssertStructHasJSONTag(t, typ, "Scheme", "scheme")
}

func TestHealthCheckConfig_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(HealthCheckConfig{})
	testutil.AssertStructHasJSONTag(t, typ, "LivenessProbe", "livenessProbe")
	testutil.AssertStructHasJSONTag(t, typ, "ReadinessProbe", "readinessProbe")
}

func TestHealthCheckConfig_Marshal(t *testing.T) {
	cfg := HealthCheckConfig{
		LivenessProbe: &HealthProbe{
			Type:                "httpGet",
			InitialDelaySeconds: 10,
			PeriodSeconds:       5,
			HTTPGet:             &HTTPGetAction{Path: "/healthz", Port: 8080},
		},
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	lp, ok := m["livenessProbe"].(map[string]interface{})
	if !ok {
		t.Fatalf("livenessProbe not a map: %v", m["livenessProbe"])
	}
	if lp["type"] != "httpGet" {
		t.Errorf("livenessProbe.type: got %v", lp["type"])
	}
	if _, ok := m["readinessProbe"]; ok {
		t.Error("readinessProbe should be omitted when nil")
	}
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
// PersistentVolume / VolumeAttachInput
// ---------------------------------------------------------------------------

func TestPersistentVolume_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(PersistentVolume{})
	testutil.AssertStructHasJSONTag(t, typ, "VolumeID", "volumeId")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "Size", "size")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "VolumeType", "volumeType")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "createdAt")
}

func TestVolumeAttachInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(VolumeAttachInput{})
	testutil.AssertStructHasJSONTag(t, typ, "VolumeID", "volumeId")
	testutil.AssertStructHasJSONTag(t, typ, "MountPath", "mountPath")
	testutil.AssertStructHasJSONTag(t, typ, "ReadOnly", "readOnly")
}

// ---------------------------------------------------------------------------
// ExecInput / ExecOutput
// ---------------------------------------------------------------------------

func TestExecInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ExecInput{})
	testutil.AssertStructHasJSONTag(t, typ, "ContainerName", "containerName")
	testutil.AssertStructHasJSONTag(t, typ, "Command", "command")
	testutil.AssertStructHasJSONTag(t, typ, "Stdin", "stdin")
	testutil.AssertStructHasJSONTag(t, typ, "Stdout", "stdout")
	testutil.AssertStructHasJSONTag(t, typ, "Stderr", "stderr")
	testutil.AssertStructHasJSONTag(t, typ, "TTY", "tty")
}

func TestExecOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ExecOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "ExecID", "execId")
	testutil.AssertStructHasJSONTag(t, typ, "Output", "output")
	testutil.AssertStructHasJSONTag(t, typ, "ExitCode", "exitCode")
	testutil.AssertStructHasJSONTag(t, typ, "Error", "error")
	testutil.AssertStructHasJSONTag(t, typ, "StartedAt", "startedAt")
	testutil.AssertStructHasJSONTag(t, typ, "FinishedAt", "finishedAt")
}

// ---------------------------------------------------------------------------
// ContainerStatus
// ---------------------------------------------------------------------------

func TestContainerStatus_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ContainerStatus{})
	testutil.AssertStructHasJSONTag(t, typ, "ContainerName", "containerName")
	testutil.AssertStructHasJSONTag(t, typ, "State", "state")
	testutil.AssertStructHasJSONTag(t, typ, "Ready", "ready")
	testutil.AssertStructHasJSONTag(t, typ, "RestartCount", "restartCount")
	testutil.AssertStructHasJSONTag(t, typ, "Image", "image")
	testutil.AssertStructHasJSONTag(t, typ, "ContainerID", "containerId")
}

// ---------------------------------------------------------------------------
// AutoScalingPolicy / ConfigureAutoScalingInput
// ---------------------------------------------------------------------------

func TestAutoScalingPolicy_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(AutoScalingPolicy{}), []string{
		"MinReplicas",
		"MaxReplicas",
	})
}

func TestAutoScalingPolicy_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(AutoScalingPolicy{})
	testutil.AssertStructHasJSONTag(t, typ, "MinReplicas", "minReplicas")
	testutil.AssertStructHasJSONTag(t, typ, "MaxReplicas", "maxReplicas")
	testutil.AssertStructHasJSONTag(t, typ, "TargetCPUUtilizationPercentage", "targetCpuUtilizationPercentage")
	testutil.AssertStructHasJSONTag(t, typ, "TargetMemoryUtilizationPercentage", "targetMemoryUtilizationPercentage")
	testutil.AssertStructHasJSONTag(t, typ, "ScaleUpStabilizationSeconds", "scaleUpStabilizationSeconds")
	testutil.AssertStructHasJSONTag(t, typ, "ScaleDownStabilizationSeconds", "scaleDownStabilizationSeconds")
}

func TestConfigureAutoScalingInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ConfigureAutoScalingInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Enabled", "enabled")
	testutil.AssertStructHasJSONTag(t, typ, "Policy", "policy")
}

func TestConfigureAutoScalingInput_Marshal(t *testing.T) {
	tr := true
	_ = tr // suppress unused warning
	input := ConfigureAutoScalingInput{
		Enabled: true,
		Policy: &AutoScalingPolicy{
			MinReplicas: 2,
			MaxReplicas: 10,
			TargetCPUUtilizationPercentage: 70,
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["enabled"] != true {
		t.Errorf("enabled: got %v", m["enabled"])
	}
	policy, ok := m["policy"].(map[string]interface{})
	if !ok {
		t.Fatalf("policy not a map: %v", m["policy"])
	}
	if policy["minReplicas"].(float64) != 2 {
		t.Errorf("minReplicas: got %v", policy["minReplicas"])
	}
	if policy["maxReplicas"].(float64) != 10 {
		t.Errorf("maxReplicas: got %v", policy["maxReplicas"])
	}
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
