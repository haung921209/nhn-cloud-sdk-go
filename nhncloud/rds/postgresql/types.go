package postgresql

type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

type Instance struct {
	ID                    string   `json:"dbInstanceId"`
	Name                  string   `json:"dbInstanceName"`
	Status                string   `json:"dbInstanceStatus"`
	Type                  string   `json:"dbInstanceType,omitempty"`
	Description           string   `json:"description,omitempty"`
	Version               string   `json:"dbVersion"`
	Port                  int      `json:"dbPort"`
	StorageType           string   `json:"storageType"`
	StorageSize           int      `json:"storageSize"`
	SubnetID              string   `json:"subnetId"`
	SecurityGroupIDs      []string `json:"dbSecurityGroupIds,omitempty"`
	FlavorID              string   `json:"dbFlavorId"`
	ParameterGroupID      string   `json:"parameterGroupId,omitempty"`
	UseDeletionProtection bool     `json:"useDeletionProtection,omitempty"`
	UseHighAvailability   bool     `json:"useHighAvailability,omitempty"`
	CreatedAt             string   `json:"createdYmdt"`
	UpdatedAt             string   `json:"updatedYmdt"`
}

type ListInstancesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Instances []Instance      `json:"dbInstances"`
}

type GetInstanceOutput struct {
	Header *ResponseHeader `json:"header"`
	Instance
}

type CreateInstanceInput struct {
	Name                  string         `json:"dbInstanceName"`
	Description           string         `json:"description,omitempty"`
	FlavorID              string         `json:"dbFlavorId"`
	Version               string         `json:"dbVersion"`
	UserName              string         `json:"dbUserName"`
	Password              string         `json:"dbUserPassword"`
	Port                  int            `json:"dbPort,omitempty"`
	ParameterGroupID      string         `json:"parameterGroupId"`
	SecurityGroupIDs      []string       `json:"dbSecurityGroupIds,omitempty"`
	Network               *NetworkConfig `json:"network"`
	Storage               *StorageConfig `json:"storage"`
	Backup                *BackupConfig  `json:"backup,omitempty"`
	UseHighAvailability   bool           `json:"useHighAvailability,omitempty"`
	UseDeletionProtection bool           `json:"useDeletionProtection,omitempty"`
}

type NetworkConfig struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
}

type StorageConfig struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

type BackupConfig struct {
	BackupPeriod    int              `json:"backupPeriod"`
	BackupSchedules []BackupSchedule `json:"backupSchedules,omitempty"`
}

type BackupSchedule struct {
	BackupWndBgnTime  string `json:"backupWndBgnTime"`
	BackupWndDuration string `json:"backupWndDuration"`
}

type ModifyInstanceInput struct {
	Name             string   `json:"dbInstanceName,omitempty"`
	Description      string   `json:"description,omitempty"`
	Port             int      `json:"dbPort,omitempty"`
	FlavorID         string   `json:"dbFlavorId,omitempty"`
	ParameterGroupID string   `json:"parameterGroupId,omitempty"`
	SecurityGroupIDs []string `json:"dbSecurityGroupIds,omitempty"`
}

type JobOutput struct {
	Header *ResponseHeader `json:"header"`
	JobID  string          `json:"jobId,omitempty"`
}

type InstanceGroup struct {
	ID              string     `json:"dbInstanceGroupId"`
	ReplicationType string     `json:"replicationType"`
	Instances       []Instance `json:"dbInstances,omitempty"`
	CreatedAt       string     `json:"createdYmdt"`
	UpdatedAt       string     `json:"updatedYmdt"`
}

type ListInstanceGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	InstanceGroups []InstanceGroup `json:"dbInstanceGroups"`
}

type InstanceGroupOutput struct {
	Header *ResponseHeader `json:"header"`
	InstanceGroup
}

type Flavor struct {
	ID    string `json:"dbFlavorId"`
	Name  string `json:"dbFlavorName"`
	RAM   int    `json:"ram"`
	VCPUs int    `json:"vcpus"`
}

type ListFlavorsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Flavors []Flavor        `json:"dbFlavors"`
}

type Version struct {
	DBVersion   string `json:"dbVersion"`
	DisplayName string `json:"dbVersionName"`
}

type ListVersionsOutput struct {
	Header   *ResponseHeader `json:"header"`
	Versions []Version       `json:"dbVersions"`
}

type StorageType struct {
	StorageType string `json:"storageType"`
	MinSize     int    `json:"minSize"`
	MaxSize     int    `json:"maxSize"`
}

type ListStorageTypesOutput struct {
	Header       *ResponseHeader `json:"header"`
	StorageTypes []StorageType   `json:"storageTypes"`
}

type SecurityGroup struct {
	ID          string              `json:"dbSecurityGroupId"`
	Name        string              `json:"dbSecurityGroupName"`
	Description string              `json:"description,omitempty"`
	Rules       []SecurityGroupRule `json:"rules,omitempty"`
	CreatedAt   string              `json:"createdYmdt,omitempty"`
	UpdatedAt   string              `json:"updatedYmdt,omitempty"`
}

type SecurityGroupRule struct {
	ID        string `json:"dbSecurityGroupRuleId"`
	Direction string `json:"direction"`
	EtherType string `json:"etherType"`
	Port      int    `json:"port,omitempty"`
	CIDR      string `json:"cidr,omitempty"`
}

type ListSecurityGroupsOutput struct {
	Header         *ResponseHeader `json:"header"`
	SecurityGroups []SecurityGroup `json:"dbSecurityGroups"`
}

type SecurityGroupOutput struct {
	Header *ResponseHeader `json:"header"`
	SecurityGroup
}

type CreateSecurityGroupInput struct {
	Name        string `json:"dbSecurityGroupName"`
	Description string `json:"description,omitempty"`
}

type SecurityGroupIDOutput struct {
	Header          *ResponseHeader `json:"header"`
	SecurityGroupID string          `json:"dbSecurityGroupId"`
}

type CreateSecurityGroupRuleInput struct {
	Direction string `json:"direction"`
	EtherType string `json:"etherType"`
	Port      int    `json:"port,omitempty"`
	CIDR      string `json:"cidr"`
}

type SecurityGroupRuleOutput struct {
	Header *ResponseHeader `json:"header"`
	RuleID string          `json:"dbSecurityGroupRuleId"`
}

type ParameterGroup struct {
	ID          string      `json:"parameterGroupId"`
	Name        string      `json:"parameterGroupName"`
	Description string      `json:"description,omitempty"`
	Version     string      `json:"dbVersion"`
	Status      string      `json:"parameterGroupStatus"`
	Parameters  []Parameter `json:"parameters,omitempty"`
	CreatedAt   string      `json:"createdYmdt,omitempty"`
	UpdatedAt   string      `json:"updatedYmdt,omitempty"`
}

type Parameter struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	DefaultValue string `json:"defaultValue,omitempty"`
	AllowedValue string `json:"allowedValue,omitempty"`
	UpdateType   string `json:"updateType,omitempty"`
}

type ListParameterGroupsOutput struct {
	Header          *ResponseHeader  `json:"header"`
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

type ParameterGroupOutput struct {
	Header *ResponseHeader `json:"header"`
	ParameterGroup
}

type CreateParameterGroupInput struct {
	Name        string `json:"parameterGroupName"`
	Description string `json:"description,omitempty"`
	DBVersion   string `json:"dbVersion"`
}

type ParameterGroupIDOutput struct {
	Header           *ResponseHeader `json:"header"`
	ParameterGroupID string          `json:"parameterGroupId"`
}

type ModifyParametersInput struct {
	Parameters []ParameterValue `json:"modifiedParameters"`
}

type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Backup struct {
	ID           string `json:"backupId"`
	Name         string `json:"backupName"`
	InstanceID   string `json:"dbInstanceId"`
	InstanceName string `json:"dbInstanceName"`
	Size         int64  `json:"backupSize"`
	Status       string `json:"backupStatus"`
	BackupType   string `json:"backupType,omitempty"`
	CreatedAt    string `json:"createdYmdt"`
}

type ListBackupsOutput struct {
	Header     *ResponseHeader `json:"header"`
	Backups    []Backup        `json:"backups"`
	TotalCount int             `json:"totalCount,omitempty"`
}

type CreateBackupInput struct {
	BackupName string `json:"backupName"`
}

type RestoreBackupInput struct {
	DBInstanceName string `json:"dbInstanceName"`
}

type DBUser struct {
	ID        string `json:"dbUserId"`
	Name      string `json:"dbUserName"`
	CreatedAt string `json:"createdYmdt,omitempty"`
	UpdatedAt string `json:"updatedYmdt,omitempty"`
}

type ListDBUsersOutput struct {
	Header  *ResponseHeader `json:"header"`
	DBUsers []DBUser        `json:"dbUsers"`
}

type CreateDBUserInput struct {
	UserName string `json:"dbUserName"`
	Password string `json:"dbUserPassword"`
}

type UpdateDBUserInput struct {
	Password string `json:"dbUserPassword,omitempty"`
}

type Database struct {
	ID        string `json:"databaseId"`
	Name      string `json:"databaseName"`
	CreatedAt string `json:"createdYmdt,omitempty"`
}

type ListDatabasesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Databases []Database      `json:"databases"`
}

type CreateDatabaseInput struct {
	DatabaseName string `json:"databaseName"`
}

type DatabaseIDOutput struct {
	Header     *ResponseHeader `json:"header"`
	DatabaseID string          `json:"databaseId"`
}

type Subnet struct {
	ID               string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	SubnetCIDR       string `json:"subnetCidr"`
	AvailabilityZone string `json:"availabilityZone"`
	VPCName          string `json:"vpcName,omitempty"`
}

type ListSubnetsOutput struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []Subnet        `json:"subnets"`
}

type NetworkInfo struct {
	SubnetID          string `json:"subnetId"`
	UsePublicAccess   bool   `json:"usePublicAccess"`
	AvailabilityZone  string `json:"availabilityZone,omitempty"`
	PublicDomainName  string `json:"publicDomainName,omitempty"`
	PrivateDomainName string `json:"privateDomainName,omitempty"`
}

type NetworkInfoOutput struct {
	Header *ResponseHeader `json:"header"`
	NetworkInfo
}

type ModifyStorageInfoInput struct {
	StorageSize int `json:"storageSize"`
}

type ModifyDeletionProtectionInput struct {
	UseDeletionProtection bool `json:"useDeletionProtection"`
}

type CreateReplicaInput struct {
	DBInstanceName   string `json:"dbInstanceName"`
	Description      string `json:"description,omitempty"`
	DBFlavorID       string `json:"dbFlavorId,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type EnableHAInput struct {
	UseHighAvailability bool `json:"useHighAvailability"`
	PingInterval        int  `json:"pingInterval,omitempty"`
}

type NotificationGroup struct {
	ID           string   `json:"notificationGroupId"`
	Name         string   `json:"notificationGroupName"`
	NotifyEmail  bool     `json:"notifyEmail"`
	NotifySMS    bool     `json:"notifySms"`
	IsEnabled    bool     `json:"isEnabled"`
	InstanceIDs  []string `json:"dbInstanceIds,omitempty"`
	UserGroupIDs []string `json:"userGroupIds,omitempty"`
	CreatedAt    string   `json:"createdYmdt,omitempty"`
	UpdatedAt    string   `json:"updatedYmdt,omitempty"`
}

type ListNotificationGroupsOutput struct {
	Header             *ResponseHeader     `json:"header"`
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

type NotificationGroupOutput struct {
	Header *ResponseHeader `json:"header"`
	NotificationGroup
}

type CreateNotificationGroupInput struct {
	Name         string   `json:"notificationGroupName"`
	NotifyEmail  bool     `json:"notifyEmail"`
	NotifySMS    bool     `json:"notifySms"`
	IsEnabled    bool     `json:"isEnabled"`
	InstanceIDs  []string `json:"dbInstanceIds,omitempty"`
	UserGroupIDs []string `json:"userGroupIds,omitempty"`
}

type NotificationGroupIDOutput struct {
	Header              *ResponseHeader `json:"header"`
	NotificationGroupID string          `json:"notificationGroupId"`
}

type LogFile struct {
	FileName     string `json:"logFileName"`
	FileSize     int64  `json:"logFileSize"`
	LastModified string `json:"lastModifiedYmdt"`
}

type ListLogFilesOutput struct {
	Header   *ResponseHeader `json:"header"`
	LogFiles []LogFile       `json:"logFiles"`
}
