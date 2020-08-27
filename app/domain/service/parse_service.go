package service

type ParseService interface {
	IsKubeKind(dict map[string]interface{}, kubeKind string) bool

	GetKubeKind(dict map[string]interface{}) string
	GetKubeName(dict map[string]interface{}) string
	GetKubeNameSpace(dict map[string]interface{}) string
}
