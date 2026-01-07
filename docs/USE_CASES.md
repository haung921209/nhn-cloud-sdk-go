# NHN Cloud SDK Use Cases

> **Verified**: 2026-01-08 | **SDK Version**: v0.1.x | **Test Status**: 10/10 Passed

This document contains verified use cases that have been tested against live NHN Cloud API.

## Table of Contents

- [Quick Setup](#quick-setup)
- [RDS MySQL Operations](#rds-mysql-operations)
- [RDS MariaDB Operations](#rds-mariadb-operations)
- [RDS PostgreSQL Operations](#rds-postgresql-operations)
- [Production Use Cases](#production-use-cases)

---

## Quick Setup

### Basic Client Configuration

```go
package main

import (
    "context"
    "log"
    
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

func main() {
    // Create OAuth credentials
    creds := credentials.NewStatic(
        "your-access-key-id",
        "your-secret-access-key",
    )

    // Configure client with service-specific app keys
    cfg := &nhncloud.Config{
        Region:      "kr1",
        Credentials: creds,
        AppKeys: map[string]string{
            "rds-mysql":      "your-mysql-appkey",
            "rds-mariadb":    "your-mariadb-appkey",
            "rds-postgresql": "your-postgresql-appkey",
        },
    }

    client, err := nhncloud.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    // Use client...
}
```

### Environment Variables Setup

```bash
export NHN_CLOUD_ACCESS_KEY_ID="your-access-key-id"
export NHN_CLOUD_SECRET_ACCESS_KEY="your-secret-access-key"
export NHN_CLOUD_MYSQL_APPKEY="your-mysql-appkey"
export NHN_CLOUD_MARIADB_APPKEY="your-mariadb-appkey"
export NHN_CLOUD_POSTGRESQL_APPKEY="your-postgresql-appkey"
```

---

## RDS MySQL Operations

### List MySQL Instances (Verified)

```go
// List all MySQL instances
instances, err := client.MySQL().ListInstances(ctx)
if err != nil {
    log.Fatal(err)
}

for _, inst := range instances.DBInstances {
    fmt.Printf("Instance: %s\n", inst.DBInstanceName)
    fmt.Printf("  ID: %s\n", inst.DBInstanceID)
    fmt.Printf("  Status: %s\n", inst.DBInstanceStatus)
    fmt.Printf("  Type: %s\n", inst.DBInstanceType)
}
```

**Example Output**:
```
Instance: cli-test-mysql-replica
  ID: ef95583a-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Status: AVAILABLE
  Type: NORMAL

Instance: test-replica-1959
  ID: 7fa4afa7-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Status: AVAILABLE
  Type: NORMAL
```

### List MySQL Flavors (Verified)

```go
// List all available MySQL flavors (instance types)
flavors, err := client.MySQL().ListFlavors(ctx)
if err != nil {
    log.Fatal(err)
}

for _, f := range flavors.DBFlavors {
    fmt.Printf("Flavor: %s\n", f.DBFlavorName)
    fmt.Printf("  VCPUs: %d, RAM: %dMB\n", f.VCPUS, f.RAM)
}
```

**Available Flavors** (as of 2026-01-08):
| Flavor Name | VCPUs | RAM (MB) |
|-------------|-------|----------|
| m2.c1m2 | 1 | 2048 |
| m2.c1m4 | 1 | 4096 |
| m2.c2m4 | 2 | 4096 |
| m2.c2m8 | 2 | 8192 |
| m2.c4m8 | 4 | 8192 |
| m2.c4m16 | 4 | 16384 |
| m2.c8m16 | 8 | 16384 |
| m2.c8m32 | 8 | 32768 |
| ... | ... | ... |

### List MySQL Versions (Verified)

```go
// List all available MySQL versions
versions, err := client.MySQL().ListVersions(ctx)
if err != nil {
    log.Fatal(err)
}

for _, v := range versions.DBVersions {
    fmt.Printf("Version: %s (%s)\n", v.DBVersion, v.DBVersionID)
}
```

**Available Versions**:
- MySQL 8.0.33
- MySQL 8.0.32
- MySQL 8.0.28
- MySQL 5.7.x (various)

### List Parameter Groups (Verified)

```go
// List MySQL parameter groups
groups, err := client.MySQL().ListParameterGroups(ctx, "")
if err != nil {
    log.Fatal(err)
}

for _, g := range groups.ParameterGroups {
    fmt.Printf("Parameter Group: %s\n", g.ParameterGroupName)
    fmt.Printf("  ID: %s\n", g.ParameterGroupID)
    fmt.Printf("  Engine: %s %s\n", g.DBEngine, g.DBEngineVersion)
}
```

### List Security Groups (Verified)

```go
// List DB security groups
secGroups, err := client.MySQL().ListSecurityGroups(ctx)
if err != nil {
    log.Fatal(err)
}

for _, sg := range secGroups.DBSecurityGroups {
    fmt.Printf("Security Group: %s\n", sg.DBSecurityGroupName)
    fmt.Printf("  ID: %s\n", sg.DBSecurityGroupID)
}
```

### List Backups (Verified)

```go
// List MySQL backups (paginated)
backups, err := client.MySQL().ListBackups(ctx, "", "", 0, 10)
if err != nil {
    log.Fatal(err)
}

for _, b := range backups.Backups {
    fmt.Printf("Backup: %s\n", b.BackupID)
    fmt.Printf("  Instance: %s\n", b.DBInstanceID)
    fmt.Printf("  Status: %s\n", b.BackupStatus)
    fmt.Printf("  Size: %dGB\n", b.BackupSize)
}
```

---

## RDS MariaDB Operations

### List MariaDB Instances (Verified)

```go
// List all MariaDB instances
instances, err := client.MariaDB().ListInstances(ctx)
if err != nil {
    log.Fatal(err)
}

for _, inst := range instances.DBInstances {
    fmt.Printf("Instance: %s (Status: %s)\n", 
        inst.DBInstanceName, inst.DBInstanceStatus)
}
```

### List MariaDB Flavors (Verified)

```go
// List all available MariaDB flavors
flavors, err := client.MariaDB().ListFlavors(ctx)
if err != nil {
    log.Fatal(err)
}

for _, f := range flavors.DBFlavors {
    fmt.Printf("Flavor: %s (VCPUs: %d, RAM: %dMB)\n", 
        f.DBFlavorName, f.VCPUS, f.RAM)
}
```

---

## RDS PostgreSQL Operations

### List PostgreSQL Instances (Verified)

```go
// List all PostgreSQL instances
instances, err := client.PostgreSQL().ListInstances(ctx)
if err != nil {
    log.Fatal(err)
}

for _, inst := range instances.DBInstances {
    fmt.Printf("Instance: %s (Status: %s)\n", 
        inst.DBInstanceName, inst.DBInstanceStatus)
}
```

### List PostgreSQL Flavors (Verified)

```go
// List all available PostgreSQL flavors
flavors, err := client.PostgreSQL().ListFlavors(ctx)
if err != nil {
    log.Fatal(err)
}

for _, f := range flavors.DBFlavors {
    fmt.Printf("Flavor: %s (VCPUs: %d, RAM: %dMB)\n", 
        f.DBFlavorName, f.VCPUS, f.RAM)
}
```

---

## Production Use Cases

### Use Case 1: Database Inventory Management

Collect inventory of all database instances across services:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

type DBInventory struct {
    MySQL      []InstanceInfo `json:"mysql"`
    MariaDB    []InstanceInfo `json:"mariadb"`
    PostgreSQL []InstanceInfo `json:"postgresql"`
}

type InstanceInfo struct {
    Name   string `json:"name"`
    ID     string `json:"id"`
    Status string `json:"status"`
    Type   string `json:"type"`
}

func main() {
    creds := credentials.NewStatic(
        os.Getenv("NHN_CLOUD_ACCESS_KEY_ID"),
        os.Getenv("NHN_CLOUD_SECRET_ACCESS_KEY"),
    )

    cfg := &nhncloud.Config{
        Region:      "kr1",
        Credentials: creds,
        AppKeys: map[string]string{
            "rds-mysql":      os.Getenv("NHN_CLOUD_MYSQL_APPKEY"),
            "rds-mariadb":    os.Getenv("NHN_CLOUD_MARIADB_APPKEY"),
            "rds-postgresql": os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY"),
        },
    }

    client, _ := nhncloud.New(cfg)
    ctx := context.Background()

    inventory := DBInventory{}

    // MySQL
    if mysql, err := client.MySQL().ListInstances(ctx); err == nil {
        for _, inst := range mysql.DBInstances {
            inventory.MySQL = append(inventory.MySQL, InstanceInfo{
                Name:   inst.DBInstanceName,
                ID:     inst.DBInstanceID,
                Status: inst.DBInstanceStatus,
                Type:   inst.DBInstanceType,
            })
        }
    }

    // MariaDB
    if mariadb, err := client.MariaDB().ListInstances(ctx); err == nil {
        for _, inst := range mariadb.DBInstances {
            inventory.MariaDB = append(inventory.MariaDB, InstanceInfo{
                Name:   inst.DBInstanceName,
                ID:     inst.DBInstanceID,
                Status: inst.DBInstanceStatus,
                Type:   inst.DBInstanceType,
            })
        }
    }

    // PostgreSQL
    if pg, err := client.PostgreSQL().ListInstances(ctx); err == nil {
        for _, inst := range pg.DBInstances {
            inventory.PostgreSQL = append(inventory.PostgreSQL, InstanceInfo{
                Name:   inst.DBInstanceName,
                ID:     inst.DBInstanceID,
                Status: inst.DBInstanceStatus,
                Type:   inst.DBInstanceType,
            })
        }
    }

    // Output as JSON
    output, _ := json.MarshalIndent(inventory, "", "  ")
    fmt.Println(string(output))
}
```

### Use Case 2: Capacity Planning - Flavor Selection

Help select appropriate flavors based on workload requirements:

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

func main() {
    // ... client setup ...

    ctx := context.Background()
    
    // Get MySQL flavors
    flavors, _ := client.MySQL().ListFlavors(ctx)
    
    // Filter by requirements
    minRAM := 8192    // 8GB minimum
    minVCPUs := 4     // 4 vCPUs minimum
    
    fmt.Println("Recommended MySQL Flavors:")
    fmt.Println("==========================")
    
    for _, f := range flavors.DBFlavors {
        if f.RAM >= minRAM && f.VCPUS >= minVCPUs {
            fmt.Printf("  %s: %d vCPUs, %d MB RAM\n", 
                f.DBFlavorName, f.VCPUS, f.RAM)
        }
    }
}
```

### Use Case 3: Backup Status Monitoring

Monitor backup status across all MySQL instances:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
)

func main() {
    // ... client setup ...

    ctx := context.Background()
    
    // Get recent backups (last 50)
    backups, err := client.MySQL().ListBackups(ctx, "", "", 0, 50)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Backup Status Report - %s\n", time.Now().Format("2006-01-02"))
    fmt.Println("=========================================")
    
    statusCount := make(map[string]int)
    for _, b := range backups.Backups {
        statusCount[b.BackupStatus]++
    }
    
    fmt.Printf("Total Backups: %d\n", len(backups.Backups))
    for status, count := range statusCount {
        fmt.Printf("  %s: %d\n", status, count)
    }
}
```

---

## Test Results Summary

| Service | Operation | Status | Latency |
|---------|-----------|--------|---------|
| MySQL | ListInstances | PASS | ~141ms |
| MySQL | ListFlavors | PASS | ~553ms |
| MySQL | ListVersions | PASS | ~48ms |
| MySQL | ListParameterGroups | PASS | ~56ms |
| MySQL | ListSecurityGroups | PASS | ~41ms |
| MySQL | ListBackups | PASS | ~36ms |
| MariaDB | ListInstances | PASS | ~71ms |
| MariaDB | ListFlavors | PASS | ~584ms |
| PostgreSQL | ListInstances | PASS | ~191ms |
| PostgreSQL | ListFlavors | PASS | ~560ms |

**Total: 10/10 Tests Passed**

---

## Limitations & Notes

### API Rate Limits
- Flavor listing APIs have higher latency (~500ms) compared to instance listing (~100ms)
- Recommend caching flavor lists as they change infrequently

### Authentication
- OAuth tokens are automatically managed by the SDK
- AppKeys are service-specific (separate keys for MySQL, MariaDB, PostgreSQL)

### Not Yet Tested
The following operations are available in the SDK but not yet verified:
- CreateInstance
- DeleteInstance
- ModifyInstance
- EnableHighAvailability
- CreateReadReplica
- Failover operations

---

## Related Documentation

- [SDK README](/README.md) - Installation and quick start
- [Examples](/examples/) - Runnable example code
- [NHN Cloud RDS API Documentation](https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide/)
