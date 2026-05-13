package mysql_test

// types_test.go — spec-compliance tests for database/mysql types.
//
// Tests verify that Go struct field names, JSON tags, and validation constraints
// match the official RDS for MySQL API v4.0 specification:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v4.0/

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// API version tests
// ---------------------------------------------------------------------------

// TestAPIVersion_IsV4 checks that instance.go uses the v4.0 path, not v3.0.
// This is a documentation/comment guard; the actual path is verified via
// the HTTP request path tests in the client tests.
func TestAPIVersion_DatabaseInstances_PathPrefix(t *testing.T) {
	// The Create/Modify/Delete paths in instances.go must use /v4.0/ not /v3.0/
	// We indirectly verify this by checking Validate passes on valid v4.0 input
	// (password up to 256 chars) and that the request struct exists.
	port := 3306
	req := &mysql.CreateInstanceRequest{
		DBInstanceName:   "test-db",
		DBFlavorID:       "flavor-uuid",
		DBVersion:        "MYSQL_V8032",
		DBUserName:       "admin",
		DBPassword:       strings.Repeat("a", 256), // v4.0 max
		DBPort:           &port,
		ParameterGroupID: "pg-uuid",
		Network: mysql.CreateInstanceNetworkConfig{
			SubnetID:         "subnet-uuid",
			AvailabilityZone: "kr-pub-a",
		},
		Storage: mysql.CreateInstanceStorageConfig{
			StorageType: "General SSD",
			StorageSize: 20,
		},
		Backup: mysql.CreateInstanceBackupConfig{
			BackupPeriod: 1,
			BackupSchedules: []mysql.CreateInstanceBackupSchedule{
				{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
			},
		},
	}
	if err := req.Validate(); err != nil {
		t.Errorf("256-char password should be valid in v4.0 spec, got error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// JSON tag tests — DatabaseInstance (response type)
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
		{"Description", "description"},
		{"DBInstanceStatus", "dbInstanceStatus"},
		{"DBVersion", "dbVersion"},
		{"DBPort", "dbPort"},
		{"DBFlavorID", "dbFlavorId"},
		{"ParameterGroupID", "parameterGroupId"},
		{"SupportUpgrade", "supportUpgrade"},
		{"CreatedYmdt", "createdYmdt"},
		{"UpdatedYmdt", "updatedYmdt"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestDatabaseInstance_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(mysql.DatabaseInstance{}), []string{
		"DBInstanceID",
		"DBInstanceName",
		"DBInstanceStatus",
		"DBVersion",
		"DBPort",
		"DBFlavorID",
		"ParameterGroupID",
		"CreatedYmdt",
		"UpdatedYmdt",
	})
}

// ---------------------------------------------------------------------------
// JSON tag tests — nested types
// ---------------------------------------------------------------------------

func TestDatabaseInstanceNetwork_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstanceNetwork{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"SubnetID", "subnetId"},
		{"AvailabilityZone", "availabilityZone"},
		{"UsePublicAccess", "usePublicAccess"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestDatabaseInstanceStorage_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstanceStorage{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"StorageType", "storageType"},
		{"StorageSize", "storageSize"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestDatabaseInstanceHA_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.DatabaseInstanceHA{})
	testutil.AssertStructHasJSONTag(t, typ, "Use", "use")
}

// ---------------------------------------------------------------------------
// CreateInstanceRequest JSON tag tests
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.CreateInstanceRequest{})
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
		{"Network", "network"},
		{"Storage", "storage"},
		{"Backup", "backup"},
		{"AuthenticationPlugin", "authenticationPlugin"},
		{"TLSOption", "tlsOption"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

// ---------------------------------------------------------------------------
// InstanceStatus enum values — v4.0 spec
// ---------------------------------------------------------------------------

func TestInstanceStatus_EnumValues(t *testing.T) {
	// v4.0 spec documents these status values for dbInstanceStatus
	expected := map[mysql.InstanceStatus]bool{
		mysql.InstanceStatusAvailable:       true,
		mysql.InstanceStatusBeforeCreate:    true,
		mysql.InstanceStatusStorageFull:     true,
		mysql.InstanceStatusFailToCreate:    true,
		mysql.InstanceStatusFailToConnect:   true,
		mysql.InstanceStatusReplicationStop: true,
		mysql.InstanceStatusFailover:        true,
		mysql.InstanceStatusShutdown:        true,
		mysql.InstanceStatusDeleted:         true,
	}
	for status := range expected {
		if string(status) == "" {
			t.Errorf("InstanceStatus constant has empty string value")
		}
	}

	// Spot-check specific string values match spec exactly
	checks := []struct {
		got  mysql.InstanceStatus
		want string
	}{
		{mysql.InstanceStatusAvailable, "AVAILABLE"},
		{mysql.InstanceStatusBeforeCreate, "BEFORE_CREATE"},
		{mysql.InstanceStatusStorageFull, "STORAGE_FULL"},
		{mysql.InstanceStatusFailToCreate, "FAIL_TO_CREATE"},
		{mysql.InstanceStatusFailToConnect, "FAIL_TO_CONNECT"},
		{mysql.InstanceStatusReplicationStop, "REPLICATION_STOP"},
		{mysql.InstanceStatusFailover, "FAILOVER"},
		{mysql.InstanceStatusShutdown, "SHUTDOWN"},
		{mysql.InstanceStatusDeleted, "DELETED"},
	}
	for _, tc := range checks {
		if string(tc.got) != tc.want {
			t.Errorf("InstanceStatus value: got %q, want %q", tc.got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Password validation — v4.0 spec: 4-256 characters
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PasswordValidation_V4Spec(t *testing.T) {
	makeReq := func(password string) *mysql.CreateInstanceRequest {
		port := 3306
		return &mysql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MYSQL_V8032",
			DBUserName:       "admin",
			DBPassword:       password,
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mysql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mysql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: mysql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mysql.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	// Too short: 3 chars — must fail
	if err := makeReq("abc").Validate(); err == nil {
		t.Error("password of 3 chars should fail validation (min is 4)")
	}

	// Min boundary: 4 chars — must pass
	if err := makeReq("abcd").Validate(); err != nil {
		t.Errorf("password of 4 chars should pass validation, got: %v", err)
	}

	// Mid-range: 16 chars — must pass (not artificially capped at 16 like old v3 spec)
	if err := makeReq(strings.Repeat("x", 16)).Validate(); err != nil {
		t.Errorf("password of 16 chars should pass validation, got: %v", err)
	}

	// v4.0 extended max: 100 chars — must pass
	if err := makeReq(strings.Repeat("x", 100)).Validate(); err != nil {
		t.Errorf("password of 100 chars should pass in v4.0 (max 256), got: %v", err)
	}

	// v4.0 max: 256 chars — must pass
	if err := makeReq(strings.Repeat("x", 256)).Validate(); err != nil {
		t.Errorf("password of 256 chars should pass in v4.0 spec, got: %v", err)
	}

	// Exceeds max: 257 chars — must fail
	if err := makeReq(strings.Repeat("x", 257)).Validate(); err == nil {
		t.Error("password of 257 chars should fail validation (max is 256 per v4.0)")
	}
}

// ---------------------------------------------------------------------------
// Port validation — v4.0 spec: 3306-43306
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PortValidation(t *testing.T) {
	makeReq := func(port int) *mysql.CreateInstanceRequest {
		return &mysql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MYSQL_V8032",
			DBUserName:       "admin",
			DBPassword:       "password1",
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mysql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mysql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: mysql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mysql.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	// Valid: default MySQL port
	if err := makeReq(3306).Validate(); err != nil {
		t.Errorf("port 3306 should be valid, got: %v", err)
	}

	// Valid: max port per spec
	if err := makeReq(43306).Validate(); err != nil {
		t.Errorf("port 43306 should be valid, got: %v", err)
	}

	// Invalid: below range
	if err := makeReq(3305).Validate(); err == nil {
		t.Error("port 3305 should fail (below 3306)")
	}

	// Invalid: above range
	if err := makeReq(43307).Validate(); err == nil {
		t.Error("port 43307 should fail (above 43306)")
	}
}

// ---------------------------------------------------------------------------
// Storage size validation — 20-2048 GB
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_StorageSizeValidation(t *testing.T) {
	makeReq := func(size int) *mysql.CreateInstanceRequest {
		port := 3306
		return &mysql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MYSQL_V8032",
			DBUserName:       "admin",
			DBPassword:       "password1",
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mysql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mysql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: size,
			},
			Backup: mysql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mysql.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	if err := makeReq(20).Validate(); err != nil {
		t.Errorf("storage 20 GB should be valid, got: %v", err)
	}
	if err := makeReq(2048).Validate(); err != nil {
		t.Errorf("storage 2048 GB should be valid, got: %v", err)
	}
	if err := makeReq(19).Validate(); err == nil {
		t.Error("storage 19 GB should fail (below minimum 20)")
	}
	if err := makeReq(2049).Validate(); err == nil {
		t.Error("storage 2049 GB should fail (above maximum 2048)")
	}
}

// ---------------------------------------------------------------------------
// Backup period validation — 0-730 days
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_BackupPeriodValidation(t *testing.T) {
	makeReq := func(period int) *mysql.CreateInstanceRequest {
		port := 3306
		return &mysql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MYSQL_V8032",
			DBUserName:       "admin",
			DBPassword:       "password1",
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mysql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mysql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: mysql.CreateInstanceBackupConfig{
				BackupPeriod: period,
				BackupSchedules: []mysql.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	if err := makeReq(0).Validate(); err != nil {
		t.Errorf("backup period 0 should be valid, got: %v", err)
	}
	if err := makeReq(730).Validate(); err != nil {
		t.Errorf("backup period 730 should be valid, got: %v", err)
	}
	if err := makeReq(731).Validate(); err == nil {
		t.Error("backup period 731 should fail (above max 730)")
	}
}

// ---------------------------------------------------------------------------
// JSON serialization round-trip
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONRoundTrip(t *testing.T) {
	inst := mysql.DatabaseInstance{
		DBInstanceID:     "uuid-001",
		DBInstanceName:   "my-mysql",
		Description:      "test instance",
		DBInstanceStatus: mysql.InstanceStatusAvailable,
		DBVersion:        "MYSQL_V8032",
		DBPort:           3306,
		DBFlavorID:       "flavor-uuid",
		ParameterGroupID: "pg-uuid",
		CreatedYmdt:      "2026-03-25T00:00:00",
		UpdatedYmdt:      "2026-03-25T00:00:00",
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded mysql.DatabaseInstance
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.DBInstanceID != inst.DBInstanceID {
		t.Errorf("DBInstanceID: got %q, want %q", decoded.DBInstanceID, inst.DBInstanceID)
	}
	if decoded.DBInstanceStatus != inst.DBInstanceStatus {
		t.Errorf("DBInstanceStatus: got %q, want %q", decoded.DBInstanceStatus, inst.DBInstanceStatus)
	}
	if decoded.DBPort != inst.DBPort {
		t.Errorf("DBPort: got %d, want %d", decoded.DBPort, inst.DBPort)
	}
}

func TestDatabaseInstance_JSONFieldNames(t *testing.T) {
	inst := mysql.DatabaseInstance{
		DBInstanceID:   "id-123",
		DBInstanceName: "db-name",
		DBFlavorID:     "flavor-1",
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Verify camelCase field names match spec
	expectedKeys := []string{"dbInstanceId", "dbInstanceName", "dbFlavorId"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q not found in marshaled output", key)
		}
	}

	// Verify old-style names are NOT present
	forbiddenKeys := []string{"DBInstanceId", "name", "id", "status"}
	for _, key := range forbiddenKeys {
		if _, ok := raw[key]; ok {
			t.Errorf("unexpected JSON key %q found — should use camelCase per API spec", key)
		}
	}
}

// ---------------------------------------------------------------------------
// MySQLResponse header wrapper
// ---------------------------------------------------------------------------

func TestMySQLResponse_HeaderJSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(mysql.MySQLResponse{}), "Header", "header")
}

// ---------------------------------------------------------------------------
// CreateInstanceRequest required field validation
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_MissingRequiredFields(t *testing.T) {
	port := 3306
	validBackup := mysql.CreateInstanceBackupConfig{
		BackupPeriod: 1,
		BackupSchedules: []mysql.CreateInstanceBackupSchedule{
			{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
		},
	}
	validNetwork := mysql.CreateInstanceNetworkConfig{SubnetID: "sn", AvailabilityZone: "az"}
	validStorage := mysql.CreateInstanceStorageConfig{StorageType: "General SSD", StorageSize: 20}

	cases := []struct {
		name    string
		req     *mysql.CreateInstanceRequest
		wantErr bool
	}{
		{
			name: "missing DBInstanceName",
			req: &mysql.CreateInstanceRequest{
				DBFlavorID: "fid", DBVersion: "MYSQL_V8032", DBUserName: "u", DBPassword: "pass1",
				DBPort: &port, ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing DBFlavorID",
			req: &mysql.CreateInstanceRequest{
				DBInstanceName: "n", DBVersion: "MYSQL_V8032", DBUserName: "u", DBPassword: "pass1",
				DBPort: &port, ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing DBVersion",
			req: &mysql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBUserName: "u", DBPassword: "pass1",
				DBPort: &port, ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing DBUserName",
			req: &mysql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBVersion: "MYSQL_V8032", DBPassword: "pass1",
				DBPort: &port, ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing DBPassword",
			req: &mysql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBVersion: "MYSQL_V8032", DBUserName: "u",
				DBPort: &port, ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing ParameterGroupID",
			req: &mysql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBVersion: "MYSQL_V8032",
				DBUserName: "u", DBPassword: "pass1", DBPort: &port,
				Network: validNetwork, Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BackupSchedule JSON tags
// ---------------------------------------------------------------------------

func TestBackupSchedule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mysql.BackupSchedule{})
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndBgnTime", "backupWndBgnTime")
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndDuration", "backupWndDuration")
}
