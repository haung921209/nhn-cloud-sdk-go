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
	username := os.Getenv("NHN_CLOUD_USERNAME")
	password := os.Getenv("NHN_CLOUD_PASSWORD")
	tenantID := os.Getenv("NHN_CLOUD_TENANT_ID")

	if username == "" || password == "" || tenantID == "" {
		log.Fatal("NHN_CLOUD_USERNAME, NHN_CLOUD_PASSWORD, and NHN_CLOUD_TENANT_ID are required")
	}

	identityCreds := credentials.NewStaticIdentity(username, password, tenantID)

	cfg := &nhncloud.Config{
		Region:              "kr1",
		IdentityCredentials: identityCreds,
		Debug:               os.Getenv("DEBUG") == "true",
	}

	client, err := nhncloud.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("=== NKS Cluster Templates ===")
	templates, err := client.NKS().ListClusterTemplates(ctx)
	if err != nil {
		log.Printf("Failed to list cluster templates: %v", err)
	} else {
		for _, t := range templates.ClusterTemplates {
			fmt.Printf("  - %s (ID: %s, COE: %s)\n", t.Name, t.ID, t.COE)
		}
	}

	fmt.Println("\n=== NKS Clusters ===")
	clusters, err := client.NKS().ListClusters(ctx)
	if err != nil {
		log.Printf("Failed to list clusters: %v", err)
	} else {
		if len(clusters.Clusters) == 0 {
			fmt.Println("  No clusters found")
		}
		for _, c := range clusters.Clusters {
			fmt.Printf("  - %s (ID: %s, Status: %s, K8s: %s, Nodes: %d)\n",
				c.Name, c.ID, c.Status, c.K8sVersion, c.NodeCount)

			nodeGroups, err := client.NKS().ListNodeGroups(ctx, c.ID)
			if err != nil {
				log.Printf("    Failed to list node groups: %v", err)
				continue
			}
			for _, ng := range nodeGroups.NodeGroups {
				fmt.Printf("    - NodeGroup: %s (Nodes: %d, Min: %d, Max: %d)\n",
					ng.Name, ng.NodeCount, ng.MinNodeCount, ng.MaxNodeCount)
			}
		}
	}

	fmt.Println("\nDone!")
}
