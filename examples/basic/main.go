package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
)

func main() {
	accessKeyID := os.Getenv("NHN_CLOUD_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("NHN_CLOUD_SECRET_ACCESS_KEY")
	username := os.Getenv("NHN_CLOUD_USERNAME")
	password := os.Getenv("NHN_CLOUD_PASSWORD")
	tenantID := os.Getenv("NHN_CLOUD_TENANT_ID")
	mysqlAppKey := os.Getenv("NHN_CLOUD_MYSQL_APPKEY")

	if accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("NHN_CLOUD_ACCESS_KEY_ID and NHN_CLOUD_SECRET_ACCESS_KEY are required")
	}

	creds := credentials.NewStatic(accessKeyID, secretAccessKey)

	var identityCreds credentials.IdentityCredentials
	if username != "" && password != "" && tenantID != "" {
		identityCreds = credentials.NewStaticIdentity(username, password, tenantID)
	}

	cfg := &nhncloud.Config{
		Region:              "kr1",
		Credentials:         creds,
		IdentityCredentials: identityCreds,
		AppKeys: map[string]string{
			"rds-mysql": mysqlAppKey,
		},
		Debug: os.Getenv("DEBUG") == "true",
	}

	client, err := nhncloud.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("=== IAM Organizations ===")
	orgs, err := client.IAM().ListOrganizations(ctx)
	if err != nil {
		log.Printf("Failed to list organizations: %v", err)
	} else {
		for _, org := range orgs.Organizations() {
			fmt.Printf("  - %s (ID: %s, Status: %s)\n", org.Name, org.ID, org.Status)
		}
	}

	if identityCreds != nil {
		fmt.Println("\n=== Compute Servers ===")
		servers, err := client.Compute().ListServers(ctx)
		if err != nil {
			log.Printf("Failed to list servers: %v", err)
		} else {
			if len(servers.Servers) == 0 {
				fmt.Println("  No servers found")
			}
			for _, server := range servers.Servers {
				fmt.Printf("  - %s (ID: %s, Status: %s)\n", server.Name, server.ID, server.Status)
			}
		}

		fmt.Println("\n=== Compute Flavors ===")
		flavors, err := client.Compute().ListFlavors(ctx)
		if err != nil {
			log.Printf("Failed to list flavors: %v", err)
		} else {
			for i, flavor := range flavors.Flavors {
				if i >= 5 {
					fmt.Printf("  ... and %d more\n", len(flavors.Flavors)-5)
					break
				}
				fmt.Printf("  - %s (VCPUs: %d, RAM: %dMB)\n", flavor.Name, flavor.VCPUs, flavor.RAM)
			}
		}
	}

	if mysqlAppKey != "" {
		fmt.Println("\n=== RDS MySQL Instances ===")
		instances, err := client.MySQL().ListInstances(ctx)
		if err != nil {
			log.Printf("Failed to list MySQL instances: %v", err)
		} else {
			if len(instances.DBInstances) == 0 {
				fmt.Println("  No MySQL instances found")
			}
			for _, inst := range instances.DBInstances {
				fmt.Printf("  - %s (ID: %s, Status: %s, Version: %s)\n",
					inst.DBInstanceName, inst.DBInstanceID, inst.DBInstanceStatus, inst.DBVersion)
			}
		}

		fmt.Println("\n=== RDS MySQL Flavors ===")
		mysqlFlavors, err := client.MySQL().ListFlavors(ctx)
		if err != nil {
			log.Printf("Failed to list MySQL flavors: %v", err)
		} else {
			for i, flavor := range mysqlFlavors.DBFlavors {
				if i >= 5 {
					fmt.Printf("  ... and %d more\n", len(mysqlFlavors.DBFlavors)-5)
					break
				}
				fmt.Printf("  - %s (VCPUs: %d, RAM: %dMB)\n", flavor.FlavorName, flavor.Vcpus, flavor.Ram)
			}
		}
	}

	fmt.Println("\nDone!")
}
