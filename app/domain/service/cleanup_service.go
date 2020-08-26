package service

type CleanUpService interface {
	IsKubeKind(dict map[string]interface{}, kubeKind string) bool

	GetKubeKind(dict map[string]interface{}) string
	GetKubeName(dict map[string]interface{}) string
	GetKubeNameSpace(dict map[string]interface{}) string

	CleanUpDeployment(dict map[string]interface{}) map[string]interface{}
	CleanUpStatus(dict map[string]interface{}) map[string]interface{}
	CleanUpMetadata(dict map[string]interface{}) map[string]interface{}
}
