package service

type CleanUpService interface {
	CleanUpDeployment(dict map[string]interface{}) map[string]interface{}
	CleanUpStatus(dict map[string]interface{}) map[string]interface{}
	CleanUpMetadata(dict map[string]interface{}) map[string]interface{}
}
