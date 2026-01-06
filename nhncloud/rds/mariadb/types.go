package mariadb

// ResponseHeader represents the common API response header
type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

// Instance represents a MariaDB database instance
type Instance struct {
	ID                    string   `json:"dbInstanceId"`
	Name                  string   `json:"dbInstanceName"`
	Status                string   `json:"dbInstanceStatus"`
	Description           string   `json:"description,omitempty"`
	Version               string   `json:"dbVersion"`
	Port                  int      `json:"dbPort"`
	StorageType           string   `json:"storageType"`
	StorageSize           int      `json:"storageSize"`
	SubnetID              string   `json:"subnetId"`
	SecurityGroupIDs      []string `json:"dbSecurityGroupIds,omitempty"`
	FlavorID              string   `json:"dbFlavorId"`
	ParameterGroupID      string   `json:"parameterGroupId,omitempty"`
	AuthenticationPlugin  string   `json:"authenticationPlugin,omitempty"`
	TLSOption             string   `json:"tlsOption,omitempty"`
	UseDeletionProtection bool     `json:"useDeletionProtection,omitempty"`
	CreatedAt             string   `json:"createdYmdt"`
	UpdatedAt             string   `json:"updatedYmdt"`
}

// InstanceGroup represents a MariaDB instance group (for HA)
type InstanceGroup struct {
	ID              string `json:"dbInstanceGroupId"`
	ReplicationType string `json:"replicationType"`
	CreatedAt       string `json:"createdYmdt"`
	UpdatedAt       string `json:"updatedYmdt"`
}

// ListInstancesOutput represents the response for listing instances
type ListInstancesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Instances []Instance      `json:"dbInstances"`
}

// ListInstanceGroupsOutput represents the response for listing instance groups
type ListInstanceGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	InstanceGroups []InstanceGroup `json:"dbInstanceGroups"`
}

// GetInstanceOutput represents the response for getting a single instance
type GetInstanceOutput struct {
	Header *ResponseHeader `json:"header"`
	Instance
}

// Network represents network configuration for an instance
type Network struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
}

// Storage represents storage configuration for an instance
type Storage struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

// BackupSchedule represents a backup time window
type BackupSchedule struct {
	BeginTime string `json:"backupWndBgnTime"`
	Duration  string `json:"backupWndDuration"`
}

// BackupConfig represents backup configuration
type BackupConfig struct {
	Period    int              `json:"backupPeriod"`
	Schedules []BackupSchedule `json:"backupSchedules"`
}

// CreateInstanceInput represents the input for creating a MariaDB instance
type CreateInstanceInput struct {
	Name                  string        `json:"dbInstanceName"`
	CandidateName         string        `json:"dbInstanceCandidateName,omitempty"`
	Description           string        `json:"description,omitempty"`
	FlavorID              string        `json:"dbFlavorId"`
	Version               string        `json:"dbVersion"`
	UserName              string        `json:"dbUserName"`
	Password              string        `json:"dbPassword"`
	Port                  int           `json:"dbPort,omitempty"`
	ParameterGroupID      string        `json:"parameterGroupId"`
	SecurityGroupIDs      []string      `json:"dbSecurityGroupIds,omitempty"`
	Network               *Network      `json:"network"`
	Storage               *Storage      `json:"storage"`
	Backup                *BackupConfig `json:"backup"`
	UseHighAvailability   bool          `json:"useHighAvailability,omitempty"`
	ReplicationMode       string        `json:"replicationMode,omitempty"`
	UseDeletionProtection bool          `json:"useDeletionProtection,omitempty"`
	AuthenticationPlugin  string        `json:"authenticationPlugin,omitempty"`
	TLSOption             string        `json:"tlsOption,omitempty"`
}

// CreateInstanceOutput represents the response for creating an instance
type CreateInstanceOutput struct {
	Header *ResponseHeader `json:"header"`
	Instance
}

// ModifyInstanceInput represents the input for modifying an instance
type ModifyInstanceInput struct {
	Name              string   `json:"dbInstanceName,omitempty"`
	CandidateName     string   `json:"dbInstanceCandidateName,omitempty"`
	Description       string   `json:"description,omitempty"`
	Port              int      `json:"dbPort,omitempty"`
	Version           string   `json:"dbVersion,omitempty"`
	FlavorID          string   `json:"dbFlavorId,omitempty"`
	ParameterGroupID  string   `json:"parameterGroupId,omitempty"`
	SecurityGroupIDs  []string `json:"dbSecurityGroupIds,omitempty"`
	ExecuteBackup     bool     `json:"executeBackup,omitempty"`
	UseOnlineFailover bool     `json:"useOnlineFailover,omitempty"`
}

// JobOutput represents a job response (for async operations)
type JobOutput struct {
	Header *ResponseHeader `json:"header"`
	JobID  string          `json:"jobId"`
}

// Flavor represents a database flavor (instance type)
type Flavor struct {
	ID    string `json:"dbFlavorId"`
	Name  string `json:"dbFlavorName"`
	RAM   int    `json:"ram"`
	VCPUs int    `json:"vcpus"`
	Disk  int    `json:"disk"`
}

// ListFlavorsOutput represents the response for listing flavors
type ListFlavorsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Flavors []Flavor        `json:"dbFlavors"`
}

// Version represents a database version
type Version struct {
	Version string `json:"dbVersion"`
	Name    string `json:"dbVersionName"`
}

// ListVersionsOutput represents the response for listing versions
type ListVersionsOutput struct {
	Header   *ResponseHeader `json:"header"`
	Versions []Version       `json:"dbVersions"`
}

// SecurityGroup represents a database security group
type SecurityGroup struct {
	ID          string         `json:"dbSecurityGroupId"`
	Name        string         `json:"dbSecurityGroupName"`
	Description string         `json:"description,omitempty"`
	Status      string         `json:"progressStatus"`
	Rules       []SecurityRule `json:"rules,omitempty"`
	CreatedAt   string         `json:"createdYmdt"`
	UpdatedAt   string         `json:"updatedYmdt"`
}

// SecurityRule represents a security group rule
type SecurityRule struct {
	ID          string `json:"ruleId"`
	Description string `json:"description,omitempty"`
	Direction   string `json:"direction"`
	EtherType   string `json:"etherType"`
	Port        Port   `json:"port"`
	CIDR        string `json:"cidr"`
	CreatedAt   string `json:"createdYmdt"`
	UpdatedAt   string `json:"updatedYmdt"`
}

// Port represents port configuration in a security rule
type Port struct {
	PortType string `json:"portType"`
	MinPort  *int   `json:"minPort,omitempty"`
	MaxPort  *int   `json:"maxPort,omitempty"`
}

// ListSecurityGroupsOutput represents the response for listing security groups
type ListSecurityGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	SecurityGroups []SecurityGroup `json:"dbSecurityGroups"`
}

// ParameterGroup represents a parameter group
type ParameterGroup struct {
	ID          string      `json:"parameterGroupId"`
	Name        string      `json:"parameterGroupName"`
	Description string      `json:"description,omitempty"`
	Version     string      `json:"dbVersion"`
	Status      string      `json:"parameterGroupStatus"`
	Parameters  []Parameter `json:"parameters,omitempty"`
	CreatedAt   string      `json:"createdYmdt"`
	UpdatedAt   string      `json:"updatedYmdt"`
}

// Parameter represents a database parameter
type Parameter struct {
	ID           string `json:"parameterId"`
	Name         string `json:"parameterName"`
	Value        string `json:"value"`
	DefaultValue string `json:"defaultValue"`
	AllowedValue string `json:"allowedValue"`
	UpdateType   string `json:"updateType"`
	ApplyType    string `json:"applyType"`
}

// ListParameterGroupsOutput represents the response for listing parameter groups
type ListParameterGroupsOutput struct {
	Header          *ResponseHeader  `json:"header"`
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

// Backup represents a database backup
type Backup struct {
	ID         string `json:"backupId"`
	Name       string `json:"backupName"`
	Status     string `json:"backupStatus"`
	InstanceID string `json:"dbInstanceId"`
	Version    string `json:"dbVersion"`
	Type       string `json:"backupType"`
	Size       int64  `json:"backupSize"`
	CreatedAt  string `json:"createdYmdt"`
	UpdatedAt  string `json:"updatedYmdt"`
}

// ListBackupsOutput represents the response for listing backups
type ListBackupsOutput struct {
	Header     *ResponseHeader `json:"header"`
	TotalCount int             `json:"totalCounts"`
	Backups    []Backup        `json:"backups"`
}

// Subnet represents a subnet
type Subnet struct {
	ID               string `json:"subnetId"`
	Name             string `json:"subnetName"`
	CIDR             string `json:"subnetCidr"`
	UsingGateway     bool   `json:"usingGateway"`
	AvailableIPCount int    `json:"availableIpCount"`
}

// ListSubnetsOutput represents the response for listing subnets
type ListSubnetsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []Subnet        `json:"subnets"`
}

// NetworkEndpoint represents a network endpoint
type NetworkEndpoint struct {
	Domain       string `json:"domain"`
	IPAddress    string `json:"ipAddress"`
	EndpointType string `json:"endPointType"`
}

// NetworkInfo represents network information for an instance
type NetworkInfo struct {
	Header           *ResponseHeader   `json:"header"`
	AvailabilityZone string            `json:"availabilityZone"`
	Subnet           Subnet            `json:"subnet"`
	Endpoints        []NetworkEndpoint `json:"endPoints"`
}
