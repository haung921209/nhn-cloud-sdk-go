package mysql

type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

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

type InstanceGroup struct {
	ID              string `json:"dbInstanceGroupId"`
	ReplicationType string `json:"replicationType"`
	CreatedAt       string `json:"createdYmdt"`
	UpdatedAt       string `json:"updatedYmdt"`
}

type ListInstancesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Instances []Instance      `json:"dbInstances"`
}

type ListInstanceGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	InstanceGroups []InstanceGroup `json:"dbInstanceGroups"`
}

type GetInstanceOutput struct {
	Header *ResponseHeader `json:"header"`
	Instance
}

type Network struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
}

type Storage struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

type BackupSchedule struct {
	BeginTime string `json:"backupWndBgnTime"`
	Duration  string `json:"backupWndDuration"`
}

type BackupConfig struct {
	Period    int              `json:"backupPeriod"`
	Schedules []BackupSchedule `json:"backupSchedules"`
}

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

type CreateInstanceOutput struct {
	Header *ResponseHeader `json:"header"`
	Instance
}

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

type JobOutput struct {
	Header *ResponseHeader `json:"header"`
	JobID  string          `json:"jobId"`
}

type Flavor struct {
	ID    string `json:"dbFlavorId"`
	Name  string `json:"dbFlavorName"`
	RAM   int    `json:"ram"`
	VCPUs int    `json:"vcpus"`
	Disk  int    `json:"disk"`
}

type ListFlavorsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Flavors []Flavor        `json:"dbFlavors"`
}

type Version struct {
	Version string `json:"dbVersion"`
	Name    string `json:"dbVersionName"`
}

type ListVersionsOutput struct {
	Header   *ResponseHeader `json:"header"`
	Versions []Version       `json:"dbVersions"`
}

type SecurityGroup struct {
	ID          string         `json:"dbSecurityGroupId"`
	Name        string         `json:"dbSecurityGroupName"`
	Description string         `json:"description,omitempty"`
	Status      string         `json:"progressStatus"`
	Rules       []SecurityRule `json:"rules,omitempty"`
	CreatedAt   string         `json:"createdYmdt"`
	UpdatedAt   string         `json:"updatedYmdt"`
}

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

type Port struct {
	PortType string `json:"portType"`
	MinPort  *int   `json:"minPort,omitempty"`
	MaxPort  *int   `json:"maxPort,omitempty"`
}

type ListSecurityGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	SecurityGroups []SecurityGroup `json:"dbSecurityGroups"`
}

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

type Parameter struct {
	ID           string `json:"parameterId"`
	Name         string `json:"parameterName"`
	Value        string `json:"value"`
	DefaultValue string `json:"defaultValue"`
	AllowedValue string `json:"allowedValue"`
	UpdateType   string `json:"updateType"`
	ApplyType    string `json:"applyType"`
}

type ListParameterGroupsOutput struct {
	Header          *ResponseHeader  `json:"header"`
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

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

type ListBackupsOutput struct {
	Header     *ResponseHeader `json:"header"`
	TotalCount int             `json:"totalCounts"`
	Backups    []Backup        `json:"backups"`
}

type Subnet struct {
	ID               string `json:"subnetId"`
	Name             string `json:"subnetName"`
	CIDR             string `json:"subnetCidr"`
	UsingGateway     bool   `json:"usingGateway"`
	AvailableIPCount int    `json:"availableIpCount"`
}

type ListSubnetsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []Subnet        `json:"subnets"`
}

type NetworkEndpoint struct {
	Domain       string `json:"domain"`
	IPAddress    string `json:"ipAddress"`
	EndpointType string `json:"endPointType"`
}

type NetworkInfo struct {
	Header           *ResponseHeader   `json:"header"`
	AvailabilityZone string            `json:"availabilityZone"`
	Subnet           Subnet            `json:"subnet"`
	Endpoints        []NetworkEndpoint `json:"endPoints"`
}
