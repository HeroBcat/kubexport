package client

import (
	"log"
	"strings"

	core "k8s.io/api/core/v1"
)

func splitNamespace(s string) (string, string) {
	str := strings.Split(s, "@")
	if len(str) == 2 {
		return str[0], str[1]
	} else if len(str) == 1 {
		return str[0], core.NamespaceDefault
	}
	log.Fatal("ERROR : Can not detect Namespace")
	return "", ""
}

func getAPIVersion(selfLink string) string {
	str := strings.Split(selfLink, "/")
	if len(str) > 2 {
		if strings.HasPrefix(str[3], "v") {
			return str[2] + "/" + str[3]
		}
		return str[2]
	}
	log.Fatal("api version not found")
	return ""
}
