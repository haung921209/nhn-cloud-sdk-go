package nhncloud

import (
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

func TestNewClient(t *testing.T) {
	creds := credentials.NewStatic("access-key", "secret-key")
	cfg := &Config{
		Region:      "kr1",
		Credentials: creds,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestNewClientNilConfig(t *testing.T) {
	_, err := New(nil)
	if err != ErrCredentialsRequired {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestClientServiceAccessors(t *testing.T) {
	creds := credentials.NewStatic("access-key", "secret-key")
	identityCreds := credentials.NewStaticIdentity("user", "pass", "tenant")

	cfg := &Config{
		Region:              "kr1",
		Credentials:         creds,
		IdentityCredentials: identityCreds,
		AppKeys: map[string]string{
			"rds-mysql":      "mysql-appkey",
			"rds-mariadb":    "mariadb-appkey",
			"rds-postgresql": "pg-appkey",
			"ncr":            "ncr-appkey",
			"ncs":            "ncs-appkey",
		},
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.IAM() == nil {
		t.Error("IAM() returned nil")
	}

	if client.Compute() == nil {
		t.Error("Compute() returned nil")
	}

	if client.MySQL() == nil {
		t.Error("MySQL() returned nil")
	}

	if client.MariaDB() == nil {
		t.Error("MariaDB() returned nil")
	}

	if client.PostgreSQL() == nil {
		t.Error("PostgreSQL() returned nil")
	}

	if client.VPC() == nil {
		t.Error("VPC() returned nil")
	}

	if client.SecurityGroup() == nil {
		t.Error("SecurityGroup() returned nil")
	}

	if client.FloatingIP() == nil {
		t.Error("FloatingIP() returned nil")
	}

	if client.LoadBalancer() == nil {
		t.Error("LoadBalancer() returned nil")
	}

	if client.BlockStorage() == nil {
		t.Error("BlockStorage() returned nil")
	}

	if client.ObjectStorage() == nil {
		t.Error("ObjectStorage() returned nil")
	}

	if client.NKS() == nil {
		t.Error("NKS() returned nil")
	}

	if client.NCR() == nil {
		t.Error("NCR() returned nil")
	}

	if client.NCS() == nil {
		t.Error("NCS() returned nil")
	}
}

func TestClientLazyInitialization(t *testing.T) {
	creds := credentials.NewStatic("access-key", "secret-key")
	cfg := &Config{
		Region:      "kr1",
		Credentials: creds,
	}

	client, _ := New(cfg)

	iam1 := client.IAM()
	iam2 := client.IAM()

	if iam1 != iam2 {
		t.Error("IAM() should return same instance (lazy initialization)")
	}
}
