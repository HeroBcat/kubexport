package service

import (
	"strings"

	"github.com/tidwall/gjson"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
)

type parseService struct {
}

func NewParseService() serv.ParseService {
	return parseService{}
}

func (s parseService) IsKubeKind(jsonContent string, kubeKind string) bool {

	kind := strings.ToLower(gjson.Get(jsonContent, "kind").String())
	if strings.ToLower(kind) == strings.ToLower(kubeKind) {
		return true
	}
	return false
}

func (s parseService) GetKubeKind(jsonContent string) string {
	return gjson.Get(jsonContent, "kind").String()
}

func (s parseService) GetKubeName(jsonContent string) string {
	return strings.ToLower(gjson.Get(jsonContent, "metadata.name").String())
}

func (s parseService) GetKubeNameSpace(jsonContent string) string {
	return strings.ToLower(gjson.Get(jsonContent, "metadata.namespace").String())
}
