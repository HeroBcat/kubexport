package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/config/constant"
)

var defaultFilePerm = os.FileMode(0664)
var helmValuesJson = "{}"

type KubeUseCase interface {
	ExportObjects(objects map[string][]string, targetDir string, helmPath string)
	ExportYaml(path, targetDir string, helmPath string)
}

type kubeUseCase struct {
	kubectl serv.KubectlService
	parse   serv.ParseService
	replace serv.ReplaceService
}

func NewKubeUseCase(kubectl serv.KubectlService, parse serv.ParseService, replace serv.ReplaceService) KubeUseCase {
	return kubeUseCase{
		kubectl,
		parse,
		replace,
	}
}

type ExportsYaml struct {
	Name      string              `yaml:"name"`
	Alias     []string            `yaml:"alias"`
	Resources []map[string]string `yaml:"resources"`
	NameSpace string              `json:"namespace"`
}

func (uc kubeUseCase) ExportYaml(path, targetDir string, helmPath string) {

	fileByte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	contents := make([]ExportsYaml, 0)

	err = yaml.Unmarshal(fileByte, &contents)
	if err != nil {
		log.Fatal(err)
	}

	isHelmChart := uc.isHelmChart(helmPath)

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

				uc.export(kind, resource, targetDir, childDir, helmPath)
			}
		}

		if !isHelmChart {
			uc.exportKustomization(targetDir, content.Name, content.NameSpace)
		}
	}

	if isHelmChart {
		uc.exportHelper(targetDir)
		uc.exportChart(targetDir)
		uc.exportToYaml(helmValuesJson, filepath.Join(targetDir, "values.yaml"), "")
	}

}

func (uc kubeUseCase) ExportObjects(objects map[string][]string, targetDir string, helmPath string) {

	for kind, resources := range objects {
		for _, resource := range resources {
			uc.export(kind, resource, targetDir, "", helmPath)
		}
	}
}

func (uc kubeUseCase) export(kind, resource, targetDir, childDir string, helmPath string) {
	json := uc.getContent(resource, kind)
	if json == nil {
		return
	}

	fmt.Println("export " + kind + ": " + resource)

	if uc.isHelmChart(helmPath) {

		filename := uc.getFileName(kind, resource, filepath.Join(targetDir, "templates"), childDir)

		p := strings.ReplaceAll(childDir, "-", "_")
		p = strings.ReplaceAll(p, ".", "_")
		uc.exportToHelmChart(*json, kind, p, uc.getValuesContent(helmPath), filename)
	} else {
		filename := uc.getFileName(kind, resource, targetDir, childDir)
		uc.exportToYaml(*json, filename, "---\n")
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

	if subDir != "" {
		subDir = strings.ReplaceAll(subDir, "-", "_")
		subDir = strings.ReplaceAll(subDir, ".", "_")
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

func (uc kubeUseCase) getContent(resource, kind string) *string {

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

	jsonContent := uc.kubectl.KubectlGet(kind, name, namespace)

	jsonContent, _ = sjson.Delete(jsonContent, "status")

	jsonContent, _ = sjson.Delete(jsonContent, "metadata.annotations")
	jsonContent, _ = sjson.Delete(jsonContent, "metadata.creationTimestamp")
	jsonContent, _ = sjson.Delete(jsonContent, "metadata.generation")
	jsonContent, _ = sjson.Delete(jsonContent, "metadata.resourceVersion")
	jsonContent, _ = sjson.Delete(jsonContent, "metadata.selfLink")
	jsonContent, _ = sjson.Delete(jsonContent, "metadata.uid")

	if uc.parse.IsKubeKind(jsonContent, constant.Deployments) {
		jsonContent, _ = sjson.Delete(jsonContent, "spec.progressDeadlineSeconds")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.revisionHistoryLimit")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.strategy")

		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.metadata")

		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.dnsPolicy")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.restartPolicy")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.schedulerName")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.terminationGracePeriodSeconds")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.securityContext")
		jsonContent, _ = sjson.Delete(jsonContent, "spec.template.spec.serviceAccount")

		for i := 0; i < len(gjson.Get(jsonContent, "spec.template.spec.containers").Array()); i++ {
			key := fmt.Sprintf("spec.template.spec.containers.%d.", i)
			jsonContent, _ = sjson.Delete(jsonContent, key+"terminationMessagePath")
			jsonContent, _ = sjson.Delete(jsonContent, key+"terminationMessagePolicy")
			jsonContent, _ = sjson.Delete(jsonContent, key+"resources")
		}

	}

	return &jsonContent
}

func (uc kubeUseCase) exportToYaml(jsonContent string, filename string, prefix string) {

	yBytes, err := yaml.JSONToYAML([]byte(jsonContent))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("%s%s", prefix, yBytes)
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
		if strings.HasPrefix(file.Name(), "_values.") {
			continue
		}
		kinds = append(kinds, file.Name())
	}

	jsonContent, _ := sjson.Set("", "resources", kinds)
	if namespace != "" {
		jsonContent, _ = sjson.Set(jsonContent, "namespace", namespace)
	}

	filename := uc.getFileName("kustomization", projectName, targetDir, projectName)
	uc.exportToYaml(jsonContent, filename, "---\n")
}

func (uc kubeUseCase) exportToHelmChart(json string, kind, project string, configs map[string]interface{}, filename string) {
	if project == "" {
		project = "global"
	}

	chart, values := uc.replace.ReplaceValues(json, helmValuesJson, kind, project, configs)
	helmValuesJson = values

	uc.exportToYaml(chart, filename, "")

}

func (uc kubeUseCase) getValuesContent(path string) map[string]interface{} {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	dict := make(map[string]interface{}, 0)
	err = yaml.Unmarshal(content, &dict)
	if err != nil {
		log.Fatal(err)
	}

	return dict
}

func (uc kubeUseCase) getJsonContent(path string) string {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func (uc kubeUseCase) exportHelper(targetDir string) {
	const defaultHelpers = `{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 24 -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 24 -}}
{{- end -}}
`

	filename := filepath.Join(targetDir, "templates", "_helpers.tpl")
	if err := ioutil.WriteFile(filename, []byte(defaultHelpers), 0644); err != nil {
		log.Fatal(err)
	}

}

func (uc kubeUseCase) exportChart(targetDir string) {

	filename := filepath.Join(targetDir, "Chart.yaml")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		jsonContent, _ := sjson.Set("", "apiVersion", "v1")
		jsonContent, _ = sjson.Set(jsonContent, "name", "openbayes")
		jsonContent, _ = sjson.Set(jsonContent, "description", "A Helm chart for Kubernetes")
		jsonContent, _ = sjson.Set(jsonContent, "version", "0.1.0")

		uc.exportToYaml(jsonContent, filename, "")
	}

}

func (uc kubeUseCase) isHelmChart(helmPath string) bool {
	_, err := os.Stat(helmPath)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return true

}
