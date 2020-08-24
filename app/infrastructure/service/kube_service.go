package service

import (
	"strings"

	k8s "k8s.io/client-go/kubernetes"

	cli "github.com/HeroBcat/kubexport/app/domain/client"
	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/config/constant"
)

type kubeService struct {
	cli.KubeClient
}

func NewKubeService(client cli.KubeClient) serv.KubeService {
	return kubeService{
		client,
	}
}

func (s kubeService) ReadKubernetesObject(kubeClient k8s.Interface, kind string, name string) string {
	var yamlFiles []string

	objects := []string{name}

	switch strings.ToLower(kind) {
	// WithNameSpace
	case strings.ToLower(constant.ControllerRevisions):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetControllerRevisions(kubeClient, kind, objects))
	case strings.ToLower(constant.DaemonSets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetDaemonSets(kubeClient, kind, objects))
	case strings.ToLower(constant.Deployments):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetDeployments(kubeClient, kind, objects))
	case strings.ToLower(constant.ReplicaSets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetReplicaSets(kubeClient, kind, objects))
	case strings.ToLower(constant.StatefulSets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetStatefulSets(kubeClient, kind, objects))
	case strings.ToLower(constant.HorizontalPodAutoscalers):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetHorizontalPodAutoscalers(kubeClient, kind, objects))
	case strings.ToLower(constant.Jobs):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetJobs(kubeClient, kind, objects))
	case strings.ToLower(constant.CronJobs):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetCronJobs(kubeClient, kind, objects))
	case strings.ToLower(constant.Leases):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetLeases(kubeClient, kind, objects))
	case strings.ToLower(constant.ConfigMaps):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetConfigMaps(kubeClient, kind, objects))
	case strings.ToLower(constant.Endpoints):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetEndpoints(kubeClient, kind, objects))
	case strings.ToLower(constant.Events):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetEvents(kubeClient, kind, objects))
	case strings.ToLower(constant.LimitRanges):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetLimitRanges(kubeClient, kind, objects))
	case strings.ToLower(constant.PersistentVolumeClaims):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPersistentVolumeClaims(kubeClient, kind, objects))
	case strings.ToLower(constant.Pods):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPods(kubeClient, kind, objects))
	case strings.ToLower(constant.PodTemplates):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPodTemplates(kubeClient, kind, objects))
	case strings.ToLower(constant.ReplicationControllers):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetReplicationControllers(kubeClient, kind, objects))
	case strings.ToLower(constant.ResourceQuotas):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetResourceQuotas(kubeClient, kind, objects))
	case strings.ToLower(constant.Secrets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetSecrets(kubeClient, kind, objects))
	case strings.ToLower(constant.Services):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetServices(kubeClient, kind, objects))
	case strings.ToLower(constant.ServiceAccounts):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetServiceAccounts(kubeClient, kind, objects))
	case strings.ToLower(constant.EndpointSlices):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetEndpointSlices(kubeClient, kind, objects))
	case strings.ToLower(constant.Ingresses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetIngresses(kubeClient, kind, objects))
	case strings.ToLower(constant.NetworkPolicies):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetNetworkPolicies(kubeClient, kind, objects))
	case strings.ToLower(constant.PodDisruptionBudgets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPodDisruptionBudgets(kubeClient, kind, objects))
	case strings.ToLower(constant.Roles):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetRoles(kubeClient, kind, objects))
	case strings.ToLower(constant.RoleBindings):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetRoleBindings(kubeClient, kind, objects))
	case strings.ToLower(constant.PodPresets):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPodPresets(kubeClient, kind, objects))

	// WithoutNameSpace
	case strings.ToLower(constant.MutatingWebhookConfigurations):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetMutatingWebhookConfigurations(kubeClient, kind, objects))
	case strings.ToLower(constant.ValidatingWebhookConfigurations):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetValidatingWebhookConfigurations(kubeClient, kind, objects))
	case strings.ToLower(constant.AuditSinks):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetAuditSinks(kubeClient, kind, objects))
	case strings.ToLower(constant.CertificateSigningRequests):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetCertificateSigningRequests(kubeClient, kind, objects))
	case strings.ToLower(constant.ComponentStatuses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetComponentStatuses(kubeClient, kind, objects))
	case strings.ToLower(constant.Namespaces):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetNamespaces(kubeClient, kind, objects))
	case strings.ToLower(constant.Nodes):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetNodes(kubeClient, kind, objects))
	case strings.ToLower(constant.PersistentVolumes):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPersistentVolumes(kubeClient, kind, objects))
	case strings.ToLower(constant.PodSecurityPolicies):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPodSecurityPolicies(kubeClient, kind, objects))
	case strings.ToLower(constant.FlowSchemas):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetFlowSchemas(kubeClient, kind, objects))
	case strings.ToLower(constant.PriorityLevelConfigurations):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPriorityLevelConfigurations(kubeClient, kind, objects))
	case strings.ToLower(constant.IngressClasses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetIngressClasses(kubeClient, kind, objects))
	case strings.ToLower(constant.RuntimeClasses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetRuntimeClasses(kubeClient, kind, objects))
	case strings.ToLower(constant.ClusterRoles):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetClusterRoles(kubeClient, kind, objects))
	case strings.ToLower(constant.ClusterRoleBindings):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetClusterRoleBindings(kubeClient, kind, objects))
	case strings.ToLower(constant.PriorityClasses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetPriorityClasses(kubeClient, kind, objects))
	case strings.ToLower(constant.CSIDrivers):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetCSIDrivers(kubeClient, kind, objects))
	case strings.ToLower(constant.CSINodes):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetCSINodes(kubeClient, kind, objects))
	case strings.ToLower(constant.StorageClasses):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetStorageClasses(kubeClient, kind, objects))
	case strings.ToLower(constant.VolumeAttachments):
		yamlFiles = appendSlice(yamlFiles, objects, s.GetVolumeAttachments(kubeClient, kind, objects))
	}

	if len(yamlFiles) > 0 {
		return yamlFiles[0]
	}
	return ""
}

func appendSlice(yamlFiles, objects, subSlice []string) []string {
	if len(objects) > 0 {
		for _, v := range subSlice {
			yamlFiles = append(yamlFiles, v)
		}
	}
	return yamlFiles
}
