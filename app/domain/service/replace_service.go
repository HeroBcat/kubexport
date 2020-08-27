package service

type ReplaceService interface {
	ReplaceValues(dict map[string]interface{}, kind, project string) (map[string]interface{}, map[string]interface{})
}
