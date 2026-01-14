// Example: List MySQL instances
//
// This example demonstrates how to list all MySQL database instances
// using the new v2 SDK.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
)

func main() {
	// Create client configuration
	cfg := mysql.Config{
		Region:    os.Getenv("NHN_CLOUD_REGION"), // e.g., "kr1"
		AppKey:    os.Getenv("NHN_CLOUD_MYSQL_APPKEY"),
		AccessKey: os.Getenv("NHN_CLOUD_ACCESS_KEY"),
		SecretKey: os.Getenv("NHN_CLOUD_SECRET_KEY"),
	}

	// Initialize MySQL client
	client, err := mysql.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List all instances
	result, err := client.ListInstances(context.Background())
	if err != nil {
		log.Fatalf("Failed to list instances: %v", err)
	}

	// Display results
	fmt.Printf("Found %d MySQL instances:\n\n", len(result.DBInstances))

	for i, inst := range result.DBInstances {
		fmt.Printf("[%d] %s\n", i+1, inst.DBInstanceName)
		fmt.Printf("    ID: %s\n", inst.DBInstanceID)
		fmt.Printf("    Status: %s\n", inst.DBInstanceStatus)
		fmt.Printf("    Version: %s\n", inst.DBVersion)
		fmt.Printf("    Port: %d\n", inst.DBPort)
		fmt.Printf("    Storage: %s (%dGB)\n", inst.Storage.StorageType, inst.Storage.StorageSize)
		fmt.Printf("    Created: %s\n", inst.CreatedAt)
		fmt.Println()
	}
}
