# NHN Cloud SDK for Go

Go SDK for NHN Cloud services.

## Installation

```bash
go get github.com/haung921209/nhn-cloud-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

func main() {
	// Create credentials
	creds := credentials.NewStatic(
		"your-access-key-id",     // NHN Cloud Access Key ID
		"your-secret-access-key", // NHN Cloud Secret Access Key
	)

	// Identity credentials for Compute/Network services
	identityCreds := credentials.NewStaticIdentity(
		"your-username",  // NHN Cloud Username (email)
		"your-password",  // NHN Cloud API Password
		"your-tenant-id", // Project Tenant ID
	)

	// Create client
	cfg := &nhncloud.Config{
		Region:              "kr1",
		Credentials:         creds,
		IdentityCredentials: identityCreds,
		AppKeys: map[string]string{
			"rds-mysql": "your-rds-mysql-appkey",
		},
	}

	client, err := nhncloud.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List IAM organizations
	orgs, err := client.IAM().ListOrganizations(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, org := range orgs.Organizations {
		fmt.Printf("Organization: %s (%s)\n", org.Name, org.ID)
	}

	// List Compute instances
	servers, err := client.Compute().ListServers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, server := range servers.Servers {
		fmt.Printf("Server: %s (%s) - %s\n", server.Name, server.ID, server.Status)
	}

	// List RDS MySQL instances
	instances, err := client.MySQL().ListInstances(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, inst := range instances.Instances {
		fmt.Printf("MySQL: %s (%s) - %s\n", inst.Name, inst.ID, inst.Status)
	}
}
```

## Authentication

### OAuth Credentials (IAM, RDS, Object Storage)

Used for services that require OAuth 2.0 authentication.

```go
creds := credentials.NewStatic(
	"access-key-id",
	"secret-access-key",
)

// Or from environment variables:
// NHN_CLOUD_ACCESS_KEY_ID
// NHN_CLOUD_SECRET_ACCESS_KEY
creds := credentials.NewEnvCredentials()
```

### Identity Credentials (Compute, Network, Block Storage)

Used for OpenStack-based services.

```go
identityCreds := credentials.NewStaticIdentity(
	"username",   // Email address
	"password",   // API password (set in NHN Cloud Console)
	"tenant-id",  // Project tenant ID
)

// Or from environment variables:
// NHN_CLOUD_USERNAME
// NHN_CLOUD_PASSWORD
// NHN_CLOUD_TENANT_ID
identityCreds := credentials.NewEnvIdentityCredentials()
```

## Supported Services

| Service | Client Method | Auth Type |
|---------|--------------|-----------|
| IAM | `client.IAM()` | OAuth |
| Compute | `client.Compute()` | Identity |
| VPC | `client.VPC()` | Identity |
| Security Group | `client.SecurityGroup()` | Identity |
| Floating IP | `client.FloatingIP()` | Identity |
| Load Balancer | `client.LoadBalancer()` | Identity |
| Block Storage | `client.BlockStorage()` | Identity |
| Object Storage | `client.ObjectStorage()` | Identity |
| RDS MySQL | `client.MySQL()` | OAuth + AppKey |
| RDS MariaDB | `client.MariaDB()` | OAuth + AppKey |
| RDS PostgreSQL | `client.PostgreSQL()` | OAuth + AppKey |
| NKS | `client.NKS()` | Identity |
| NCR | `client.NCR()` | OAuth + AppKey |
| NCS | `client.NCS()` | OAuth + AppKey |

## Examples

### RDS MySQL

```go
ctx := context.Background()

// List instances
instances, _ := client.MySQL().ListInstances(ctx)

// Get instance details
instance, _ := client.MySQL().GetInstance(ctx, "instance-id")

// Create instance
input := &mysql.CreateInstanceInput{
	Name:     "my-database",
	FlavorID: "flavor-id",
	Version:  "MYSQL_V8033",
	UserName: "admin",
	Password: "SecurePassword123!",
	Network: &mysql.Network{
		SubnetID: "subnet-id",
	},
	Storage: &mysql.Storage{
		StorageType: "General SSD",
		StorageSize: 20,
	},
	Backup: &mysql.BackupConfig{
		Period: 1,
		Schedules: []mysql.BackupSchedule{
			{BeginTime: "00:00", Duration: "02:00"},
		},
	},
	ParameterGroupID: "parameter-group-id",
}
result, _ := client.MySQL().CreateInstance(ctx, input)

// Delete instance
client.MySQL().DeleteInstance(ctx, "instance-id")
```

### Compute

```go
ctx := context.Background()

// List servers
servers, _ := client.Compute().ListServers(ctx)

// Create server
input := &compute.CreateServerInput{
	Name:      "my-server",
	ImageRef:  "image-id",
	FlavorRef: "flavor-id",
	KeyName:   "my-keypair",
	Networks: []compute.ServerNetwork{
		{UUID: "network-id"},
	},
}
result, _ := client.Compute().CreateServer(ctx, input)

// Server actions
client.Compute().StopServer(ctx, "server-id")
client.Compute().StartServer(ctx, "server-id")
client.Compute().RebootServer(ctx, "server-id", false) // soft reboot
```

## Configuration Options

```go
cfg := &nhncloud.Config{
	// Required
	Region:      "kr1",              // kr1, kr2, jp1
	Credentials: creds,              // OAuth credentials

	// Optional
	IdentityCredentials: identityCreds, // For Compute/Network
	AppKeys: map[string]string{         // Service-specific app keys
		"rds-mysql":      "...",
		"rds-mariadb":    "...",
		"rds-postgresql": "...",
		"ncr":            "...",
		"ncs":            "...",
	},
	HTTPClient: customHTTPClient,  // Custom HTTP client
	Debug:      true,               // Enable debug logging
	UserAgent:  "my-app/1.0",      // Custom user agent
}
```

## Error Handling

```go
result, err := client.MySQL().GetInstance(ctx, "invalid-id")
if err != nil {
	if apiErr, ok := err.(*client.APIError); ok {
		fmt.Printf("API Error: %d - %s\n", apiErr.StatusCode, apiErr.Message)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
```

## License

Apache License 2.0
