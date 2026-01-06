package ncs

type ResponseHeader struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

type Workload struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Description       string            `json:"description,omitempty"`
	Type              string            `json:"type"`
	Replicas          int               `json:"replicas"`
	AvailableReplicas int               `json:"availableReplicas"`
	Status            string            `json:"status"`
	Containers        []Container       `json:"containers,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

type Container struct {
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	Command      []string              `json:"command,omitempty"`
	Args         []string              `json:"args,omitempty"`
	Ports        []ContainerPort       `json:"ports,omitempty"`
	Env          []EnvVar              `json:"env,omitempty"`
	Resources    *ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts []VolumeMount         `json:"volumeMounts,omitempty"`
}

type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	GPU    string `json:"nvidia.com/gpu,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type ListWorkloadsOutput struct {
	Header    *ResponseHeader `json:"header"`
	Workloads []Workload      `json:"workloads"`
}

type GetWorkloadOutput struct {
	Header *ResponseHeader `json:"header"`
	Workload
}

type CreateWorkloadInput struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Replicas    int               `json:"replicas,omitempty"`
	Containers  []Container       `json:"containers"`
	Labels      map[string]string `json:"labels,omitempty"`
	Volumes     []Volume          `json:"volumes,omitempty"`
}

type Volume struct {
	Name                  string                 `json:"name"`
	EmptyDir              *EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	ConfigMap             *ConfigMapVolumeSource `json:"configMap,omitempty"`
	Secret                *SecretVolumeSource    `json:"secret,omitempty"`
	PersistentVolumeClaim *PVCVolumeSource       `json:"persistentVolumeClaim,omitempty"`
}

type EmptyDirVolumeSource struct {
	Medium    string `json:"medium,omitempty"`
	SizeLimit string `json:"sizeLimit,omitempty"`
}

type ConfigMapVolumeSource struct {
	Name string `json:"name"`
}

type SecretVolumeSource struct {
	SecretName string `json:"secretName"`
}

type PVCVolumeSource struct {
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type CreateWorkloadOutput struct {
	Header *ResponseHeader `json:"header"`
	Workload
}

type UpdateWorkloadInput struct {
	Description string      `json:"description,omitempty"`
	Replicas    int         `json:"replicas,omitempty"`
	Containers  []Container `json:"containers,omitempty"`
}

type Template struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Version     string      `json:"version"`
	Type        string      `json:"type"`
	IsPublic    bool        `json:"isPublic"`
	Containers  []Container `json:"containers,omitempty"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type ListTemplatesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Templates []Template      `json:"templates"`
}

type GetTemplateOutput struct {
	Header *ResponseHeader `json:"header"`
	Template
}

type Service struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Type       string            `json:"type"`
	ClusterIP  string            `json:"clusterIP,omitempty"`
	ExternalIP string            `json:"externalIP,omitempty"`
	Ports      []ServicePort     `json:"ports,omitempty"`
	Selector   map[string]string `json:"selector,omitempty"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
}

type ServicePort struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	NodePort   int    `json:"nodePort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}

type ListServicesOutput struct {
	Header   *ResponseHeader `json:"header"`
	Services []Service       `json:"services"`
}

type GetServiceOutput struct {
	Header *ResponseHeader `json:"header"`
	Service
}

type CreateServiceInput struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace,omitempty"`
	Type      string            `json:"type"`
	Ports     []ServicePort     `json:"ports"`
	Selector  map[string]string `json:"selector,omitempty"`
}

type CreateServiceOutput struct {
	Header *ResponseHeader `json:"header"`
	Service
}
