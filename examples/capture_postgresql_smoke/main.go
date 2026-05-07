// Example: rds-postgresql smoke capture
//
// Exercises read-only PostgreSQL SDK endpoints (List/Get only — NO billable
// operations) so the capture middleware (NHN_SDK_CAPTURE_DIR, see
// nhncloud/internal/capture) mirrors response bodies to disk.
//
// The captured JSON files are consumed by Phase A.2.1 of the multi-engine
// clone plan to resolve UUIDs (flavor / subnet / parameter group / db
// version) used by the rds-postgresql-p99 scenario.
//
// Credentials are loaded from environment variables first; if any of the
// three required values are missing, we fall back to ~/.nhncloud/credentials
// (INI-ish format used by the CLI). The PostgreSQL app key may live under
// either the env var NHN_CLOUD_POSTGRESQL_APPKEY or the credentials key
// `rds_postgresql_app_key` (legacy CLI convention — shared across RDS engines).
//
// Run:
//
//	NHN_SDK_CAPTURE_DIR=$PWD/tests/captured-responses/rds-postgresql \
//	    go run ./examples/capture_postgresql_smoke
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
)

// loadCredentialsFile parses ~/.nhncloud/credentials (INI-style: `key = value`,
// with `#` comments and `[section]` headers ignored). Returns a flat map of
// the [default] section.
func loadCredentialsFile() map[string]string {
	out := map[string]string{}
	home, err := os.UserHomeDir()
	if err != nil {
		return out
	}
	path := filepath.Join(home, ".nhncloud", "credentials")
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		k := strings.TrimSpace(line[:eq])
		v := strings.TrimSpace(line[eq+1:])
		out[k] = v
	}
	return out
}

// resolveCreds picks values from env first, then the credentials file.
func resolveCreds() (region, appKey, accessKey, secretKey string) {
	file := loadCredentialsFile()

	region = os.Getenv("NHN_CLOUD_REGION")
	if region == "" {
		region = file["region"]
	}
	if region == "" {
		region = "kr1"
	}
	region = strings.ToLower(region)

	appKey = os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY")
	if appKey == "" {
		appKey = os.Getenv("NHN_POSTGRESQL_APP_KEY")
	}
	if appKey == "" {
		// CLI convention: rds_postgresql_app_key is PostgreSQL-specific; falls
		// back to postgresql_app_key, then to rds_app_key (which is shared
		// across RDS engines historically).
		if v := file["rds_postgresql_app_key"]; v != "" {
			appKey = v
		} else if v := file["postgresql_app_key"]; v != "" {
			appKey = v
		} else if v := file["rds_app_key"]; v != "" {
			appKey = v
		}
	}

	accessKey = os.Getenv("NHN_CLOUD_ACCESS_KEY_ID")
	if accessKey == "" {
		accessKey = os.Getenv("NHN_CLOUD_ACCESS_KEY")
	}
	if accessKey == "" {
		accessKey = file["access_key_id"]
	}

	secretKey = os.Getenv("NHN_CLOUD_SECRET_ACCESS_KEY")
	if secretKey == "" {
		secretKey = os.Getenv("NHN_CLOUD_SECRET_KEY")
	}
	if secretKey == "" {
		secretKey = file["secret_access_key"]
	}
	return
}

// callResult prints a one-line summary in the format the runner expects.
// Pattern: "[ok] Name -> N items" or "[err] Name: <message>".
func callResult(name string, count int, err error) {
	if err != nil {
		fmt.Printf("[err] %s: %v\n", name, err)
		return
	}
	fmt.Printf("[ok]  %s -> %d items\n", name, count)
}

func main() {
	region, appKey, accessKey, secretKey := resolveCreds()

	missing := []string{}
	if region == "" {
		missing = append(missing, "region")
	}
	if appKey == "" {
		missing = append(missing, "postgresql_app_key (NHN_CLOUD_POSTGRESQL_APPKEY or rds_postgresql_app_key)")
	}
	if accessKey == "" {
		missing = append(missing, "access_key_id")
	}
	if secretKey == "" {
		missing = append(missing, "secret_access_key")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "missing credentials: %s\n", strings.Join(missing, ", "))
		fmt.Fprintln(os.Stderr, "Set env vars or populate ~/.nhncloud/credentials")
		os.Exit(2)
	}

	captureDir := os.Getenv("NHN_SDK_CAPTURE_DIR")
	if captureDir == "" {
		fmt.Fprintln(os.Stderr, "warn: NHN_SDK_CAPTURE_DIR is not set; HTTP responses will NOT be mirrored to disk")
	} else {
		fmt.Fprintf(os.Stderr, "capture dir: %s\n", captureDir)
	}

	cfg := postgresql.Config{
		Region:    region,
		AppKey:    appKey,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	client, err := postgresql.NewClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client init failed: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// 1) Catalog/reference endpoints (always safe — no per-instance state).
	flavorsResp, ferr := client.ListFlavors(ctx)
	if ferr == nil && flavorsResp != nil {
		callResult("ListFlavors", len(flavorsResp.DBFlavors), nil)
	} else {
		callResult("ListFlavors", 0, ferr)
	}

	versionsResp, verr := client.ListVersions(ctx)
	if verr == nil && versionsResp != nil {
		callResult("ListVersions", len(versionsResp.DBVersions), nil)
	} else {
		callResult("ListVersions", 0, verr)
	}

	storageResp, serr := client.ListStorageTypes(ctx)
	if serr == nil && storageResp != nil {
		callResult("ListStorageTypes", 1, nil)
		_ = storageResp
	} else {
		callResult("ListStorageTypes", 0, serr)
	}

	subnetsResp, snErr := client.ListSubnets(ctx)
	if snErr == nil && subnetsResp != nil {
		callResult("ListSubnets", len(subnetsResp.Subnets), nil)
	} else {
		callResult("ListSubnets", 0, snErr)
	}

	// 2) Account-scoped list endpoints.
	pgResp, pgErr := client.ListParameterGroups(ctx)
	if pgErr == nil && pgResp != nil {
		callResult("ListParameterGroups", 1, nil)
		_ = pgResp
	} else {
		callResult("ListParameterGroups", 0, pgErr)
	}

	ngResp, ngErr := client.ListNotificationGroups(ctx)
	if ngErr == nil && ngResp != nil {
		callResult("ListNotificationGroups", 1, nil)
		_ = ngResp
	} else {
		callResult("ListNotificationGroups", 0, ngErr)
	}

	ugResp, ugErr := client.ListUserGroups(ctx)
	if ugErr == nil && ugResp != nil {
		callResult("ListUserGroups", 1, nil)
		_ = ugResp
	} else {
		callResult("ListUserGroups", 0, ugErr)
	}

	sgResp, sgErr := client.ListSecurityGroups(ctx)
	if sgErr == nil && sgResp != nil {
		callResult("ListSecurityGroups", 1, nil)
		_ = sgResp
	} else {
		callResult("ListSecurityGroups", 0, sgErr)
	}

	metricsResp, mErr := client.ListMetrics(ctx)
	if mErr == nil && metricsResp != nil {
		callResult("ListMetrics", 1, nil)
		_ = metricsResp
	} else {
		callResult("ListMetrics", 0, mErr)
	}

	// 3) Instances list — and if any exist, drill into the first one for
	// per-instance reads (still all GETs).
	instancesResp, iErr := client.ListInstances(ctx)
	var firstInstanceID string
	if iErr == nil && instancesResp != nil {
		callResult("ListInstances", len(instancesResp.DBInstances), nil)
		if len(instancesResp.DBInstances) > 0 {
			firstInstanceID = instancesResp.DBInstances[0].DBInstanceID
		}
	} else {
		callResult("ListInstances", 0, iErr)
	}

	if firstInstanceID != "" {
		fmt.Fprintf(os.Stderr, "drilling into first instance: %s\n", firstInstanceID)

		if _, err := client.GetInstance(ctx, firstInstanceID); err != nil {
			callResult("GetInstance", 0, err)
		} else {
			callResult("GetInstance", 1, nil)
		}

		if br, err := client.ListBackups(ctx, firstInstanceID); err != nil {
			callResult("ListBackups", 0, err)
		} else {
			_ = br
			callResult("ListBackups", 1, nil)
		}

		if _, err := client.GetNetworkInfo(ctx, firstInstanceID); err != nil {
			callResult("GetNetworkInfo", 0, err)
		} else {
			callResult("GetNetworkInfo", 1, nil)
		}

		if _, err := client.GetStorageInfo(ctx, firstInstanceID); err != nil {
			callResult("GetStorageInfo", 0, err)
		} else {
			callResult("GetStorageInfo", 1, nil)
		}

		// PostgreSQL-specific: list databases, users, extensions, hba rules.
		if _, err := client.ListDatabases(ctx, firstInstanceID); err != nil {
			callResult("ListDatabases", 0, err)
		} else {
			callResult("ListDatabases", 1, nil)
		}

		if _, err := client.ListDBUsers(ctx, firstInstanceID); err != nil {
			callResult("ListDBUsers", 0, err)
		} else {
			callResult("ListDBUsers", 1, nil)
		}

		if _, err := client.ListExtensions(ctx, firstInstanceID); err != nil {
			callResult("ListExtensions", 0, err)
		} else {
			callResult("ListExtensions", 1, nil)
		}

		if _, err := client.ListHBARules(ctx, firstInstanceID); err != nil {
			callResult("ListHBARules", 0, err)
		} else {
			callResult("ListHBARules", 1, nil)
		}
	} else {
		fmt.Fprintln(os.Stderr, "note: no instances in account — per-instance endpoints skipped")
	}

	fmt.Fprintln(os.Stderr, "done.")
}
