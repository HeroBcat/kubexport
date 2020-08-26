package service

import (
	serv "github.com/HeroBcat/kubexport/app/domain/service"
)

type cleanUpService struct {
}

func NewCleanUpService() serv.CleanUpService {
	return cleanUpService{}
}

func (s cleanUpService) IsKubeKind(dict map[string]interface{}, kubeKind string) bool {
	if kind := s.checkKey(dict, "kind"); kind != "" {
		if kind == kubeKind {
			return true
		}
	}
	return false
}

func (s cleanUpService) GetKubeKind(dict map[string]interface{}) string {
	if kind := s.checkKey(dict, "kind"); kind != "" {
		return kind
	}
	return ""
}

func (s cleanUpService) GetKubeName(dict map[string]interface{}) string {
	if metadata := s.checkMapKey(dict, "metadata"); metadata != nil {
		if name := s.checkKey(metadata, "name"); name != "" {
			return name
		}
	}
	return ""
}

func (s cleanUpService) GetKubeNameSpace(dict map[string]interface{}) string {
	if metadata := s.checkMapKey(dict, "metadata"); metadata != nil {
		if namespace := s.checkKey(metadata, "namespace"); namespace != "" {
			return namespace
		}
	}
	return ""
}

func (s cleanUpService) CleanUpDeployment(dict map[string]interface{}) map[string]interface{} {

	// Deployment.spec
	if spec := s.checkMapKey(dict, "spec"); spec != nil {
		delete(spec, "progressDeadlineSeconds")
		delete(spec, "revisionHistoryLimit")
		delete(spec, "strategy")

		// Deployment.spec.template
		if template := s.checkMapKey(spec, "template"); template != nil {
			delete(template, "metadata")

			// Deployment.spec.template.spec
			if tSpec := s.checkMapKey(template, "spec"); tSpec != nil {
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
	if metadata := s.checkMapKey(dict, "metadata"); metadata != nil {
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

	if containers := s.checkListKey(dict, "containers"); containers != nil {
		for _, container := range containers {
			if containerDict := s.isMap(container); containerDict != nil {
				delete(containerDict, "terminationMessagePath")
				delete(containerDict, "terminationMessagePolicy")
				delete(containerDict, "resources")
			}
		}
	}
	return dict
}

func (s cleanUpService) checkKey(dict map[string]interface{}, key string) string {
	if sub, isOK := dict[key].(string); isOK {
		return sub
	}
	return ""
}

func (s cleanUpService) checkMapKey(dict map[string]interface{}, key string) map[string]interface{} {
	if subDict, isOK := dict[key].(map[string]interface{}); isOK {
		return subDict
	}
	return nil
}

func (s cleanUpService) isMap(dict interface{}) map[string]interface{} {
	if subDict, isOK := dict.(map[string]interface{}); isOK {
		return subDict
	}
	return nil
}

func (s cleanUpService) checkListKey(dict map[string]interface{}, key string) []interface{} {
	if subDict, isOK := dict[key].([]interface{}); isOK {
		return subDict
	}
	return nil
}
