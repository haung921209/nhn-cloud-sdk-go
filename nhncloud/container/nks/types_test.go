package nks

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// Cluster
// ---------------------------------------------------------------------------

func TestCluster_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(Cluster{}), []string{
		"ID",
		"Name",
		"TenantID",
		"Status",
		"APIAddress",
		"MasterAddresses",
		"K8sVersion",
		"NodeCount",
		"MasterCount",
		"FlavorID",
		"MasterFlavorID",
		"KeyPair",
		"ClusterTemplateID",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestCluster_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(Cluster{})
	// API spec fields from NKS docs
	testutil.AssertStructHasJSONTag(t, typ, "ID", "uuid")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "TenantID", "tenant_id")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "StatusReason", "status_reason")
	testutil.AssertStructHasJSONTag(t, typ, "HealthStatus", "health_status")
	testutil.AssertStructHasJSONTag(t, typ, "APIAddress", "api_address")
	testutil.AssertStructHasJSONTag(t, typ, "MasterAddresses", "master_addresses")
	testutil.AssertStructHasJSONTag(t, typ, "NodeAddresses", "node_addresses")
	testutil.AssertStructHasJSONTag(t, typ, "K8sVersion", "coe_version")
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MasterCount", "master_count")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorID", "flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "MasterFlavorID", "master_flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "KeyPair", "keypair")
	testutil.AssertStructHasJSONTag(t, typ, "ClusterTemplateID", "cluster_template_id")
	testutil.AssertStructHasJSONTag(t, typ, "NetworkID", "fixed_network")
	testutil.AssertStructHasJSONTag(t, typ, "SubnetID", "fixed_subnet")
	testutil.AssertStructHasJSONTag(t, typ, "Labels", "labels")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "created_at")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updated_at")
}

func TestCluster_Parse(t *testing.T) {
	raw := `{
		"uuid": "cluster-uuid-001",
		"name": "my-cluster",
		"tenant_id": "tenant-001",
		"status": "CREATE_COMPLETE",
		"health_status": "FRESH",
		"api_address": "https://1.2.3.4:6443",
		"master_addresses": ["10.0.0.1"],
		"node_addresses": ["10.0.0.2","10.0.0.3"],
		"coe_version": "v1.29.3",
		"node_count": 2,
		"master_count": 3,
		"flavor_id": "flavor-worker-uuid",
		"master_flavor_id": "flavor-master-uuid",
		"keypair": "my-keypair",
		"cluster_template_id": "iaas_console",
		"fixed_network": "vpc-uuid",
		"fixed_subnet": "subnet-uuid",
		"labels": {"kube_tag":"v1.29.3"},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-06-01T00:00:00Z"
	}`
	var c Cluster
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.ID != "cluster-uuid-001" {
		t.Errorf("ID: got %q, want cluster-uuid-001", c.ID)
	}
	if c.Name != "my-cluster" {
		t.Errorf("Name: got %q, want my-cluster", c.Name)
	}
	if c.TenantID != "tenant-001" {
		t.Errorf("TenantID: got %q, want tenant-001", c.TenantID)
	}
	if c.Status != "CREATE_COMPLETE" {
		t.Errorf("Status: got %q, want CREATE_COMPLETE", c.Status)
	}
	if c.HealthStatus != "FRESH" {
		t.Errorf("HealthStatus: got %q, want FRESH", c.HealthStatus)
	}
	if c.K8sVersion != "v1.29.3" {
		t.Errorf("K8sVersion: got %q, want v1.29.3", c.K8sVersion)
	}
	if c.NodeCount != 2 {
		t.Errorf("NodeCount: got %d, want 2", c.NodeCount)
	}
	if c.MasterCount != 3 {
		t.Errorf("MasterCount: got %d, want 3", c.MasterCount)
	}
	if len(c.MasterAddresses) != 1 || c.MasterAddresses[0] != "10.0.0.1" {
		t.Errorf("MasterAddresses: got %v", c.MasterAddresses)
	}
	if len(c.NodeAddresses) != 2 {
		t.Errorf("NodeAddresses: got %d elements, want 2", len(c.NodeAddresses))
	}
	if c.Labels["kube_tag"] != "v1.29.3" {
		t.Errorf("Labels[kube_tag]: got %q, want v1.29.3", c.Labels["kube_tag"])
	}
	if c.NetworkID != "vpc-uuid" {
		t.Errorf("NetworkID: got %q, want vpc-uuid", c.NetworkID)
	}
	if c.SubnetID != "subnet-uuid" {
		t.Errorf("SubnetID: got %q, want subnet-uuid", c.SubnetID)
	}
}

// ---------------------------------------------------------------------------
// ListClustersOutput
// ---------------------------------------------------------------------------

func TestListClustersOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(ListClustersOutput{}), []string{
		"Clusters",
	})
}

func TestListClustersOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ListClustersOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "Clusters", "clusters")
}

func TestListClustersOutput_Parse(t *testing.T) {
	raw := `{
		"clusters": [
			{
				"uuid": "c1",
				"name": "cluster-1",
				"tenant_id": "t1",
				"status": "CREATE_COMPLETE",
				"api_address": "https://1.1.1.1:6443",
				"master_addresses": ["10.0.0.1"],
				"coe_version": "v1.28.0",
				"node_count": 1,
				"master_count": 1,
				"flavor_id": "f1",
				"master_flavor_id": "f2",
				"keypair": "kp1",
				"cluster_template_id": "iaas_console",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}
		]
	}`
	var out ListClustersOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Clusters) != 1 {
		t.Fatalf("Clusters: got %d, want 1", len(out.Clusters))
	}
	if out.Clusters[0].ID != "c1" {
		t.Errorf("Clusters[0].ID: got %q, want c1", out.Clusters[0].ID)
	}
}

// ---------------------------------------------------------------------------
// CreateClusterInput
// ---------------------------------------------------------------------------

func TestCreateClusterInput_RequiredFields(t *testing.T) {
	// API spec requires: name, cluster_template_id, keypair, node_count (String per spec),
	// flavor_id, fixed_network, fixed_subnet, labels
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(CreateClusterInput{}), []string{
		"Name",
		"ClusterTemplateID",
		"KeyPair",
		"NodeCount",
		"FlavorID",
		"NetworkID",
		"SubnetID",
		"Labels",
	})
}

func TestCreateClusterInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateClusterInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "ClusterTemplateID", "cluster_template_id")
	testutil.AssertStructHasJSONTag(t, typ, "MasterCount", "master_count")
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MasterFlavorID", "master_flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorID", "flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "KeyPair", "keypair")
	testutil.AssertStructHasJSONTag(t, typ, "NetworkID", "fixed_network")
	testutil.AssertStructHasJSONTag(t, typ, "SubnetID", "fixed_subnet")
	testutil.AssertStructHasJSONTag(t, typ, "Labels", "labels")
}

func TestCreateClusterInput_Marshal(t *testing.T) {
	input := CreateClusterInput{
		Name:              "test-cluster",
		ClusterTemplateID: "iaas_console",
		KeyPair:           "my-keypair",
		NodeCount:         "3",
		FlavorID:          "flavor-uuid",
		NetworkID:         "vpc-uuid",
		SubnetID:          "subnet-uuid",
		Labels: map[string]string{
			"kube_tag":          "v1.29.3",
			"availability_zone": "kr1-pub-a",
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["name"] != "test-cluster" {
		t.Errorf("name: got %v", m["name"])
	}
	if m["cluster_template_id"] != "iaas_console" {
		t.Errorf("cluster_template_id: got %v", m["cluster_template_id"])
	}
	if m["keypair"] != "my-keypair" {
		t.Errorf("keypair: got %v", m["keypair"])
	}
	if m["node_count"] != "3" {
		t.Errorf("node_count: got %v, want \"3\"", m["node_count"])
	}
	if m["fixed_network"] != "vpc-uuid" {
		t.Errorf("fixed_network: got %v", m["fixed_network"])
	}
	if m["fixed_subnet"] != "subnet-uuid" {
		t.Errorf("fixed_subnet: got %v", m["fixed_subnet"])
	}
	labels, ok := m["labels"].(map[string]interface{})
	if !ok {
		t.Fatalf("labels not a map: %v", m["labels"])
	}
	if labels["kube_tag"] != "v1.29.3" {
		t.Errorf("labels.kube_tag: got %v", labels["kube_tag"])
	}
}

// ---------------------------------------------------------------------------
// UpdateClusterInput
// ---------------------------------------------------------------------------

func TestUpdateClusterInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(UpdateClusterInput{})
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
}

func TestUpdateClusterInput_Marshal(t *testing.T) {
	input := UpdateClusterInput{NodeCount: 5}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["node_count"].(float64) != 5 {
		t.Errorf("node_count: got %v", m["node_count"])
	}
}

// ---------------------------------------------------------------------------
// NodeGroup
// ---------------------------------------------------------------------------

func TestNodeGroup_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(NodeGroup{}), []string{
		"ID",
		"Name",
		"ClusterID",
		"NodeCount",
		"MinNodeCount",
		"MaxNodeCount",
		"FlavorID",
		"ImageID",
		"Role",
		"Status",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestNodeGroup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(NodeGroup{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "uuid")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "ClusterID", "cluster_id")
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MinNodeCount", "min_node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MaxNodeCount", "max_node_count")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorID", "flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "ImageID", "image_id")
	testutil.AssertStructHasJSONTag(t, typ, "Role", "role")
	testutil.AssertStructHasJSONTag(t, typ, "Status", "status")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "created_at")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updated_at")
}

func TestNodeGroup_Parse(t *testing.T) {
	raw := `{
		"uuid": "ng-uuid-001",
		"name": "default-worker",
		"cluster_id": "cluster-uuid-001",
		"node_count": 3,
		"min_node_count": 1,
		"max_node_count": 10,
		"flavor_id": "flavor-uuid",
		"image_id": "image-uuid",
		"role": "worker",
		"status": "UPDATE_COMPLETE",
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-06-01T00:00:00Z"
	}`
	var ng NodeGroup
	if err := json.Unmarshal([]byte(raw), &ng); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ng.ID != "ng-uuid-001" {
		t.Errorf("ID: got %q, want ng-uuid-001", ng.ID)
	}
	if ng.ClusterID != "cluster-uuid-001" {
		t.Errorf("ClusterID: got %q", ng.ClusterID)
	}
	if ng.NodeCount != 3 {
		t.Errorf("NodeCount: got %d, want 3", ng.NodeCount)
	}
	if ng.MinNodeCount != 1 {
		t.Errorf("MinNodeCount: got %d, want 1", ng.MinNodeCount)
	}
	if ng.MaxNodeCount != 10 {
		t.Errorf("MaxNodeCount: got %d, want 10", ng.MaxNodeCount)
	}
	if ng.Role != "worker" {
		t.Errorf("Role: got %q, want worker", ng.Role)
	}
}

// ---------------------------------------------------------------------------
// ListNodeGroupsOutput
// ---------------------------------------------------------------------------

func TestListNodeGroupsOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ListNodeGroupsOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "NodeGroups", "nodegroups")
}

func TestListNodeGroupsOutput_Parse(t *testing.T) {
	raw := `{
		"nodegroups": [
			{
				"uuid": "ng-1",
				"name": "ng-default",
				"cluster_id": "c1",
				"node_count": 2,
				"min_node_count": 1,
				"max_node_count": 5,
				"flavor_id": "f1",
				"image_id": "img1",
				"role": "worker",
				"status": "CREATE_COMPLETE",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}
		]
	}`
	var out ListNodeGroupsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.NodeGroups) != 1 {
		t.Fatalf("NodeGroups: got %d, want 1", len(out.NodeGroups))
	}
	if out.NodeGroups[0].Name != "ng-default" {
		t.Errorf("NodeGroups[0].Name: got %q, want ng-default", out.NodeGroups[0].Name)
	}
}

// ---------------------------------------------------------------------------
// CreateNodeGroupInput
// ---------------------------------------------------------------------------

func TestCreateNodeGroupInput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(CreateNodeGroupInput{}), []string{
		"Name",
		"FlavorID",
		"NodeCount",
	})
}

func TestCreateNodeGroupInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(CreateNodeGroupInput{})
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorID", "flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "ImageID", "image_id")
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MinNodeCount", "min_node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MaxNodeCount", "max_node_count")
}

func TestCreateNodeGroupInput_Marshal(t *testing.T) {
	input := CreateNodeGroupInput{
		Name:         "new-ng",
		FlavorID:     "flavor-uuid",
		ImageID:      "image-uuid",
		NodeCount:    2,
		MinNodeCount: 1,
		MaxNodeCount: 5,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if m["name"] != "new-ng" {
		t.Errorf("name: got %v", m["name"])
	}
	if m["flavor_id"] != "flavor-uuid" {
		t.Errorf("flavor_id: got %v", m["flavor_id"])
	}
	if m["node_count"].(float64) != 2 {
		t.Errorf("node_count: got %v", m["node_count"])
	}
	if m["min_node_count"].(float64) != 1 {
		t.Errorf("min_node_count: got %v", m["min_node_count"])
	}
}

// ---------------------------------------------------------------------------
// UpdateNodeGroupInput
// ---------------------------------------------------------------------------

func TestUpdateNodeGroupInput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(UpdateNodeGroupInput{})
	testutil.AssertStructHasJSONTag(t, typ, "NodeCount", "node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MinNodeCount", "min_node_count")
	testutil.AssertStructHasJSONTag(t, typ, "MaxNodeCount", "max_node_count")
}

// ---------------------------------------------------------------------------
// ClusterTemplate
// ---------------------------------------------------------------------------

func TestClusterTemplate_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(ClusterTemplate{}), []string{
		"ID",
		"Name",
		"TenantID",
		"COE",
		"ServerType",
		"NetworkDriver",
		"VolumeDriver",
		"DockerStorageDriver",
		"ExternalNetworkID",
		"MasterFlavor",
		"Flavor",
		"CreatedAt",
		"UpdatedAt",
	})
}

func TestClusterTemplate_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ClusterTemplate{})
	testutil.AssertStructHasJSONTag(t, typ, "ID", "uuid")
	testutil.AssertStructHasJSONTag(t, typ, "Name", "name")
	testutil.AssertStructHasJSONTag(t, typ, "TenantID", "tenant_id")
	testutil.AssertStructHasJSONTag(t, typ, "Public", "public")
	testutil.AssertStructHasJSONTag(t, typ, "COE", "coe")
	testutil.AssertStructHasJSONTag(t, typ, "ServerType", "server_type")
	testutil.AssertStructHasJSONTag(t, typ, "NetworkDriver", "network_driver")
	testutil.AssertStructHasJSONTag(t, typ, "VolumeDriver", "volume_driver")
	testutil.AssertStructHasJSONTag(t, typ, "DockerStorageDriver", "docker_storage_driver")
	testutil.AssertStructHasJSONTag(t, typ, "DockerVolumeSize", "docker_volume_size")
	testutil.AssertStructHasJSONTag(t, typ, "ExternalNetworkID", "external_network_id")
	testutil.AssertStructHasJSONTag(t, typ, "FixedNetwork", "fixed_network")
	testutil.AssertStructHasJSONTag(t, typ, "FixedSubnet", "fixed_subnet")
	testutil.AssertStructHasJSONTag(t, typ, "DNSNameserver", "dns_nameserver")
	testutil.AssertStructHasJSONTag(t, typ, "MasterFlavor", "master_flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "Flavor", "flavor_id")
	testutil.AssertStructHasJSONTag(t, typ, "Labels", "labels")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedAt", "created_at")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedAt", "updated_at")
}

func TestClusterTemplate_Parse(t *testing.T) {
	raw := `{
		"uuid": "tmpl-uuid-001",
		"name": "iaas_console",
		"tenant_id": "tenant-001",
		"public": true,
		"coe": "kubernetes",
		"server_type": "vm",
		"network_driver": "flannel",
		"volume_driver": "cinder",
		"docker_storage_driver": "overlay2",
		"docker_volume_size": 20,
		"external_network_id": "ext-net-uuid",
		"fixed_network": "vpc-uuid",
		"fixed_subnet": "subnet-uuid",
		"dns_nameserver": "8.8.8.8",
		"master_flavor_id": "master-flavor-uuid",
		"flavor_id": "worker-flavor-uuid",
		"labels": {},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-06-01T00:00:00Z"
	}`
	var tmpl ClusterTemplate
	if err := json.Unmarshal([]byte(raw), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tmpl.ID != "tmpl-uuid-001" {
		t.Errorf("ID: got %q", tmpl.ID)
	}
	if tmpl.COE != "kubernetes" {
		t.Errorf("COE: got %q, want kubernetes", tmpl.COE)
	}
	if !tmpl.Public {
		t.Error("Public should be true")
	}
	if tmpl.DockerVolumeSize != 20 {
		t.Errorf("DockerVolumeSize: got %d, want 20", tmpl.DockerVolumeSize)
	}
	if tmpl.MasterFlavor != "master-flavor-uuid" {
		t.Errorf("MasterFlavor: got %q", tmpl.MasterFlavor)
	}
	if tmpl.Flavor != "worker-flavor-uuid" {
		t.Errorf("Flavor: got %q", tmpl.Flavor)
	}
}

// ---------------------------------------------------------------------------
// ListClusterTemplatesOutput
// ---------------------------------------------------------------------------

func TestListClusterTemplatesOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(ListClusterTemplatesOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "ClusterTemplates", "clustertemplates")
}

func TestListClusterTemplatesOutput_Parse(t *testing.T) {
	raw := `{
		"clustertemplates": [
			{
				"uuid": "tmpl-1",
				"name": "iaas_console",
				"tenant_id": "t1",
				"public": true,
				"coe": "kubernetes",
				"server_type": "vm",
				"network_driver": "flannel",
				"volume_driver": "cinder",
				"docker_storage_driver": "overlay2",
				"docker_volume_size": 50,
				"external_network_id": "ext-net",
				"fixed_network": "",
				"fixed_subnet": "",
				"dns_nameserver": "8.8.8.8",
				"master_flavor_id": "mf1",
				"flavor_id": "f1",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}
		]
	}`
	var out ListClusterTemplatesOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.ClusterTemplates) != 1 {
		t.Fatalf("ClusterTemplates: got %d, want 1", len(out.ClusterTemplates))
	}
	if out.ClusterTemplates[0].Name != "iaas_console" {
		t.Errorf("ClusterTemplates[0].Name: got %q", out.ClusterTemplates[0].Name)
	}
}

// ---------------------------------------------------------------------------
// GetKubeconfigOutput
// ---------------------------------------------------------------------------

func TestGetKubeconfigOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(GetKubeconfigOutput{}), []string{
		"Kubeconfig",
	})
}

func TestGetKubeconfigOutput_Value(t *testing.T) {
	out := GetKubeconfigOutput{Kubeconfig: "apiVersion: v1\nclusters: []"}
	if out.Kubeconfig == "" {
		t.Error("Kubeconfig should not be empty")
	}
}

// ---------------------------------------------------------------------------
// GetSupportedVersionsOutput
// ---------------------------------------------------------------------------

func TestGetSupportedVersionsOutput_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(GetSupportedVersionsOutput{}), []string{
		"SupportedK8s",
	})
}

func TestGetSupportedVersionsOutput_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(GetSupportedVersionsOutput{})
	testutil.AssertStructHasJSONTag(t, typ, "SupportedK8s", "supported_k8s")
}

func TestGetSupportedVersionsOutput_Parse(t *testing.T) {
	raw := `{"supported_k8s":{"v1.27.3":true,"v1.28.3":true,"v1.29.3":true}}`
	var out GetSupportedVersionsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.SupportedK8s) != 3 {
		t.Fatalf("SupportedK8s: got %d entries, want 3", len(out.SupportedK8s))
	}
	if !out.SupportedK8s["v1.29.3"] {
		t.Error("v1.29.3 should be supported")
	}
}

// ---------------------------------------------------------------------------
// Round-trip tests for key request types
// ---------------------------------------------------------------------------

func TestCreateClusterInput_RoundTrip(t *testing.T) {
	input := CreateClusterInput{
		Name:              "round-trip-cluster",
		ClusterTemplateID: "iaas_console",
		KeyPair:           "my-keypair",
		NodeCount:         "2",
		FlavorID:          "f-uuid",
		NetworkID:         "n-uuid",
		SubnetID:          "s-uuid",
		Labels:            map[string]string{"kube_tag": "v1.29.3"},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got CreateClusterInput
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != input.Name {
		t.Errorf("Name mismatch: got %q, want %q", got.Name, input.Name)
	}
	if got.ClusterTemplateID != input.ClusterTemplateID {
		t.Errorf("ClusterTemplateID mismatch: got %q", got.ClusterTemplateID)
	}
	if got.Labels["kube_tag"] != "v1.29.3" {
		t.Errorf("Labels[kube_tag] mismatch: got %q", got.Labels["kube_tag"])
	}
}

func TestUpdateNodeGroupInput_RoundTrip(t *testing.T) {
	input := UpdateNodeGroupInput{
		NodeCount:    5,
		MinNodeCount: 2,
		MaxNodeCount: 20,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got UpdateNodeGroupInput
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.NodeCount != 5 {
		t.Errorf("NodeCount: got %d, want 5", got.NodeCount)
	}
	if got.MinNodeCount != 2 {
		t.Errorf("MinNodeCount: got %d, want 2", got.MinNodeCount)
	}
	if got.MaxNodeCount != 20 {
		t.Errorf("MaxNodeCount: got %d, want 20", got.MaxNodeCount)
	}
}
