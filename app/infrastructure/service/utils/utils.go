package utils

func IsKey(dict map[string]interface{}, key string) string {
	if sub, isOK := dict[key].(string); isOK {
		return sub
	}
	return ""
}

func ISMapKey(dict map[string]interface{}, key string) map[string]interface{} {
	if subDict, isOK := dict[key].(map[string]interface{}); isOK {
		return subDict
	}
	return nil
}

func IsListKey(dict map[string]interface{}, key string) []interface{} {
	if subDict, isOK := dict[key].([]interface{}); isOK {
		return subDict
	}
	return nil
}

func IsMap(obj interface{}) map[string]interface{} {
	if dict, isOK := obj.(map[string]interface{}); isOK {
		return dict
	}
	return nil
}

func IsObject(obj interface{}) interface{} {
	if result := IsMap(obj); result != nil {
		return nil
	}
	if result := IsList(obj); result != nil {
		return nil
	}
	return obj
}

func IsList(obj interface{}) []interface{} {
	if list, isOK := obj.([]interface{}); isOK {
		return list
	}
	return nil
}

func IsListObject(obj interface{}) bool {

	if list, isOK := obj.([]interface{}); isOK {
		for _, obj := range list {
			if object := IsObject(obj); object == nil {
				return false
			}
		}
	}
	return true
}
