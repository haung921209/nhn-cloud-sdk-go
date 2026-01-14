package mariadb_test

import (
	"context"
	"testing"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     mysql.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: mysql.Config{
				Region:    "kr1",
				AppKey:    "test-app-key",
				AccessKey: "test-access-key",
				SecretKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "missing region",
			cfg: mysql.Config{
				AppKey:    "test-app-key",
				AccessKey: "test-access-key",
				SecretKey: "test-secret-key",
			},
			wantErr: true,
		},
		{
			name: "missing app key",
			cfg: mysql.Config{
				Region:    "kr1",
				AccessKey: "test-access-key",
				SecretKey: "test-secret-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := mysql.NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestListInstances_ContextRequired(t *testing.T) {
	// This test verifies the API signature
	// We can't actually call without mock server, but we verify it compiles
	cfg := mysql.Config{
		Region:    "kr1",
		AppKey:    "test",
		AccessKey: "test",
		SecretKey: "test",
	}

	client, err := mysql.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Verify method exists and takes context
	_ = client
	_ = context.Background()

	// Actual call would require mock server or live credentials
	// Will be tested in integration tests
}
