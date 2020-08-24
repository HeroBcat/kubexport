package service

import (
	k8s "k8s.io/client-go/kubernetes"
)

type KubeService interface {
	ReadKubernetesObject(kubeClient k8s.Interface, kind string, name string) string
}
