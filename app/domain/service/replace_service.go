package service

type ReplaceService interface {
	ReplaceValues(jsonContent, valuesJson string, kind, project string, configs map[string]interface{}) (string, string)
}
