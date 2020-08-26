package service

import (
	"bytes"
	"fmt"
	"os/exec"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
)

type kubectlService struct {
}

func NewKubectlService() serv.KubectlService {
	return kubectlService{}
}

func (s kubectlService) KubectlGet(kind, name, namespace string) string {

	if namespace == "" {
		namespace = "default"
	}

	cmd := exec.Command("kubectl", "get", kind, name, "-n", namespace, "-o", "json")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return ""
	}
	return out.String()

}
