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
	replace serv.ReplaceService
}

func NewKubeUseCase(kubectl serv.KubectlService, cleanup serv.CleanUpService, parse serv.ParseService, replace serv.ReplaceService) KubeUseCase {
	return kubeUseCase{
		kubectl,
		cleanup,
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

		if !isHelmChart {
			uc.exportKustomization(targetDir, content.Name, content.NameSpace)
		}
	}

	if isHelmChart {
		uc.appendValues(targetDir)
		uc.exportHelper(targetDir)
		uc.exportChart(targetDir)
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

	fmt.Println("export " + kind + ": " + resource)

	filename := uc.getFileName(kind, resource, targetDir, childDir)

	if isHelmChart {
		valuesName := uc.getValuesName(kind, resource, targetDir, childDir)
		filename = uc.getFileName(kind, resource, filepath.Join(targetDir, "templates"), childDir)

		p := strings.ReplaceAll(childDir, "-", "_")
		p = strings.ReplaceAll(p, ".", "_")
		uc.exportToHelmChart(dict, kind, p, filename, valuesName)
	} else {
		uc.exportToYaml(dict, filename, "---\n")
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

func (uc kubeUseCase) getValuesName(kind, resource, targetDir, subDir string) string {

	if subDir != "" {
		targetDir = filepath.Join(targetDir, subDir)
	}
	uc.createDirIfNotExist(targetDir)
	filename := "_values." + resource + "." + strings.ToLower(kind) + ".yaml"
	if subDir != "" {
		filename = "_values." + strings.ToLower(kind) + ".yaml"
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

func (uc kubeUseCase) exportToYaml(dict interface{}, filename string, prefix string) {
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

	dict := make(map[string]interface{}, 0)
	dict["resources"] = kinds
	if namespace != "" {
		dict["namespace"] = namespace
	}

	filename := uc.getFileName("kustomization", projectName, targetDir, projectName)
	uc.exportToYaml(dict, filename, "---\n")
}

func (uc kubeUseCase) exportToHelmChart(dict map[string]interface{}, kind, project, filename, valuesName string) {
	if project == "" {
		project = "global"
	}
	chart, values := uc.replace.ReplaceValues(dict, kind, project)

	uc.exportToYaml(chart, filename, "")
	uc.exportToYaml(values, valuesName, "")

}

func (uc kubeUseCase) appendValues(targetDir string) {

	values := make(map[string]map[string]interface{}, 0)

	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasPrefix(info.Name(), "_values.") && strings.HasSuffix(info.Name(), ".yaml") {

			ps := strings.Split(path, "/")
			if len(ps) > 1 {
				project := ps[len(ps)-2]
				project = strings.ReplaceAll(project, "-", "_")
				project = strings.ReplaceAll(project, ".", "_")
				name := ps[len(ps)-1]
				if strings.HasSuffix(targetDir, project) {
					project = "global"
				}

				fs := strings.Split(name, ".")
				if len(fs) > 1 {
					name = fs[len(fs)-2]
				}

				dict := uc.getValuesContent(path)
				if projDict, ok := values[project]; ok {
					projDict[name] = dict[name]
				} else {
					values[project] = make(map[string]interface{}, 0)
					values[project][name] = dict[name]
				}

			}

		}

		return nil
	})

	uc.exportToYaml(values, filepath.Join(targetDir, "values.yaml"), "")

	states, err := ioutil.ReadDir(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, state := range states {
		if state.Name() == "templates" {
			continue
		} else if state.Name() == "values.yaml" {
			continue
		}

		os.RemoveAll(filepath.Join(targetDir, state.Name()))
	}

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
		dict := make(map[string]string, 0)
		dict["apiVersion"] = "v1"
		dict["name"] = "openbayes"
		dict["description"] = "A Helm chart for Kubernetes"
		dict["version"] = "0.1.0"

		uc.exportToYaml(dict, filename, "")
	}

}
