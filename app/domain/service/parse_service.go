package service

type ParseService interface {
	IsKubeKind(jsonContent string, kubeKind string) bool

	GetKubeKind(jsonContent string) string
	GetKubeName(jsonContent string) string
	GetKubeNameSpace(jsonContent string) string
}
