package client

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch "k8s.io/api/batch/v1"
	batchv1b1 "k8s.io/api/batch/v1beta1"
	certificates "k8s.io/api/certificates/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	flowcontrol "k8s.io/api/flowcontrol/v1alpha1"
	policy "k8s.io/api/policy/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	cli "github.com/HeroBcat/kubexport/app/domain/client"
)

type kubeClient struct {
}

func NewKubeClient() cli.KubeClient {
	return kubeClient{}
}

func (c kubeClient) DefaultKubeClient() (k8s.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get kubernetes config: %s", err)
	}
	return k8s.NewForConfig(config)
}

// WithNameSpace
func (c kubeClient) GetControllerRevisions(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AppsV1().ControllerRevisions(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetDaemonSets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AppsV1().DaemonSets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}

		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = apps.DaemonSetStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetDeployments(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AppsV1().Deployments(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}

		if result.Kind == "" {
			result.Kind = "Deployment"
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = apps.DeploymentStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetReplicaSets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AppsV1().ReplicaSets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = apps.ReplicaSetStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetStatefulSets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AppsV1().StatefulSets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = apps.StatefulSetStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetHorizontalPodAutoscalers(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = autoscaling.HorizontalPodAutoscalerStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetJobs(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.BatchV1().Jobs(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = batch.JobStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetCronJobs(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.BatchV1beta1().CronJobs(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = batchv1b1.CronJobStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetLeases(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoordinationV1().Leases(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetConfigMaps(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetEndpoints(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().Endpoints(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetEvents(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().Events(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetLimitRanges(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().LimitRanges(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPersistentVolumeClaims(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPods(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().Pods(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.PodStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPodTemplates(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().PodTemplates(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetReplicationControllers(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().ReplicationControllers(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.ReplicationControllerStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetResourceQuotas(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().ResourceQuotas(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.ResourceQuotaStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetSecrets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().Secrets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetServices(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.ServiceStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetServiceAccounts(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.CoreV1().ServiceAccounts(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetEndpointSlices(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.DiscoveryV1beta1().EndpointSlices(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetIngresses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.ExtensionsV1beta1().Ingresses(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = extensions.IngressStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetNetworkPolicies(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.ExtensionsV1beta1().NetworkPolicies(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPodDisruptionBudgets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.PolicyV1beta1().PodDisruptionBudgets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = policy.PodDisruptionBudgetStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetRoles(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.RbacV1().Roles(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetRoleBindings(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.RbacV1().RoleBindings(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPodPresets(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		objectName, namespace := splitNamespace(v)
		result, err := kubeClient.SettingsV1alpha1().PodPresets(namespace).Get(context.Background(), objectName, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

// WithoutNameSpace
func (c kubeClient) GetMutatingWebhookConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetValidatingWebhookConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetAuditSinks(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.AuditregistrationV1alpha1().AuditSinks().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetCertificateSigningRequests(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.CertificatesV1beta1().CertificateSigningRequests().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		result.Status = certificates.CertificateSigningRequestStatus{}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetComponentStatuses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.CoreV1().ComponentStatuses().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetNamespaces(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.NamespaceStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetNodes(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.CoreV1().Nodes().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.NodeStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPersistentVolumes(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.CoreV1().PersistentVolumes().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = core.PersistentVolumeStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPodSecurityPolicies(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.PolicyV1beta1().PodSecurityPolicies().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetFlowSchemas(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.FlowcontrolV1alpha1().FlowSchemas().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = flowcontrol.FlowSchemaStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPriorityLevelConfigurations(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.FlowcontrolV1alpha1().PriorityLevelConfigurations().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		result.Status = flowcontrol.PriorityLevelConfigurationStatus{}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetIngressClasses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.NetworkingV1beta1().IngressClasses().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetRuntimeClasses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.NodeV1beta1().RuntimeClasses().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetClusterRoles(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.RbacV1().ClusterRoles().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetClusterRoleBindings(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.RbacV1().ClusterRoleBindings().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetPriorityClasses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.SchedulingV1().PriorityClasses().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetCSIDrivers(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.StorageV1().CSIDrivers().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetCSINodes(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.StorageV1().CSINodes().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetStorageClasses(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.StorageV1().StorageClasses().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (c kubeClient) GetVolumeAttachments(kubeClient k8s.Interface, kind string, objects []string) []string {
	var yamlFiles []string
	for _, v := range objects {
		result, err := kubeClient.StorageV1().VolumeAttachments().Get(context.Background(), v, meta.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if result.Kind == "" {
			result.Kind = kind
		}
		if result.APIVersion == "" {
			result.APIVersion = getAPIVersion(result.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}
