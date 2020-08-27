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

func IsMap(dict interface{}) map[string]interface{} {
	if subDict, isOK := dict.(map[string]interface{}); isOK {
		return subDict
	}
	return nil
}
