package service

import (
	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/app/infrastructure/service/utils"
)

type parseService struct {
}

func NewParseService() serv.ParseService {
	return parseService{}
}

func (s parseService) IsKubeKind(dict map[string]interface{}, kubeKind string) bool {
	if kind := utils.IsKey(dict, "kind"); kind != "" {
		if kind == kubeKind {
			return true
		}
	}
	return false
}

func (s parseService) GetKubeKind(dict map[string]interface{}) string {
	if kind := utils.IsKey(dict, "kind"); kind != "" {
		return kind
	}
	return ""
}

func (s parseService) GetKubeName(dict map[string]interface{}) string {
	if metadata := utils.ISMapKey(dict, "metadata"); metadata != nil {
		if name := utils.IsKey(metadata, "name"); name != "" {
			return name
		}
	}
	return ""
}

func (s parseService) GetKubeNameSpace(dict map[string]interface{}) string {
	if metadata := utils.ISMapKey(dict, "metadata"); metadata != nil {
		if namespace := utils.IsKey(metadata, "namespace"); namespace != "" {
			return namespace
		}
	}
	return ""
}
