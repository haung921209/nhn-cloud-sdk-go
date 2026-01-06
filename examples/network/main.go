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

	fmt.Println("=== VPCs ===")
	vpcs, err := client.VPC().ListVPCs(ctx)
	if err != nil {
		log.Printf("Failed to list VPCs: %v", err)
	} else {
		for _, v := range vpcs.VPCs {
			fmt.Printf("  - %s (ID: %s, CIDR: %s, State: %s)\n", v.Name, v.ID, v.CIDRv4, v.State)
		}
	}

	fmt.Println("\n=== Subnets ===")
	subnets, err := client.VPC().ListSubnets(ctx)
	if err != nil {
		log.Printf("Failed to list subnets: %v", err)
	} else {
		for _, s := range subnets.Subnets {
			fmt.Printf("  - %s (ID: %s, CIDR: %s)\n", s.Name, s.ID, s.CIDR)
		}
	}

	fmt.Println("\n=== Security Groups ===")
	sgs, err := client.SecurityGroup().ListSecurityGroups(ctx)
	if err != nil {
		log.Printf("Failed to list security groups: %v", err)
	} else {
		for _, sg := range sgs.SecurityGroups {
			fmt.Printf("  - %s (ID: %s, Rules: %d)\n", sg.Name, sg.ID, len(sg.Rules))
		}
	}

	fmt.Println("\n=== Floating IPs ===")
	fips, err := client.FloatingIP().ListFloatingIPs(ctx)
	if err != nil {
		log.Printf("Failed to list floating IPs: %v", err)
	} else {
		if len(fips.FloatingIPs) == 0 {
			fmt.Println("  No floating IPs found")
		}
		for _, fip := range fips.FloatingIPs {
			fmt.Printf("  - %s (Status: %s, Port: %v)\n", fip.FloatingIPAddress, fip.Status, fip.PortID)
		}
	}

	fmt.Println("\n=== Load Balancers ===")
	lbs, err := client.LoadBalancer().ListLoadBalancers(ctx)
	if err != nil {
		log.Printf("Failed to list load balancers: %v", err)
	} else {
		if len(lbs.LoadBalancers) == 0 {
			fmt.Println("  No load balancers found")
		}
		for _, lb := range lbs.LoadBalancers {
			fmt.Printf("  - %s (ID: %s, VIP: %s, Status: %s)\n", lb.Name, lb.ID, lb.VIPAddress, lb.ProvisioningStatus)
		}
	}

	fmt.Println("\nDone!")
}
