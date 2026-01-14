// Package postgresql provides a client for NHN Cloud RDS for PostgreSQL v1.0 API.
//
// Official API Documentation:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/
//
// Key Differences from PostgreSQL/MariaDB:
//   - API Version: v1.0 (not v3.0)
//   - Authentication: Bearer Token (not OAuth2 headers)
//   - Database Management: /databases (not /db-schemas)
//   - PostgreSQL-specific: Extensions, HBA Rules
//   - Port: 5432 (not 3306)
//   - Instance creation requires databaseName field
//
// Example:
//
//	cfg := postgresql.Config{
//	    Region: "kr1",
//	    AppKey: "your-app-key",
//	    Token:  "your-bearer-token", // Obtained via OAuth2
//	}
//
//	client, err := postgresql.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	instances, err := client.ListInstances(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, inst := range instances.DBInstances {
//	    fmt.Printf("%s: %s\n", inst.DBInstanceName, inst.DBInstanceStatus)
//	}
package postgresql

import (
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Client is the PostgreSQL API client
type Client struct {
	core *core.Client
}

// Config holds PostgreSQL client configuration
type Config struct {
	Region string
	AppKey string
	Token  string // Bearer token (obtained via OAuth2 token exchange)
}

// NewClient creates a new PostgreSQL client.
//
// Note: PostgreSQL v1.0 uses Bearer Token authentication, different from
// PostgreSQL/MariaDB v3.0 which uses OAuth2 headers.
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// PostgreSQL uses different base URL: rds-postgres (not rds-postgresql)
	baseURL := fmt.Sprintf("%s-rds-postgres.api.nhncloudservice.com", cfg.Region)

	// PostgreSQL uses Bearer token authentication
	authenticator := auth.NewBearerAuth(cfg.AppKey, cfg.Token)

	coreClient := core.NewClient(baseURL, authenticator, nil)

	return &Client{
		core: coreClient,
	}, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Region == "" {
		return &core.ValidationError{Field: "Region", Message: "region is required"}
	}
	if c.AppKey == "" {
		return &core.ValidationError{Field: "AppKey", Message: "app key is required"}
	}
	if c.Token == "" {
		return &core.ValidationError{Field: "Token", Message: "bearer token is required"}
	}
	return nil
}
