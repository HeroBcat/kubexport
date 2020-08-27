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
	ExportObjects(objects map[string][]string, targetDir string, isHelmChart bool)
	ExportYaml(path, targetDir string, isHelmChart bool)
}

type kubeUseCase struct {
	kubectl serv.KubectlService
	cleanup serv.CleanUpService
	parse   serv.ParseService
}

func NewKubeUseCase(kubectl serv.KubectlService, cleanup serv.CleanUpService, parse serv.ParseService) KubeUseCase {
	return kubeUseCase{
		kubectl,
		cleanup,
		parse,
	}
}

type ExportsYaml struct {
	Name      string              `yaml:"name"`
	Alias     []string            `yaml:"alias"`
	Resources []map[string]string `yaml:"resources"`
	NameSpace string              `json:"namespace"`
}

func (uc kubeUseCase) ExportYaml(path, targetDir string, isHelmChart bool) {

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
				childDir := ""
				projectName := content.Name
				if strings.HasPrefix(resource, projectName) {
					childDir = projectName
				}

				if childDir == "" {
					for _, alias := range content.Alias {
						if strings.HasPrefix(resource, alias) {
							childDir = projectName
							break
						}
					}
				}

				uc.export(kind, resource, targetDir, childDir, isHelmChart)
			}
		}

		uc.exportKustomization(targetDir, content.Name, content.NameSpace)
	}

}

func (uc kubeUseCase) ExportObjects(objects map[string][]string, targetDir string, isHelmChart bool) {

	for kind, resources := range objects {
		for _, resource := range resources {
			uc.export(kind, resource, targetDir, "", isHelmChart)
		}
	}
}

func (uc kubeUseCase) export(kind, resource, targetDir, childDir string, isHelmChart bool) {
	dict := uc.getContent(resource, kind)
	if dict == nil {
		return
	}

	filename := uc.getFileName(kind, resource, targetDir, childDir)

	if isHelmChart {
		uc.exportToHelmChart(dict, filename)
	} else {
		uc.exportToYaml(dict, filename)
	}
}

func (uc kubeUseCase) createDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return err
}

func (uc kubeUseCase) getFileName(kind, resource, targetDir, subDir string) string {
	fmt.Println("export " + kind + ": " + resource)
	if subDir != "" {
		targetDir = filepath.Join(targetDir, subDir)
	}
	uc.createDirIfNotExist(targetDir)
	filename := resource + "." + strings.ToLower(kind) + ".yaml"
	if subDir != "" {
		filename = strings.ToLower(kind) + ".yaml"
	}
	filename = filepath.Join(targetDir, filename)
	return filename
}

func (uc kubeUseCase) getContent(resource, kind string) map[string]interface{} {

	name := ""
	namespace := ""

	rs := strings.Split(resource, "@")
	switch len(rs) {
	case 0:
		return nil
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

	if uc.parse.IsKubeKind(dict, constant.Deployments) {
		dict = uc.cleanup.CleanUpDeployment(dict)
	}
	return dict
}

func (uc kubeUseCase) exportToYaml(dict map[string]interface{}, filename string) {
	jBytes, err := json.Marshal(dict)
	if err != nil {
		log.Fatal(err)
	}

	yBytes, err := yaml.JSONToYAML(jBytes)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("---\n%s", yBytes)
	file.WriteString(content)
	file.Close()
}

func (uc kubeUseCase) exportKustomization(targetDir, projectName, namespace string) {

	dir := filepath.Join(targetDir, projectName)

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	kinds := make([]string, 0)
	for _, file := range files {
		kinds = append(kinds, file.Name())
	}

	dict := make(map[string]interface{}, 0)
	dict["resources"] = kinds
	if namespace != "" {
		dict["namespace"] = namespace
	}

	filename := uc.getFileName("kustomization", projectName, targetDir, projectName)
	uc.exportToYaml(dict, filename)
}

func (uc kubeUseCase) exportToHelmChart(dict map[string]interface{}, filename string) {

}
