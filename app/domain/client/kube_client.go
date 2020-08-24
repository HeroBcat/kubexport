package client

import (
	k8s "k8s.io/client-go/kubernetes"
)

type KubeClient interface {
	DefaultKubeClient() (k8s.Interface, error)
	// WithNameSpace
	GetControllerRevisions(kubeClient k8s.Interface, kind string, objects []string) []string
	GetDaemonSets(kubeClient k8s.Interface, kind string, objects []string) []string
	GetDeployments(kubeClient k8s.Interface, kind string, objects []string) []string
	GetReplicaSets(kubeClient k8s.Interface, kind string, objects []string) []string
	GetStatefulSets(kubeClient k8s.Interface, kind string, objects []string) []string
	GetHorizontalPodAutoscalers(kubeClient k8s.Interface, kind string, objects []string) []string
	GetJobs(kubeClient k8s.Interface, kind string, objects []string) []string
	GetCronJobs(kubeClient k8s.Interface, kind string, objects []string) []string
	GetLeases(kubeClient k8s.Interface, kind string, objects []string) []string
	GetConfigMaps(kubeClient k8s.Interface, kind string, objects []string) []string
	GetEndpoints(kubeClient k8s.Interface, kind string, objects []string) []string
	GetEvents(kubeClient k8s.Interface, kind string, objects []string) []string
	GetLimitRanges(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPersistentVolumeClaims(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPods(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPodTemplates(kubeClient k8s.Interface, kind string, objects []string) []string
	GetReplicationControllers(kubeClient k8s.Interface, kind string, objects []string) []string
	GetResourceQuotas(kubeClient k8s.Interface, kind string, objects []string) []string
	GetSecrets(kubeClient k8s.Interface, kind string, objects []string) []string
	GetServices(kubeClient k8s.Interface, kind string, objects []string) []string
	GetServiceAccounts(kubeClient k8s.Interface, kind string, objects []string) []string
	GetEndpointSlices(kubeClient k8s.Interface, kind string, objects []string) []string
	GetIngresses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetNetworkPolicies(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPodDisruptionBudgets(kubeClient k8s.Interface, kind string, objects []string) []string
	GetRoles(kubeClient k8s.Interface, kind string, objects []string) []string
	GetRoleBindings(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPodPresets(kubeClient k8s.Interface, kind string, objects []string) []string
	// WithoutNameSpace
	GetMutatingWebhookConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string
	GetValidatingWebhookConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string
	GetAuditSinks(kubeClient k8s.Interface, kind string, objects []string) []string
	GetCertificateSigningRequests(kubeClient k8s.Interface, kind string, objects []string) []string
	GetComponentStatuses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetNamespaces(kubeClient k8s.Interface, kind string, objects []string) []string
	GetNodes(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPersistentVolumes(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPodSecurityPolicies(kubeClient k8s.Interface, kind string, objects []string) []string
	GetFlowSchemas(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPriorityLevelConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string
	GetIngressClasses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetRuntimeClasses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetClusterRoles(kubeClient k8s.Interface, kind string, objects []string) []string
	GetClusterRoleBindings(kubeClient k8s.Interface, kind string, objects []string) []string
	GetPriorityClasses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetCSIDrivers(kubeClient k8s.Interface, kind string, objects []string) []string
	GetCSINodes(kubeClient k8s.Interface, kind string, objects []string) []string
	GetStorageClasses(kubeClient k8s.Interface, kind string, objects []string) []string
	GetVolumeAttachments(kubeClient k8s.Interface, kind string, objects []string) []string
}
