package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/config/constant"
)

var defaultFilePerm = os.FileMode(0664)

type KubeUseCase interface {
	ExportObjects(objects map[string][]string, targetDir string)
	ExportYaml(path, targetDir string)
}

type kubeUseCase struct {
	kubectl serv.KubectlService
	cleanup serv.CleanUpService
}

func NewKubeUseCase(kubectl serv.KubectlService, cleanup serv.CleanUpService) KubeUseCase {
	return kubeUseCase{
		kubectl,
		cleanup,
	}
}

type ExportsYaml struct {
	Name      string              `yaml:"name"`
	Resources []map[string]string `yaml:"resources"`
}

func (uc kubeUseCase) ExportYaml(path, targetDir string) {

	fileByte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	contents := make([]ExportsYaml, 0)

	err = yaml.Unmarshal(fileByte, &contents)
	if err != nil {
		log.Fatal(err)
	}

	for _, content := range contents {
		for _, resources := range content.Resources {
			for kind, resource := range resources {
				uc.export(resource, kind, targetDir)
			}
		}
	}
}

func (uc kubeUseCase) ExportObjects(objects map[string][]string, targetDir string) {

	for kind, resources := range objects {
		for _, resource := range resources {
			uc.export(resource, kind, targetDir)
		}
	}
}

func (uc kubeUseCase) export(resource, kind, targetDir string) {

	name := ""
	namespace := ""

	rs := strings.Split(resource, "@")
	switch len(rs) {
	case 0:
		return
	case 1:
		name = resource
	case 2:
		name = rs[0]
		namespace = rs[1]
	}

	jsonFile := uc.kubectl.KubectlGet(kind, name, namespace)

	dict := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonFile), &dict)
	if err != nil {
		log.Fatal(err)
	}

	dict = uc.cleanup.CleanUpStatus(dict)
	dict = uc.cleanup.CleanUpMetadata(dict)

	if uc.cleanup.IsKubeKind(dict, constant.Deployments) {
		dict = uc.cleanup.CleanUpDeployment(dict)
	}

	jBytes, err := json.Marshal(dict)
	if err != nil {
		log.Fatal(err)
	}

	yBytes, err := yaml.JSONToYAML(jBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("export " + kind + ": " + resource)
	createDirIfNotExist(targetDir)
	filename := resource + "." + strings.ToLower(kind) + ".yaml"
	filename = filepath.Join(targetDir, filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(yBytes)
	file.Close()

}

func createDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return err
}
