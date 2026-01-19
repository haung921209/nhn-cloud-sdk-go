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
