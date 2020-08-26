package service

type CleanUpService interface {
	IsKubeKind(dict map[string]interface{}, kubeKind string) bool
	CleanUpDeployment(dict map[string]interface{}) map[string]interface{}
	CleanUpStatus(dict map[string]interface{}) map[string]interface{}
	CleanUpMetadata(dict map[string]interface{}) map[string]interface{}
}
