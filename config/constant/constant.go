package constant

type WithNameSpace = string
type WithoutNameSpace = string

const (
	ControllerRevisions      WithNameSpace = "ControllerRevision"
	DaemonSets               WithNameSpace = "DaemonSet"
	Deployments              WithNameSpace = "Deployment"
	ReplicaSets              WithNameSpace = "ReplicaSet"
	StatefulSets             WithNameSpace = "StatefulSet"
	HorizontalPodAutoscalers WithNameSpace = "HorizontalPodAutoscaler"
	Jobs                     WithNameSpace = "Job"
	CronJobs                 WithNameSpace = "CronJob"
	Leases                   WithNameSpace = "Lease"
	ConfigMaps               WithNameSpace = "ConfigMap"
	Endpoints                WithNameSpace = "Endpoint"
	Events                   WithNameSpace = "Event"
	LimitRanges              WithNameSpace = "LimitRange"
	PersistentVolumeClaims   WithNameSpace = "PersistentVolumeClaim"
	Pods                     WithNameSpace = "Pod"
	PodTemplates             WithNameSpace = "PodTemplate"
	ReplicationControllers   WithNameSpace = "ReplicationController"
	ResourceQuotas           WithNameSpace = "ResourceQuota"
	Secrets                  WithNameSpace = "Secret"
	Services                 WithNameSpace = "Service"
	ServiceAccounts          WithNameSpace = "ServiceAccount"
	EndpointSlices           WithNameSpace = "EndpointSlice"
	Ingresses                WithNameSpace = "Ingress"
	NetworkPolicies          WithNameSpace = "NetworkPolicy"
	PodDisruptionBudgets     WithNameSpace = "PodDisruptionBudget"
	Roles                    WithNameSpace = "Role"
	RoleBindings             WithNameSpace = "RoleBinding"
	PodPresets               WithNameSpace = "PodPreset"
)

const (
	MutatingWebhookConfigurations   WithoutNameSpace = "MutatingWebhookConfiguration"
	ValidatingWebhookConfigurations WithoutNameSpace = "ValidatingWebhookConfiguration"
	AuditSinks                      WithoutNameSpace = "AuditSink"
	CertificateSigningRequests      WithoutNameSpace = "CertificateSigningRequest"
	ComponentStatuses               WithoutNameSpace = "ComponentStatus"
	Namespaces                      WithoutNameSpace = "Namespace"
	Nodes                           WithoutNameSpace = "Node"
	PersistentVolumes               WithoutNameSpace = "PersistentVolume"
	PodSecurityPolicies             WithoutNameSpace = "PodSecurityPolicy"
	FlowSchemas                     WithoutNameSpace = "FlowSchema"
	PriorityLevelConfigurations     WithoutNameSpace = "PriorityLevelConfiguration"
	IngressClasses                  WithoutNameSpace = "IngressClass"
	RuntimeClasses                  WithoutNameSpace = "RuntimeClass"
	ClusterRoles                    WithoutNameSpace = "ClusterRole"
	ClusterRoleBindings             WithoutNameSpace = "ClusterRoleBinding"
	PriorityClasses                 WithoutNameSpace = "PriorityClass"
	CSIDrivers                      WithoutNameSpace = "CSIDriver"
	CSINodes                        WithoutNameSpace = "CSINode"
	StorageClasses                  WithoutNameSpace = "StorageClass"
	VolumeAttachments               WithoutNameSpace = "VolumeAttachment"
)

var KubeKinds = []string{
	// WithNameSpace
	ControllerRevisions,
	DaemonSets,
	Deployments,
	ReplicaSets,
	StatefulSets,
	HorizontalPodAutoscalers,
	Jobs,
	CronJobs,
	Leases,
	ConfigMaps,
	Endpoints,
	Events,
	LimitRanges,
	PersistentVolumeClaims,
	Pods,
	PodTemplates,
	ReplicationControllers,
	ResourceQuotas,
	Secrets,
	Services,
	ServiceAccounts,
	EndpointSlices,
	Ingresses,
	NetworkPolicies,
	PodDisruptionBudgets,
	Roles,
	RoleBindings,
	PodPresets,
	// WithoutNameSpace
	MutatingWebhookConfigurations,
	ValidatingWebhookConfigurations,
	AuditSinks,
	CertificateSigningRequests,
	ComponentStatuses,
	Namespaces,
	Nodes,
	PersistentVolumes,
	PodSecurityPolicies,
	FlowSchemas,
	PriorityLevelConfigurations,
	IngressClasses,
	RuntimeClasses,
	ClusterRoles,
	ClusterRoleBindings,
	PriorityClasses,
	CSIDrivers,
	CSINodes,
	StorageClasses,
	VolumeAttachments,
}
