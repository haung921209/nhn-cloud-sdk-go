package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

type TestResult struct {
	Name    string
	Passed  bool
	Message string
	Latency time.Duration
}

var results []TestResult

func logTest(name string, passed bool, message string, latency time.Duration) {
	results = append(results, TestResult{name, passed, message, latency})
	if passed {
		fmt.Printf("✅ PASS: %s (%v)\n", name, latency)
	} else {
		fmt.Printf("❌ FAIL: %s - %s\n", name, message)
	}
}

func main() {
	accessKeyID := os.Getenv("NHN_CLOUD_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("NHN_CLOUD_SECRET_ACCESS_KEY")
	mysqlAppKey := os.Getenv("NHN_CLOUD_MYSQL_APPKEY")
	mariadbAppKey := os.Getenv("NHN_CLOUD_MARIADB_APPKEY")
	pgAppKey := os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY")

	if accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("NHN_CLOUD_ACCESS_KEY_ID and NHN_CLOUD_SECRET_ACCESS_KEY required")
	}

	fmt.Println("========================================")
	fmt.Println("NHN Cloud SDK - Integration Test")
	fmt.Println("========================================")
	fmt.Printf("Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	creds := credentials.NewStatic(accessKeyID, secretAccessKey)
	cfg := &nhncloud.Config{
		Region:      "kr1",
		Credentials: creds,
		AppKeys: map[string]string{
			"rds-mysql":      mysqlAppKey,
			"rds-mariadb":    mariadbAppKey,
			"rds-postgresql": pgAppKey,
		},
	}

	client, err := nhncloud.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("--- MySQL Tests ---")
	if mysqlAppKey != "" {
		testMySQLListInstances(ctx, client)
		testMySQLListFlavors(ctx, client)
		testMySQLListVersions(ctx, client)
		testMySQLListParameterGroups(ctx, client)
		testMySQLListSecurityGroups(ctx, client)
		testMySQLListBackups(ctx, client)
	} else {
		fmt.Println("⏭️  Skipped (no MySQL AppKey)")
	}

	fmt.Println("\n--- MariaDB Tests ---")
	if mariadbAppKey != "" {
		testMariaDBListInstances(ctx, client)
		testMariaDBListFlavors(ctx, client)
	} else {
		fmt.Println("⏭️  Skipped (no MariaDB AppKey)")
	}

	fmt.Println("\n--- PostgreSQL Tests ---")
	if pgAppKey != "" {
		testPostgreSQLListInstances(ctx, client)
		testPostgreSQLListFlavors(ctx, client)
	} else {
		fmt.Println("⏭️  Skipped (no PostgreSQL AppKey)")
	}

	printSummary()
}

func testMySQLListInstances(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	instances, err := client.MySQL().ListInstances(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListInstances", false, err.Error(), latency)
		return
	}
	logTest("MySQL ListInstances", true, fmt.Sprintf("%d instances found", len(instances.DBInstances)), latency)
}

func testMySQLListFlavors(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	flavors, err := client.MySQL().ListFlavors(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListFlavors", false, err.Error(), latency)
		return
	}
	if len(flavors.DBFlavors) == 0 {
		logTest("MySQL ListFlavors", false, "no flavors returned", latency)
		return
	}
	logTest("MySQL ListFlavors", true, fmt.Sprintf("%d flavors found", len(flavors.DBFlavors)), latency)
}

func testMySQLListVersions(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	versions, err := client.MySQL().ListVersions(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListVersions", false, err.Error(), latency)
		return
	}
	if len(versions.DBVersions) == 0 {
		logTest("MySQL ListVersions", false, "no versions returned", latency)
		return
	}
	logTest("MySQL ListVersions", true, fmt.Sprintf("%d versions found", len(versions.DBVersions)), latency)
}

func testMySQLListParameterGroups(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	groups, err := client.MySQL().ListParameterGroups(ctx, "")
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListParameterGroups", false, err.Error(), latency)
		return
	}
	logTest("MySQL ListParameterGroups", true, fmt.Sprintf("%d groups found", len(groups.ParameterGroups)), latency)
}

func testMySQLListSecurityGroups(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	groups, err := client.MySQL().ListSecurityGroups(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListSecurityGroups", false, err.Error(), latency)
		return
	}
	logTest("MySQL ListSecurityGroups", true, fmt.Sprintf("%d groups found", len(groups.DBSecurityGroups)), latency)
}

func testMySQLListBackups(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	backups, err := client.MySQL().ListBackups(ctx, "", "", 0, 10)
	latency := time.Since(start)

	if err != nil {
		logTest("MySQL ListBackups", false, err.Error(), latency)
		return
	}
	logTest("MySQL ListBackups", true, fmt.Sprintf("%d backups found", len(backups.Backups)), latency)
}

func testMariaDBListInstances(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	instances, err := client.MariaDB().ListInstances(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MariaDB ListInstances", false, err.Error(), latency)
		return
	}
	logTest("MariaDB ListInstances", true, fmt.Sprintf("%d instances found", len(instances.DBInstances)), latency)
}

func testMariaDBListFlavors(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	flavors, err := client.MariaDB().ListFlavors(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("MariaDB ListFlavors", false, err.Error(), latency)
		return
	}
	logTest("MariaDB ListFlavors", true, fmt.Sprintf("%d flavors found", len(flavors.DBFlavors)), latency)
}

func testPostgreSQLListInstances(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	instances, err := client.PostgreSQL().ListInstances(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("PostgreSQL ListInstances", false, err.Error(), latency)
		return
	}
	logTest("PostgreSQL ListInstances", true, fmt.Sprintf("%d instances found", len(instances.DBInstances)), latency)
}

func testPostgreSQLListFlavors(ctx context.Context, client *nhncloud.Client) {
	start := time.Now()
	flavors, err := client.PostgreSQL().ListFlavors(ctx)
	latency := time.Since(start)

	if err != nil {
		logTest("PostgreSQL ListFlavors", false, err.Error(), latency)
		return
	}
	logTest("PostgreSQL ListFlavors", true, fmt.Sprintf("%d flavors found", len(flavors.DBFlavors)), latency)
}

func printSummary() {
	fmt.Println("\n========================================")
	fmt.Println("Test Summary")
	fmt.Println("========================================")

	passed := 0
	failed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total:  %d\n", len(results))

	if failed > 0 {
		fmt.Println("\n❌ Some tests failed!")
		os.Exit(1)
	} else {
		fmt.Println("\n✅ All tests passed!")
	}
}
