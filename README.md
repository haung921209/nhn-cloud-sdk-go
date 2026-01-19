# NHN Cloud SDK for Go

This is an unofficial, community-driven Go SDK for NHN Cloud, designed to provide a clean, idiomatic interface for managing NHN Cloud resources.

> **Status**: ✅ **Refactored & Verified** (January 2026) -> "v2-rebuild"

## Supported Services

| Service | Package | Status | Notes |
|---------|---------|--------|-------|
| **Compute** | `nhncloud/compute` | 🟢 Verified | instances, flavors, images, keypairs |
| **Network** | `nhncloud/network` | 🟢 Verified | vpc, subnets, floating-ips |
| **NKS (Kubernetes)** | `nhncloud/container/nks` | 🟢 Verified | clusters, nodegroups, `describe-versions` |
| **NCR (Registry)** | `nhncloud/container/ncr` | 🟢 Verified | registries, images |
| **NCS (Container)** | `nhncloud/container/ncs` | 🟢 Logic Verified | workloads, services (Region Availability Issues) |
| **Object Storage** | `nhncloud/storage/objectstorage` | 🟢 Verified | containers, objects |
| **NAS** | `nhncloud/storage/nas` | 🟢 Verified | volumes, snapshots |
| **RDS MySQL** | `nhncloud/database/mysql` | 🟢 Verified | instances, backups |
| **RDS MariaDB** | `nhncloud/database/mariadb` | 🟢 Verified | instances, backups |
| **RDS PostgreSQL** | `nhncloud/database/postgresql` | 🟢 Verified | instances, backups |
| **Cloud Monitoring** | `nhncloud/monitoring` | 🟢 Verified | alarms (AppKey auth) |

## Installation

```bash
go get github.com/haung921209/nhn-cloud-sdk-go
```

## Configuration

The SDK unifies authentication for all services.

### 1. Environment Variables
Set these to avoid checking credentials into code:
- `NHN_CLOUD_REGION` (e.g. `kr1`)
- `NHN_CLOUD_APPKEY`
- `NHN_CLOUD_TENANT_ID`
- `NHN_CLOUD_USERNAME` / `NHN_CLOUD_PASSWORD`

### 2. TLS/SSL Setup (RDS)
To securely connect to RDS MySQL/MariaDB/Postgres, you must register the Root CA.

```go
// Download: https://static.toastoven.net/toastcloud/sdk_download/rds/ca-certificate.crt
rootCertPool := x509.NewCertPool()
pem, _ := os.ReadFile("ca-certificate.crt")
rootCertPool.AppendCertsFromPEM(pem)
// Use this config in your driver
mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
```

## Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
)

func main() {
	creds := credentials.NewUserCredentials("tenant-id", "username", "password")
	client := compute.NewClient("kr1", creds, http.DefaultClient, true)

	instances, err := client.ListInstances(context.Background())
	if err != nil {
		panic(err)
	}

	for _, instance := range instances.Servers {
		fmt.Printf("Instance: %s (%s)\n", instance.Name, instance.ID)
	}
}
```
