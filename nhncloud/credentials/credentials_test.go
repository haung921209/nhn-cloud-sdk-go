package credentials

import (
	"os"
	"testing"
)

func TestNewStatic(t *testing.T) {
	creds := NewStatic("access-key", "secret-key")

	if got := creds.GetAccessKeyID(); got != "access-key" {
		t.Errorf("GetAccessKeyID() = %v, want %v", got, "access-key")
	}

	if got := creds.GetSecretAccessKey(); got != "secret-key" {
		t.Errorf("GetSecretAccessKey() = %v, want %v", got, "secret-key")
	}
}

func TestNewStaticIdentity(t *testing.T) {
	creds := NewStaticIdentity("user@example.com", "password123", "tenant-123")

	if got := creds.GetUsername(); got != "user@example.com" {
		t.Errorf("GetUsername() = %v, want %v", got, "user@example.com")
	}

	if got := creds.GetPassword(); got != "password123" {
		t.Errorf("GetPassword() = %v, want %v", got, "password123")
	}

	if got := creds.GetTenantID(); got != "tenant-123" {
		t.Errorf("GetTenantID() = %v, want %v", got, "tenant-123")
	}
}

func TestNewEnv(t *testing.T) {
	os.Setenv("NHNCLOUD_ACCESS_KEY_ID", "env-access-key")
	os.Setenv("NHNCLOUD_SECRET_ACCESS_KEY", "env-secret-key")
	defer func() {
		os.Unsetenv("NHNCLOUD_ACCESS_KEY_ID")
		os.Unsetenv("NHNCLOUD_SECRET_ACCESS_KEY")
	}()

	creds := NewEnv()

	if got := creds.GetAccessKeyID(); got != "env-access-key" {
		t.Errorf("GetAccessKeyID() = %v, want %v", got, "env-access-key")
	}

	if got := creds.GetSecretAccessKey(); got != "env-secret-key" {
		t.Errorf("GetSecretAccessKey() = %v, want %v", got, "env-secret-key")
	}
}

func TestNewEnvIdentity(t *testing.T) {
	os.Setenv("NHNCLOUD_USERNAME", "env-user")
	os.Setenv("NHNCLOUD_PASSWORD", "env-password")
	os.Setenv("NHNCLOUD_TENANT_ID", "env-tenant")
	defer func() {
		os.Unsetenv("NHNCLOUD_USERNAME")
		os.Unsetenv("NHNCLOUD_PASSWORD")
		os.Unsetenv("NHNCLOUD_TENANT_ID")
	}()

	creds := NewEnvIdentity()

	if got := creds.GetUsername(); got != "env-user" {
		t.Errorf("GetUsername() = %v, want %v", got, "env-user")
	}

	if got := creds.GetPassword(); got != "env-password" {
		t.Errorf("GetPassword() = %v, want %v", got, "env-password")
	}

	if got := creds.GetTenantID(); got != "env-tenant" {
		t.Errorf("GetTenantID() = %v, want %v", got, "env-tenant")
	}
}

func TestStaticImplementsCredentials(t *testing.T) {
	var _ Credentials = (*Static)(nil)
}

func TestStaticIdentityImplementsIdentityCredentials(t *testing.T) {
	var _ IdentityCredentials = (*StaticIdentity)(nil)
}

func TestEnvImplementsCredentials(t *testing.T) {
	var _ Credentials = (*Env)(nil)
}

func TestEnvIdentityImplementsIdentityCredentials(t *testing.T) {
	var _ IdentityCredentials = (*EnvIdentity)(nil)
}
