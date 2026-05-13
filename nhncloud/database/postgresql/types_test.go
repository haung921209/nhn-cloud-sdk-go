package postgresql_test

// types_test.go — spec-compliance tests for database/postgresql types.
//
// Tests verify that Go struct field names, JSON tags, and validation constraints
// match the official RDS for PostgreSQL API v1.0 specification:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/
//
// Key PostgreSQL-specific differences from MySQL/MariaDB:
//   - Password: 4-16 chars (NOT 4-256)
//   - Port range: 5432-45432
//   - databaseName field is REQUIRED on create
//   - No authenticationPlugin field (no NATIVE/SHA256/ED25519 variants)
//   - Extensions and HBA rules are PostgreSQL-only features
//   - API version: v1.0 (not v4.0)

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/testutil"
)

// ---------------------------------------------------------------------------
// InstanceStatus enum values — v1.0 spec
// ---------------------------------------------------------------------------

func TestInstanceStatus_EnumValues(t *testing.T) {
	checks := []struct {
		got  postgresql.InstanceStatus
		want string
	}{
		{postgresql.InstanceStatusAvailable, "AVAILABLE"},
		{postgresql.InstanceStatusBeforeCreate, "BEFORE_CREATE"},
		{postgresql.InstanceStatusCreating, "CREATING"},
		{postgresql.InstanceStatusModifying, "MODIFYING"},
		{postgresql.InstanceStatusDeleting, "DELETING"},
		{postgresql.InstanceStatusFailed, "FAILED"},
		{postgresql.InstanceStatusFailToCreate, "FAIL_TO_CREATE"},
		{postgresql.InstanceStatusStopped, "STOPPED"},
		{postgresql.InstanceStatusStopping, "STOPPING"},
		{postgresql.InstanceStatusStarting, "STARTING"},
		{postgresql.InstanceStatusRestarting, "RESTARTING"},
		{postgresql.InstanceStatusBackingUp, "BACKING_UP"},
		{postgresql.InstanceStatusRestoring, "RESTORING"},
	}
	for _, tc := range checks {
		if string(tc.got) != tc.want {
			t.Errorf("InstanceStatus value: got %q, want %q", tc.got, tc.want)
		}
	}
}

func TestInstanceType_EnumValues(t *testing.T) {
	if string(postgresql.InstanceTypeMaster) != "MASTER" {
		t.Errorf("InstanceTypeMaster: got %q, want %q", postgresql.InstanceTypeMaster, "MASTER")
	}
	if string(postgresql.InstanceTypeReplica) != "REPLICA" {
		t.Errorf("InstanceTypeReplica: got %q, want %q", postgresql.InstanceTypeReplica, "REPLICA")
	}
}

// ---------------------------------------------------------------------------
// JSON tag tests — DatabaseInstance (response type)
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.DatabaseInstance{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"DBInstanceID", "dbInstanceId"},
		{"DBInstanceName", "dbInstanceName"},
		{"DBInstanceStatus", "dbInstanceStatus"},
		{"DBInstanceType", "dbInstanceType"},
		{"DBVersion", "dbVersion"},
		{"DBPort", "dbPort"},
		{"DBFlavorID", "dbFlavorId"},
		{"ParameterGroupID", "parameterGroupId"},
		{"CreatedAt", "createdAt"},
		{"UpdatedAt", "updatedAt"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestDatabaseInstance_RequiredFields(t *testing.T) {
	testutil.AssertAllRequiredFields(t, reflect.TypeOf(postgresql.DatabaseInstance{}), []string{
		"DBInstanceID",
		"DBInstanceName",
		"DBInstanceStatus",
		"DBInstanceType",
		"DBVersion",
		"DBPort",
		"DBFlavorID",
		"ParameterGroupID",
		"CreatedAt",
		"UpdatedAt",
	})
}

// ---------------------------------------------------------------------------
// CreateInstanceRequest JSON tag tests
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.CreateInstanceRequest{})
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
		{"DatabaseName", "databaseName"},
		{"Network", "network"},
		{"Storage", "storage"},
		{"Backup", "backup"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

// ---------------------------------------------------------------------------
// PostgreSQL-specific: databaseName is REQUIRED
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_DatabaseNameRequired(t *testing.T) {
	// Verify field exists with correct JSON tag
	typ := reflect.TypeOf(postgresql.CreateInstanceRequest{})
	field, ok := typ.FieldByName("DatabaseName")
	if !ok {
		t.Fatal("CreateInstanceRequest must have a DatabaseName field (required by PostgreSQL spec)")
	}
	tag := field.Tag.Get("json")
	tagName := strings.Split(tag, ",")[0]
	if tagName != "databaseName" {
		t.Errorf("DatabaseName json tag: got %q, want %q", tagName, "databaseName")
	}

	// Validate that missing DatabaseName triggers an error
	port := 5432
	req := &postgresql.CreateInstanceRequest{
		DBInstanceName:   "test",
		DBFlavorID:       "f-id",
		DBVersion:        "POSTGRESQL_V14_6",
		DBUserName:       "admin",
		DBPassword:       "pass1",
		DBPort:           &port,
		ParameterGroupID: "pg-id",
		// DatabaseName intentionally omitted
		Network: postgresql.CreateInstanceNetworkConfig{
			SubnetID:         "sn-id",
			AvailabilityZone: "kr-pub-a",
		},
		Storage: postgresql.CreateInstanceStorageConfig{
			StorageType: "General SSD",
			StorageSize: 20,
		},
		Backup: postgresql.CreateInstanceBackupConfig{
			BackupPeriod: 1,
			BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
				{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
			},
		},
	}
	if err := req.Validate(); err == nil {
		t.Error("missing DatabaseName should fail validation for PostgreSQL")
	}
}

// ---------------------------------------------------------------------------
// Password validation — PostgreSQL v1.0 spec: 4-16 chars (different from MySQL!)
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PasswordValidation_PostgreSQL(t *testing.T) {
	makeReq := func(password string) *postgresql.CreateInstanceRequest {
		port := 5432
		return &postgresql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "POSTGRESQL_V14_6",
			DBUserName:       "admin",
			DBPassword:       password,
			DBPort:           &port,
			DatabaseName:     "mydb",
			ParameterGroupID: "pg-id",
			Network: postgresql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: postgresql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: postgresql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
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

	// Max boundary: 16 chars — must pass
	if err := makeReq(strings.Repeat("x", 16)).Validate(); err != nil {
		t.Errorf("password of 16 chars should pass (max for PostgreSQL), got: %v", err)
	}

	// Above max: 17 chars — must fail
	// PostgreSQL password max is 16 chars (unlike MySQL/MariaDB v4.0 which allow 256)
	if err := makeReq(strings.Repeat("x", 17)).Validate(); err == nil {
		t.Error("password of 17 chars should fail for PostgreSQL (max is 16 per v1.0 spec)")
	}

	// Well above max: 256 chars — must fail
	if err := makeReq(strings.Repeat("x", 256)).Validate(); err == nil {
		t.Error("password of 256 chars should fail for PostgreSQL (max is 16, not 256)")
	}
}

// ---------------------------------------------------------------------------
// Port validation — PostgreSQL spec: 5432-45432
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_PortValidation(t *testing.T) {
	makeReq := func(port int) *postgresql.CreateInstanceRequest {
		return &postgresql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "POSTGRESQL_V14_6",
			DBUserName:       "admin",
			DBPassword:       "pass1",
			DBPort:           &port,
			DatabaseName:     "mydb",
			ParameterGroupID: "pg-id",
			Network: postgresql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: postgresql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: 20,
			},
			Backup: postgresql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
					{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
				},
			},
		}
	}

	// Valid: default PostgreSQL port
	if err := makeReq(5432).Validate(); err != nil {
		t.Errorf("port 5432 should be valid for PostgreSQL, got: %v", err)
	}

	// Valid: max port per spec
	if err := makeReq(45432).Validate(); err != nil {
		t.Errorf("port 45432 should be valid for PostgreSQL, got: %v", err)
	}

	// Invalid: below range
	if err := makeReq(5431).Validate(); err == nil {
		t.Error("port 5431 should fail (below 5432 per PostgreSQL spec)")
	}

	// Invalid: above range
	if err := makeReq(45433).Validate(); err == nil {
		t.Error("port 45433 should fail (above 45432 per PostgreSQL spec)")
	}

	// Invalid: MySQL/MariaDB default port — wrong for PostgreSQL
	if err := makeReq(3306).Validate(); err == nil {
		t.Error("port 3306 should fail for PostgreSQL (that is a MySQL/MariaDB port)")
	}
}

// ---------------------------------------------------------------------------
// Storage size validation — 20-2048 GB
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_StorageSizeValidation(t *testing.T) {
	makeReq := func(size int) *postgresql.CreateInstanceRequest {
		port := 5432
		return &postgresql.CreateInstanceRequest{
			DBInstanceName:   "test",
			DBFlavorID:       "f-id",
			DBVersion:        "POSTGRESQL_V14_6",
			DBUserName:       "admin",
			DBPassword:       "pass1",
			DBPort:           &port,
			DatabaseName:     "mydb",
			ParameterGroupID: "pg-id",
			Network: postgresql.CreateInstanceNetworkConfig{
				SubnetID:         "sn-id",
				AvailabilityZone: "kr-pub-a",
			},
			Storage: postgresql.CreateInstanceStorageConfig{
				StorageType: "General SSD",
				StorageSize: size,
			},
			Backup: postgresql.CreateInstanceBackupConfig{
				BackupPeriod: 1,
				BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
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
// Extension types — PostgreSQL-specific
// ---------------------------------------------------------------------------

func TestExtension_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.Extension{})
	testutil.AssertStructHasJSONTag(t, typ, "ExtensionID", "extensionId")
	testutil.AssertStructHasJSONTag(t, typ, "ExtensionName", "extensionName")
	testutil.AssertStructHasJSONTag(t, typ, "ExtensionStatus", "extensionStatus")
	testutil.AssertStructHasJSONTag(t, typ, "Databases", "databases")
}

func TestExtensionDatabase_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.ExtensionDatabase{})
	testutil.AssertStructHasJSONTag(t, typ, "DBInstanceGroupExtensionID", "dbInstanceGroupExtensionId")
	testutil.AssertStructHasJSONTag(t, typ, "DatabaseID", "databaseId")
	testutil.AssertStructHasJSONTag(t, typ, "DatabaseName", "databaseName")
}

func TestInstallExtensionRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.InstallExtensionRequest{})
	testutil.AssertStructHasJSONTag(t, typ, "DatabaseID", "databaseId")
	testutil.AssertStructHasJSONTag(t, typ, "SchemaName", "schemaName")
	testutil.AssertStructHasJSONTag(t, typ, "WithCascade", "withCascade")
}

// ---------------------------------------------------------------------------
// HBA rule types — PostgreSQL-specific
// ---------------------------------------------------------------------------

func TestCreateHBARuleRequest_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.CreateHBARuleRequest{})
	cases := []struct {
		field   string
		jsonTag string
	}{
		{"DatabaseApplyType", "databaseApplyType"},
		{"DBUserApplyType", "dbUserApplyType"},
		{"Address", "address"},
		{"AuthMethod", "authMethod"},
	}
	for _, tc := range cases {
		testutil.AssertStructHasJSONTag(t, typ, tc.field, tc.jsonTag)
	}
}

func TestHBARule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.HBARule{})
	testutil.AssertStructHasJSONTag(t, typ, "HBARuleID", "hbaRuleId")
	// HBARuleStatus not in SDK — KNOWN DRIFT: field missing from HBARule struct
	testutil.AssertStructHasJSONTag(t, typ, "Order", "order")
	testutil.AssertStructHasJSONTag(t, typ, "DatabaseApplyType", "databaseApplyType")
	testutil.AssertStructHasJSONTag(t, typ, "Address", "address")
	testutil.AssertStructHasJSONTag(t, typ, "AuthMethod", "authMethod")
}

// ---------------------------------------------------------------------------
// PostgreSQLResponse header wrapper
// ---------------------------------------------------------------------------

func TestPostgreSQLResponse_HeaderJSONTag(t *testing.T) {
	testutil.AssertStructHasJSONTag(t, reflect.TypeOf(postgresql.PostgreSQLResponse{}), "Header", "header")
}

// ---------------------------------------------------------------------------
// JSON serialization round-trip — DatabaseInstance
// ---------------------------------------------------------------------------

func TestDatabaseInstance_JSONRoundTrip(t *testing.T) {
	inst := postgresql.DatabaseInstance{
		DBInstanceID:     "uuid-001",
		DBInstanceName:   "my-postgres",
		DBInstanceStatus: postgresql.InstanceStatusAvailable,
		DBInstanceType:   postgresql.InstanceTypeMaster,
		DBVersion:        "POSTGRESQL_V14_6",
		DBPort:           5432,
		DBFlavorID:       "flavor-uuid",
		ParameterGroupID: "pg-uuid",
		CreatedAt:        "2026-03-25T00:00:00",
		UpdatedAt:        "2026-03-25T00:00:00",
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded postgresql.DatabaseInstance
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.DBInstanceID != inst.DBInstanceID {
		t.Errorf("DBInstanceID: got %q, want %q", decoded.DBInstanceID, inst.DBInstanceID)
	}
	if decoded.DBInstanceStatus != inst.DBInstanceStatus {
		t.Errorf("DBInstanceStatus: got %q, want %q", decoded.DBInstanceStatus, inst.DBInstanceStatus)
	}
	if decoded.DBInstanceType != inst.DBInstanceType {
		t.Errorf("DBInstanceType: got %q, want %q", decoded.DBInstanceType, inst.DBInstanceType)
	}
	if decoded.DBPort != inst.DBPort {
		t.Errorf("DBPort: got %d, want %d", decoded.DBPort, inst.DBPort)
	}
}

func TestDatabaseInstance_JSONFieldNames(t *testing.T) {
	inst := postgresql.DatabaseInstance{
		DBInstanceID:   "id-123",
		DBInstanceName: "db-name",
		DBFlavorID:     "flavor-1",
		DBInstanceType: postgresql.InstanceTypeMaster,
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	expectedKeys := []string{"dbInstanceId", "dbInstanceName", "dbFlavorId", "dbInstanceType"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q not found in marshaled output", key)
		}
	}
}

// ---------------------------------------------------------------------------
// PostgreSQL-specific: no TLSOption, no AuthenticationPlugin on instance create
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_NoMySQLOnlyFields(t *testing.T) {
	typ := reflect.TypeOf(postgresql.CreateInstanceRequest{})

	// TLSOption is MySQL-only
	if _, ok := typ.FieldByName("TLSOption"); ok {
		t.Error("PostgreSQL CreateInstanceRequest must not have TLSOption field")
	}

	// AuthenticationPlugin is MySQL/MariaDB-only
	if _, ok := typ.FieldByName("AuthenticationPlugin"); ok {
		t.Error("PostgreSQL CreateInstanceRequest must not have AuthenticationPlugin field")
	}
}

// ---------------------------------------------------------------------------
// HighAvailability type (PostgreSQL-specific fields)
// ---------------------------------------------------------------------------

func TestHighAvailability_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.DatabaseInstanceHA{})
	testutil.AssertStructHasJSONTag(t, typ, "Use", "use")
}

// ---------------------------------------------------------------------------
// BackupSchedule JSON tags
// ---------------------------------------------------------------------------

func TestBackupSchedule_JSONTags(t *testing.T) {
	typ := reflect.TypeOf(postgresql.BackupSchedule{})
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndBgnTime", "backupWndBgnTime")
	testutil.AssertStructHasJSONTag(t, typ, "BackupWndDuration", "backupWndDuration")
}

// ---------------------------------------------------------------------------
// CreateInstanceRequest required field validation
// ---------------------------------------------------------------------------

func TestCreateInstanceRequest_MissingRequiredFields(t *testing.T) {
	port := 5432
	validBackup := postgresql.CreateInstanceBackupConfig{
		BackupPeriod: 1,
		BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
			{BackupWndBgnTime: "00:00:00", BackupWndDuration: "ONE_HOUR"},
		},
	}
	validNetwork := postgresql.CreateInstanceNetworkConfig{SubnetID: "sn", AvailabilityZone: "az"}
	validStorage := postgresql.CreateInstanceStorageConfig{StorageType: "General SSD", StorageSize: 20}

	cases := []struct {
		name    string
		req     *postgresql.CreateInstanceRequest
		wantErr bool
	}{
		{
			name: "missing DBInstanceName",
			req: &postgresql.CreateInstanceRequest{
				DBFlavorID: "fid", DBVersion: "POSTGRESQL_V14_6", DBUserName: "u",
				DBPassword: "pass1", DBPort: &port, DatabaseName: "db",
				ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "missing DatabaseName (PostgreSQL-specific required field)",
			req: &postgresql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBVersion: "POSTGRESQL_V14_6",
				DBUserName: "u", DBPassword: "pass1", DBPort: &port,
				ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: true,
		},
		{
			name: "valid request",
			req: &postgresql.CreateInstanceRequest{
				DBInstanceName: "n", DBFlavorID: "fid", DBVersion: "POSTGRESQL_V14_6",
				DBUserName: "u", DBPassword: "pass1", DBPort: &port,
				DatabaseName: "db", ParameterGroupID: "pg", Network: validNetwork,
				Storage: validStorage, Backup: validBackup,
			},
			wantErr: false,
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
