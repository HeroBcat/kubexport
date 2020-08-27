package service

import (
	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/app/infrastructure/service/utils"
)

type cleanUpService struct {
}

func NewCleanUpService() serv.CleanUpService {
	return cleanUpService{}
}

func (s cleanUpService) CleanUpDeployment(dict map[string]interface{}) map[string]interface{} {

	// Deployment.spec
	if spec := utils.ISMapKey(dict, "spec"); spec != nil {
		delete(spec, "progressDeadlineSeconds")
		delete(spec, "revisionHistoryLimit")
		delete(spec, "strategy")

		// Deployment.spec.template
		if template := utils.ISMapKey(spec, "template"); template != nil {
			delete(template, "metadata")

			// Deployment.spec.template.spec
			if tSpec := utils.ISMapKey(template, "spec"); tSpec != nil {
				tSpec = s.cleanUpSpecContainers(tSpec)
				delete(tSpec, "dnsPolicy")
				delete(tSpec, "restartPolicy")
				delete(tSpec, "schedulerName")
				delete(tSpec, "terminationGracePeriodSeconds")
				delete(tSpec, "securityContext")
				delete(tSpec, "serviceAccount")
			}
		}
	}

	return dict
}

func (s cleanUpService) CleanUpStatus(dict map[string]interface{}) map[string]interface{} {
	delete(dict, "status")
	return dict
}

func (s cleanUpService) CleanUpMetadata(dict map[string]interface{}) map[string]interface{} {
	if metadata := utils.ISMapKey(dict, "metadata"); metadata != nil {
		delete(metadata, "annotations")
		delete(metadata, "creationTimestamp")
		delete(metadata, "generation")
		delete(metadata, "resourceVersion")
		delete(metadata, "selfLink")
		delete(metadata, "uid")
	}
	return dict
}

func (s cleanUpService) cleanUpSpecContainers(dict map[string]interface{}) map[string]interface{} {

	if containers := utils.IsListKey(dict, "containers"); containers != nil {
		for _, container := range containers {
			if containerDict := utils.IsMap(container); containerDict != nil {
				delete(containerDict, "terminationMessagePath")
				delete(containerDict, "terminationMessagePolicy")
				delete(containerDict, "resources")
			}
		}
	}
	return dict
}
