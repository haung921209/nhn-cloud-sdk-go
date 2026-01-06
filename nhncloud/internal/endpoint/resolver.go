package endpoint

import (
	"fmt"
	"strings"
)

type ServiceType string

const (
	ServiceIAM        ServiceType = "iam"
	ServiceCompute    ServiceType = "compute"
	ServiceNetwork    ServiceType = "network"
	ServiceRDSMySQL   ServiceType = "rds-mysql"
	ServiceRDSMariaDB ServiceType = "rds-mariadb"
	ServiceRDSPG      ServiceType = "rds-postgresql"
	ServiceVPC        ServiceType = "vpc"
	ServiceBlock      ServiceType = "block-storage"
	ServiceObject     ServiceType = "object-storage"
	ServiceNKS        ServiceType = "nks"
	ServiceNCR        ServiceType = "ncr"
	ServiceNCS        ServiceType = "ncs"
)

var baseURLTemplates = map[ServiceType]string{
	ServiceIAM:        "https://iam.api.nhncloudservice.com",
	ServiceRDSMySQL:   "https://%s-rds-mysql.api.nhncloudservice.com",
	ServiceRDSMariaDB: "https://%s-rds-mariadb.api.nhncloudservice.com",
	ServiceRDSPG:      "https://%s-rds-postgres.api.nhncloudservice.com",
}

func Resolve(serviceType ServiceType, region string) string {
	region = strings.ToLower(region)

	template, ok := baseURLTemplates[serviceType]
	if !ok {
		return ""
	}

	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, region)
	}

	return template
}

func ResolveWithAppKey(serviceType ServiceType, region, appKey string) string {
	baseURL := Resolve(serviceType, region)
	if baseURL == "" || appKey == "" {
		return baseURL
	}

	switch serviceType {
	case ServiceRDSMySQL, ServiceRDSMariaDB, ServiceRDSPG:
		return baseURL + "/v3.0/appkeys/" + appKey
	case ServiceNCR, ServiceNCS:
		return baseURL + "/v2.0/appkeys/" + appKey
	default:
		return baseURL
	}
}
