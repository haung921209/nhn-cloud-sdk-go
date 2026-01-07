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
	mysqlAppKey := os.Getenv("NHN_CLOUD_MYSQL_APPKEY")
	mariadbAppKey := os.Getenv("NHN_CLOUD_MARIADB_APPKEY")
	pgAppKey := os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY")

	if accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("NHN_CLOUD_ACCESS_KEY_ID and NHN_CLOUD_SECRET_ACCESS_KEY are required")
	}

	creds := credentials.NewStatic(accessKeyID, secretAccessKey)

	cfg := &nhncloud.Config{
		Region:      "kr1",
		Credentials: creds,
		AppKeys: map[string]string{
			"rds-mysql":      mysqlAppKey,
			"rds-mariadb":    mariadbAppKey,
			"rds-postgresql": pgAppKey,
		},
		Debug: os.Getenv("DEBUG") == "true",
	}

	client, err := nhncloud.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	if mysqlAppKey != "" {
		fmt.Println("=== MySQL Instances ===")
		instances, err := client.MySQL().ListInstances(ctx)
		if err != nil {
			log.Printf("Failed to list MySQL instances: %v", err)
		} else {
			if len(instances.DBInstances) == 0 {
				fmt.Println("  No MySQL instances found")
			}
			for _, inst := range instances.DBInstances {
				fmt.Printf("  - %s (ID: %s, Status: %s, Version: %s, Port: %d)\n",
					inst.DBInstanceName, inst.DBInstanceID, inst.DBInstanceStatus, inst.DBVersion, inst.DBPort)
			}
		}

		fmt.Println("\n=== MySQL Flavors ===")
		flavors, err := client.MySQL().ListFlavors(ctx)
		if err != nil {
			log.Printf("Failed to list MySQL flavors: %v", err)
		} else {
			for i, f := range flavors.DBFlavors {
				if i >= 5 {
					fmt.Printf("  ... and %d more\n", len(flavors.DBFlavors)-5)
					break
				}
				fmt.Printf("  - %s (VCPUs: %d, RAM: %dMB)\n", f.FlavorName, f.Vcpus, f.Ram)
			}
		}

		fmt.Println("\n=== MySQL Versions ===")
		versions, err := client.MySQL().ListVersions(ctx)
		if err != nil {
			log.Printf("Failed to list MySQL versions: %v", err)
		} else {
			for _, v := range versions.DBVersions {
				fmt.Printf("  - %s (%s)\n", v.DBVersion, v.DBVersionName)
			}
		}
	}

	if mariadbAppKey != "" {
		fmt.Println("\n=== MariaDB Instances ===")
		instances, err := client.MariaDB().ListInstances(ctx)
		if err != nil {
			log.Printf("Failed to list MariaDB instances: %v", err)
		} else {
			if len(instances.DBInstances) == 0 {
				fmt.Println("  No MariaDB instances found")
			}
			for _, inst := range instances.DBInstances {
				fmt.Printf("  - %s (ID: %s, Status: %s, Version: %s)\n",
					inst.DBInstanceName, inst.DBInstanceID, inst.DBInstanceStatus, inst.DBVersion)
			}
		}
	}

	if pgAppKey != "" {
		fmt.Println("\n=== PostgreSQL Instances ===")
		instances, err := client.PostgreSQL().ListInstances(ctx)
		if err != nil {
			log.Printf("Failed to list PostgreSQL instances: %v", err)
		} else {
			if len(instances.DBInstances) == 0 {
				fmt.Println("  No PostgreSQL instances found")
			}
			for _, inst := range instances.DBInstances {
				fmt.Printf("  - %s (ID: %s, Status: %s, Version: %s)\n",
					inst.DBInstanceName, inst.DBInstanceID, inst.DBInstanceStatus, inst.DBVersion)
			}
		}

	}

	fmt.Println("\nDone!")
}
