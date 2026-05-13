package mariadb_test

// types_test.go — spec-compliance tests for database/mariadb types.
//
// Tests verify that Go struct field names, JSON tags, and validation constraints
// match the official RDS for MariaDB API v4.0 specification:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v4.0/

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// InstanceStatus enum values — v4.0 spec
// ---------------------------------------------------------------------------

func TestInstanceStatus_EnumValues(t *testing.T) {
	checks := []struct {
		got  mariadb.InstanceStatus
		want string
	}{
		{mariadb.InstanceStatusAvailable, "AVAILABLE"},
		{mariadb.InstanceStatusBeforeCreate, "BEFORE_CREATE"},
		{mariadb.InstanceStatusCreating, "CREATING"},
		{mariadb.InstanceStatusModifying, "MODIFYING"},
		{mariadb.InstanceStatusDeleting, "DELETING"},
		{mariadb.InstanceStatusFailed, "FAILED"},
		{mariadb.InstanceStatusFailToCreate, "FAIL_TO_CREATE"},
		{mariadb.InstanceStatusStopped, "STOPPED"},
		{mariadb.InstanceStatusStopping, "STOPPING"},
		{mariadb.InstanceStatusStarting, "STARTING"},
		{mariadb.InstanceStatusRestarting, "RESTARTING"},
		{mariadb.InstanceStatusBackingUp, "BACKING_UP"},
		{mariadb.InstanceStatusRestoring, "RESTORING"},
		{mariadb.InstanceStatusReplicating, "REPLICATING"},
		{mariadb.InstanceStatusFailoverIng, "FAILOVER_ING"},
	}
	for _, tc := range checks {
		if string(tc.got) != tc.want {
			t.Errorf("InstanceStatus value: got %q, want %q", tc.got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// JSON tag tests — DatabaseInstance (response type)
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mariadb.DatabaseInstance{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"DBInstanceID", "dbInstanceId"},
		{"DBInstanceName", "dbInstanceName"},
		{"DBInstanceDescription", "description"},
		{"DBInstanceStatus", "dbInstanceStatus"},
		{"DBVersion", "dbVersion"},
		{"DBPort", "dbPort"},
		{"DBFlavorID", "dbFlavorId"},
		{"ParameterGroupID", "parameterGroupId"},
		{"CreatedAt", "createdYmdt"},
		{"UpdatedAt", "updatedYmdt"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestDatabaseInstance_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(mariadb.DatabaseInstance{}), []string{
		"DBInstanceID",
		"DBInstanceName",
		"DBInstanceStatus",
		"DBVersion",
		"DBPort",
		"DBFlavorID",
		"ParameterGroupID",
		"CreatedAt",
		"UpdatedAt",
	})
}

// ---------------------------------------------------------------------------
// JSON tag tests — nested types
// ---------------------------------------------------------------------------

func TestDatabaseInstanceNetwork_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mariadb.DatabaseInstanceNetwork{})
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
	typ := reflect.TypeOf(mariadb.DatabaseInstanceStorage{})
	testutil.AssertStructHasJSONTag(t, typ, "StorageType", "storageType")
	testutil.AssertStructHasJSONTag(t, typ, "StorageSize", "storageSize")
}

// ---------------------------------------------------------------------------
// CreateInstanceRequest JSON tag tests
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mariadb.CreateInstanceRequest{})
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
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

// ---------------------------------------------------------------------------
// Auth plugins — MariaDB v4.0 spec: NATIVE, ED25519
// (NOT SHA256/CACHING_SHA2 which are MySQL-only)
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_AuthPlugin_MariaDB_NoTLSOption(t *testing.T) {
	// MariaDB spec does NOT include TLSOption; the field must not be present
	// in CreateInstanceRequest for MariaDB.
	typ := reflect.TypeOf(mariadb.CreateInstanceRequest{})
	if _, ok := typ.FieldByName("TLSOption"); ok {
		t.Error("MariaDB CreateInstanceRequest must not have TLSOption field (MariaDB v4.0 spec does not include it)")
	}
}

func TestCreateInstanceRequest_AuthPlugin_JSONSerialization(t *testing.T) {
	port := 3306
	req := mariadb.CreateInstanceRequest{
		DBInstanceName:       "test-db",
		DBFlavorID:           "flavor-uuid",
		DBVersion:            "MARIADB_V10611",
		DBUserName:           "admin",
		DBPassword:           "password1",
		DBPort:               &port,
		ParameterGroupID:     "pg-uuid",
		AuthenticationPlugin: "NATIVE",
		Network: mariadb.CreateInstanceNetworkConfig{
			SubnetID:         "subnet-uuid",
			AvailabilityZone: "kr-pub-a",
		},
		Storage: mariadb.CreateInstanceStorageConfig{
			StorageType: "General SSD",
			StorageSize: 20,
		},
		Backup: mariadb.CreateInstanceBackupConfig{
			BackupPeriod: 1,
			BackupSchedules: []mariadb.CreateInstanceBackupSchedule{
				{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if v, ok := raw["authenticationPlugin"]; !ok || v != "NATIVE" {
		t.Errorf("authenticationPlugin: got %v, want %q", v, "NATIVE")
	}
}

// ---------------------------------------------------------------------------
// Password validation — v4.0 spec: 4-256 characters
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PasswordValidation_V4Spec(t *testing.T) {
	makeReq := func(password string) *mariadb.CreateInstanceRequest {
		port := 3306
		return &mariadb.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MARIADB_V10611",
			DBUserName:       "admin",
			DBPassword:       password,
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mariadb.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mariadb.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: mariadb.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mariadb.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	// Below min: 3 chars — must fail
	if err := makeReq("abc").Validate(); err == nil {
		t.Error("password of 3 chars should fail (min is 4)")
	}

	// Min boundary: 4 chars — must pass
	if err := makeReq("abcd").Validate(); err != nil {
		t.Errorf("password of 4 chars should pass, got: %v", err)
	}

	// Mid: 16 chars — must pass (not capped at 16 in v4.0)
	if err := makeReq(strings.Repeat("x", 16)).Validate(); err != nil {
		t.Errorf("password of 16 chars should pass, got: %v", err)
	}

	// v4.0 max: 256 chars — must pass
	if err := makeReq(strings.Repeat("x", 256)).Validate(); err != nil {
		t.Errorf("password of 256 chars should pass per v4.0 spec, got: %v", err)
	}

	// Above max: 257 chars — must fail
	if err := makeReq(strings.Repeat("x", 257)).Validate(); err == nil {
		t.Error("password of 257 chars should fail (max is 256 per v4.0)")
	}
}

// ---------------------------------------------------------------------------
// Port validation — MariaDB spec: 3306-43306
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PortValidation(t *testing.T) {
	makeReq := func(port int) *mariadb.CreateInstanceRequest {
		return &mariadb.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MARIADB_V10611",
			DBUserName:       "admin",
			DBPassword:       "password1",
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mariadb.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mariadb.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: mariadb.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mariadb.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	// Valid bounds
	if err := makeReq(3306).Validate(); err != nil {
		t.Errorf("port 3306 should be valid, got: %v", err)
	}
	if err := makeReq(43306).Validate(); err != nil {
		t.Errorf("port 43306 should be valid, got: %v", err)
	}

	// Invalid: below range
	if err := makeReq(3305).Validate(); err == nil {
		t.Error("port 3305 should fail (below 3306 per MariaDB spec)")
	}

	// Invalid: above range
	if err := makeReq(43307).Validate(); err == nil {
		t.Error("port 43307 should fail (above 43306 per MariaDB spec)")
	}

	// Note: 5432 is within 3306-43306 range, so it's technically valid for MariaDB
	// even though it's the PostgreSQL default port
	if err := makeReq(5432).Validate(); err != nil {
		t.Errorf("port 5432 is within 3306-43306 range and should be valid, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Storage size validation — 20-2048 GB
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_StorageSizeValidation(t *testing.T) {
	makeReq := func(size int) *mariadb.CreateInstanceRequest {
		port := 3306
		return &mariadb.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "MARIADB_V10611",
			DBUserName:       "admin",
			DBPassword:       "password1",
			DBPort:           &port,
			ParameterGroupID: "pg-id",
			Network: mariadb.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: mariadb.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: size,
			},
			Backup: mariadb.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []mariadb.CreateInstanceBackupSchedule{
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
		t.Error("storage 19 GB should fail (below 20)")
	}
	if err := makeReq(2049).Validate(); err == nil {
		t.Error("storage 2049 GB should fail (above 2048)")
	}
}

// ---------------------------------------------------------------------------
// JSON serialization round-trip — DatabaseInstance
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONRoundTrip(t *testing.T) {
	inst := mariadb.DatabaseInstance{
		DBInstanceID:          "uuid-001",
		DBInstanceName:        "my-mariadb",
		DBInstanceDescription: "test instance",
		DBInstanceStatus:      mariadb.InstanceStatusAvailable,
		DBVersion:             "MARIADB_V10611",
		DBPort:                3306,
		DBFlavorID:            "flavor-uuid",
		ParameterGroupID:      "pg-uuid",
		CreatedAt:             "2026-03-25T00:00:00",
		UpdatedAt:             "2026-03-25T00:00:00",
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded mariadb.DatabaseInstance
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

// ---------------------------------------------------------------------------
// MariaDBResponse header wrapper
// ---------------------------------------------------------------------------

func TestMariaDBResponse_HeaderJSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(mariadb.MariaDBResponse{}), "Header", "header")
}

// ---------------------------------------------------------------------------
// BackupSchedule JSON tags
// ---------------------------------------------------------------------------

func TestBackupSchedule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mariadb.BackupSchedule{})
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndBgnTime", "backupWndBgnTime")
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndDuration", "backupWndDuration")
}

// ---------------------------------------------------------------------------
// HA struct JSON tags
// ---------------------------------------------------------------------------

func TestDatabaseInstanceHA_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(mariadb.DatabaseInstanceHA{})
	testutil.AssertStructHasJSONTag(t, typ, "Use", "use")
}
