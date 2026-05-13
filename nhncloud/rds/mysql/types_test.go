package mysql_test

// types_test.go — spec-compliance tests for rds/mysql types.
//
// The rds/mysql package provides the lower-level (legacy-compatible) types
// used by the CLI layer. These tests verify JSON tag correctness and
// response parsing of nested objects (storage, network).

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mysql"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// ResponseHeader
// ---------------------------------------------------------------------------

func TestResponseHeader_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.ResponseHeader{})
	testutil.AssertStructHasJSONTag(t, typ, "ResultCode", "resultCode")
	testutil.AssertStructHasJSONTag(t, typ, "ResultMessage", "resultMessage")
	testutil.AssertStructHasJSONTag(t, typ, "IsSuccessful", "isSuccessful")
}

// ---------------------------------------------------------------------------
// DatabaseInstance — core response fields
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstance{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"DBInstanceID", "dbInstanceId"},
		{"DBInstanceGroupID", "dbInstanceGroupId"},
		{"DBInstanceName", "dbInstanceName"},
		{"DBInstanceStatus", "dbInstanceStatus"},
		{"DBInstanceType", "dbInstanceType"},
		{"DBVersion", "dbVersion"},
		{"DBPort", "dbPort"},
		{"ProgressStatus", "progressStatus"},
		{"Storage", "storage"},
		{"Network", "network"},
		{"DBFlavorID", "dbFlavorId"},
		{"AuthenticationPlugin", "authenticationPlugin"},
		{"TLSOption", "tlsOption"},
		{"UseDeletionProtection", "useDeletionProtection"},
		{"CreatedYmdt", "createdYmdt"},
		{"UpdatedYmdt", "updatedYmdt"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestInstanceStorage_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.InstanceStorage{})
	testutil.AssertStructHasJSONTag(t, typ, "StorageType", "storageType")
	testutil.AssertStructHasJSONTag(t, typ, "StorageSize", "storageSize")
}

func TestInstanceNetwork_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.InstanceNetwork{})
	testutil.AssertStructHasJSONTag(t, typ, "SubnetID", "subnetId")
}

func TestDatabaseInstance_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(mysql.DatabaseInstance{}), []string{
		"DBInstanceID",
		"DBInstanceName",
		"DBInstanceStatus",
		"DBVersion",
		"DBPort",
		"DBFlavorID",
		"CreatedYmdt",
		"UpdatedYmdt",
	})
}

// ---------------------------------------------------------------------------
// DatabaseInstancesResponse — list response parsing
// ---------------------------------------------------------------------------

func TestDatabaseInstancesResponse_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstancesResponse{})
	testutil.AssertStructHasJSONTag(t, typ, "DBInstances", "dbInstances")
}

func TestDatabaseInstancesResponse_ParseNestedObjects(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"dbInstances": [
			{
				"dbInstanceId": "inst-001",
				"dbInstanceGroupId": "grp-001",
				"dbInstanceName": "my-db",
				"dbInstanceStatus": "AVAILABLE",
				"dbInstanceType": "MASTER",
				"dbVersion": "MYSQL_V8032",
				"dbPort": 3306,
				"storage": {"storageType": "General SSD", "storageSize": 100},
				"network": {"subnetId": "subnet-abc"},
				"dbFlavorId": "flavor-xyz",
				"authenticationPlugin": "CACHING_SHA2",
				"tlsOption": "NONE",
				"useDeletionProtection": false,
				"progressStatus": "NONE",
				"createdYmdt": "2026-03-25T00:00:00",
				"updatedYmdt": "2026-03-25T00:00:00"
			}
		]
	}`

	var resp mysql.DatabaseInstancesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if resp.Header == nil {
		t.Fatal("Header should not be nil")
	}
	if !resp.Header.IsSuccessful {
		t.Error("Header.IsSuccessful should be true")
	}
	if len(resp.DBInstances) != 1 {
		t.Fatalf("DBInstances: got %d items, want 1", len(resp.DBInstances))
	}

	inst := resp.DBInstances[0]
	if inst.DBInstanceID != "inst-001" {
		t.Errorf("DBInstanceID: got %q, want %q", inst.DBInstanceID, "inst-001")
	}
	if inst.DBInstanceGroupID != "grp-001" {
		t.Errorf("DBInstanceGroupID: got %q, want %q", inst.DBInstanceGroupID, "grp-001")
	}
	if inst.DBInstanceName != "my-db" {
		t.Errorf("DBInstanceName: got %q, want %q", inst.DBInstanceName, "my-db")
	}
	if inst.DBInstanceStatus != "AVAILABLE" {
		t.Errorf("DBInstanceStatus: got %q, want %q", inst.DBInstanceStatus, "AVAILABLE")
	}
	if inst.DBInstanceType != "MASTER" {
		t.Errorf("DBInstanceType: got %q, want %q", inst.DBInstanceType, "MASTER")
	}
	if inst.DBPort != 3306 {
		t.Errorf("DBPort: got %d, want %d", inst.DBPort, 3306)
	}
	if inst.Storage == nil {
		t.Fatal("Storage should not be nil")
	}
	if inst.Storage.StorageType != "General SSD" {
		t.Errorf("Storage.StorageType: got %q, want %q", inst.Storage.StorageType, "General SSD")
	}
	if inst.Storage.StorageSize != 100 {
		t.Errorf("Storage.StorageSize: got %d, want %d", inst.Storage.StorageSize, 100)
	}
	if inst.Network == nil {
		t.Fatal("Network should not be nil")
	}
	if inst.Network.SubnetID != "subnet-abc" {
		t.Errorf("Network.SubnetID: got %q, want %q", inst.Network.SubnetID, "subnet-abc")
	}
	if inst.AuthenticationPlugin != "CACHING_SHA2" {
		t.Errorf("AuthenticationPlugin: got %q, want %q", inst.AuthenticationPlugin, "CACHING_SHA2")
	}
	if inst.ProgressStatus != "NONE" {
		t.Errorf("ProgressStatus: got %q, want %q", inst.ProgressStatus, "NONE")
	}
}

// ---------------------------------------------------------------------------
// DBFlavor — flavor response parsing
// ---------------------------------------------------------------------------

func TestDBFlavor_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DBFlavor{})
	testutil.AssertStructHasJSONTag(t, typ, "FlavorID", "dbFlavorId")
	testutil.AssertStructHasJSONTag(t, typ, "FlavorName", "dbFlavorName")
	testutil.AssertStructHasJSONTag(t, typ, "Ram", "ram")
	testutil.AssertStructHasJSONTag(t, typ, "Vcpus", "vcpus")
}

func TestDBFlavorsResponse_ParseNestedObjects(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"dbFlavors": [
			{
				"dbFlavorId": "flavor-001",
				"dbFlavorName": "m2.c1m2",
				"ram": 2048,
				"vcpus": 1,
				"disk": 20
			}
		]
	}`

	var resp mysql.DBFlavorsResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if len(resp.DBFlavors) != 1 {
		t.Fatalf("DBFlavors: got %d items, want 1", len(resp.DBFlavors))
	}
	f := resp.DBFlavors[0]
	if f.FlavorID != "flavor-001" {
		t.Errorf("FlavorID: got %q, want %q", f.FlavorID, "flavor-001")
	}
	if f.FlavorName != "m2.c1m2" {
		t.Errorf("FlavorName: got %q, want %q", f.FlavorName, "m2.c1m2")
	}
	if f.Ram != 2048 {
		t.Errorf("Ram: got %d, want 2048", f.Ram)
	}
	if f.Vcpus != 1 {
		t.Errorf("Vcpus: got %d, want 1", f.Vcpus)
	}
}

// ---------------------------------------------------------------------------
// Network — NetworkSubnet parsing
// ---------------------------------------------------------------------------

func TestNetworkSubnet_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.NetworkSubnet{})
	testutil.AssertStructHasJSONTag(t, typ, "SubnetID", "subnetId")
	testutil.AssertStructHasJSONTag(t, typ, "SubnetName", "subnetName")
	testutil.AssertStructHasJSONTag(t, typ, "SubnetCidr", "subnetCidr")
	testutil.AssertStructHasJSONTag(t, typ, "UsingGateway", "usingGateway")
	testutil.AssertStructHasJSONTag(t, typ, "AvailableIpCount", "availableIpCount")
}

func TestNetworkSubnetsResponse_ParseNestedObjects(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"subnets": [
			{
				"subnetId": "sn-001",
				"subnetName": "my-subnet",
				"subnetCidr": "192.168.0.0/24",
				"usingGateway": true,
				"availableIpCount": 250
			}
		]
	}`

	var resp mysql.NetworkSubnetsResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if len(resp.Subnets) != 1 {
		t.Fatalf("Subnets: got %d items, want 1", len(resp.Subnets))
	}
	sn := resp.Subnets[0]
	if sn.SubnetID != "sn-001" {
		t.Errorf("SubnetID: got %q, want %q", sn.SubnetID, "sn-001")
	}
	if sn.SubnetCidr != "192.168.0.0/24" {
		t.Errorf("SubnetCidr: got %q, want %q", sn.SubnetCidr, "192.168.0.0/24")
	}
	if sn.AvailableIpCount != 250 {
		t.Errorf("AvailableIpCount: got %d, want 250", sn.AvailableIpCount)
	}
}

// ---------------------------------------------------------------------------
// DBSecurityGroup
// ---------------------------------------------------------------------------

func TestDBSecurityGroup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DBSecurityGroup{})
	testutil.AssertStructHasJSONTag(t, typ, "DBSecurityGroupID", "dbSecurityGroupId")
	testutil.AssertStructHasJSONTag(t, typ, "DBSecurityGroupName", "dbSecurityGroupName")
	testutil.AssertStructHasJSONTag(t, typ, "ProgressStatus", "progressStatus")
}

// ---------------------------------------------------------------------------
// ParameterGroup
// ---------------------------------------------------------------------------

func TestParameterGroup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.ParameterGroup{})
	testutil.AssertStructHasJSONTag(t, typ, "ParameterGroupID", "parameterGroupId")
	testutil.AssertStructHasJSONTag(t, typ, "ParameterGroupName", "parameterGroupName")
	testutil.AssertStructHasJSONTag(t, typ, "DBVersion", "dbVersion")
	testutil.AssertStructHasJSONTag(t, typ, "ParameterGroupStatus", "parameterGroupStatus")
}

// ---------------------------------------------------------------------------
// Backup
// ---------------------------------------------------------------------------

func TestBackup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.Backup{})
	testutil.AssertStructHasJSONTag(t, typ, "BackupID", "backupId")
	testutil.AssertStructHasJSONTag(t, typ, "BackupName", "backupName")
	testutil.AssertStructHasJSONTag(t, typ, "BackupStatus", "backupStatus")
	testutil.AssertStructHasJSONTag(t, typ, "DBInstanceID", "dbInstanceId")
	testutil.AssertStructHasJSONTag(t, typ, "BackupType", "backupType")
	testutil.AssertStructHasJSONTag(t, typ, "BackupSize", "backupSize")
}

// ---------------------------------------------------------------------------
// DBUser
// ---------------------------------------------------------------------------

func TestDBUser_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DBUser{})
	testutil.AssertStructHasJSONTag(t, typ, "DBUserID", "dbUserId")
	testutil.AssertStructHasJSONTag(t, typ, "DBUserName", "dbUserName")
	testutil.AssertStructHasJSONTag(t, typ, "Host", "host")
	testutil.AssertStructHasJSONTag(t, typ, "AuthorityType", "authorityType")
	testutil.AssertStructHasJSONTag(t, typ, "DBUserStatus", "dbUserStatus")
	testutil.AssertStructHasJSONTag(t, typ, "AuthenticationPlugin", "authenticationPlugin")
	testutil.AssertStructHasJSONTag(t, typ, "TLSOption", "tlsOption")
}

// ---------------------------------------------------------------------------
// JobIDResponse
// ---------------------------------------------------------------------------

func TestJobIDResponse_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.JobIDResponse{})
	testutil.AssertStructHasJSONTag(t, typ, "JobID", "jobId")
}

func TestJobIDResponse_ParseJSON(t *testing.T) {
	raw := `{"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true}, "jobId": "job-abc-123"}`
	var resp mysql.JobIDResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if resp.JobID != "job-abc-123" {
		t.Errorf("JobID: got %q, want %q", resp.JobID, "job-abc-123")
	}
}

// ---------------------------------------------------------------------------
// CreateDatabaseInstanceRequest required fields
// ---------------------------------------------------------------------------

func TestCreateDatabaseInstanceRequest_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(mysql.CreateDatabaseInstanceRequest{}), []string{
		"DBInstanceName",
		"DBFlavorID",
		"DBVersion",
		"DBUserName",
		"DBPassword",
		"ParameterGroupID",
		"Network",
		"Storage",
		"Backup",
	})
}

func TestCreateDatabaseInstanceRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.CreateDatabaseInstanceRequest{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"DBInstanceName", "dbInstanceName"},
		{"DBFlavorID", "dbFlavorId"},
		{"DBVersion", "dbVersion"},
		{"DBUserName", "dbUserName"},
		{"DBPassword", "dbPassword"},
		{"DBPort", "dbPort"},
		{"ParameterGroupID", "parameterGroupId"},
		{"DBSecurityGroupIDs", "dbSecurityGroupIds"},
		{"UserGroupIDs", "userGroupIds"},
		{"NotificationGroupIDs", "notificationGroupIds"},
		{"AuthenticationPlugin", "authenticationPlugin"},
		{"TLSOption", "tlsOption"},
		{"UseHighAvailability", "useHighAvailability"},
		{"UseDeletionProtection", "useDeletionProtection"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

// ---------------------------------------------------------------------------
// DatabaseInstanceGroup
// ---------------------------------------------------------------------------

func TestDatabaseInstanceGroup_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstanceGroup{})
	testutil.AssertStructHasJSONTag(t, typ, "DBInstanceGroupID", "dbInstanceGroupId")
	testutil.AssertStructHasJSONTag(t, typ, "ReplicationType", "replicationType")
	testutil.AssertStructHasJSONTag(t, typ, "CreatedYmdt", "createdYmdt")
	testutil.AssertStructHasJSONTag(t, typ, "UpdatedYmdt", "updatedYmdt")
}

// ---------------------------------------------------------------------------
// NetworkInfoResponse — nested endpoint parsing
// ---------------------------------------------------------------------------

func TestNetworkInfoResponse_ParseNestedEndpoints(t *testing.T) {
	raw := `{
		"header": {"resultCode": 0, "resultMessage": "SUCCESS", "isSuccessful": true},
		"availabilityZone": "kr-pub-a",
		"subnet": {
			"subnetId": "sn-001",
			"subnetName": "my-subnet",
			"subnetCidr": "10.0.0.0/24"
		},
		"endPoints": [
			{"domain": "db.example.com", "ipAddress": "10.0.0.5", "endPointType": "PRIVATE"},
			{"domain": "db-pub.example.com", "ipAddress": "1.2.3.4", "endPointType": "PUBLIC"}
		]
	}`

	var resp mysql.NetworkInfoResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if resp.AvailabilityZone != "kr-pub-a" {
		t.Errorf("AvailabilityZone: got %q, want %q", resp.AvailabilityZone, "kr-pub-a")
	}
	if resp.Subnet.SubnetID != "sn-001" {
		t.Errorf("Subnet.SubnetID: got %q, want %q", resp.Subnet.SubnetID, "sn-001")
	}
	if len(resp.EndPoints) != 2 {
		t.Fatalf("EndPoints: got %d items, want 2", len(resp.EndPoints))
	}
	if resp.EndPoints[0].Domain != "db.example.com" {
		t.Errorf("EndPoints[0].Domain: got %q, want %q", resp.EndPoints[0].Domain, "db.example.com")
	}
	if resp.EndPoints[1].EndPointType != "PUBLIC" {
		t.Errorf("EndPoints[1].EndPointType: got %q, want %q", resp.EndPoints[1].EndPointType, "PUBLIC")
	}
}

// ---------------------------------------------------------------------------
// Type alias sanity checks — ensure aliases point to correct base types
// ---------------------------------------------------------------------------

func TestTypeAliases_ListInstances(t *testing.T) {
	// ListInstancesOutput must be DatabaseInstancesResponse
	var _ mysql.ListInstancesOutput = mysql.DatabaseInstancesResponse{}
}

func TestTypeAliases_ListFlavors(t *testing.T) {
	var _ mysql.ListFlavorsOutput = mysql.DBFlavorsResponse{}
}

func TestTypeAliases_CreateInstance(t *testing.T) {
	var _ mysql.CreateInstanceInput = mysql.CreateDatabaseInstanceRequest{}
}

func TestTypeAliases_JobOutput(t *testing.T) {
	var _ mysql.JobOutput = mysql.JobIDResponse{}
}
