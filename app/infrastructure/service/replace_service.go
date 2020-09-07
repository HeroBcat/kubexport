package service

import (
	"fmt"
	"strings"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
	"github.com/HeroBcat/kubexport/app/infrastructure/service/utils"
)

type replaceService struct {
	parse serv.ParseService
}

func NewReplaceService(parse serv.ParseService) serv.ReplaceService {
	return replaceService{
		parse: parse,
	}
}

func (s replaceService) ReplaceValues(dict map[string]interface{}, kind, project string) (map[string]interface{}, map[string]interface{}) {

	chart := make(map[string]interface{}, 0)
	values := make(map[string]interface{}, 0)

	kind = strings.ToLower(kind)

	chartKey := fmt.Sprintf(".Values.%s.%s", project, kind)

	for key, value := range dict {
		if key == "apiVersion" {
			chart[key] = value
		} else if key == "kind" {
			chart[key] = value
		} else if obj := utils.IsObject(value); obj != nil {
			chart[key] = s.getChartKey(s.addChartKey(chartKey, key))
			values[key] = value
		} else if subDict := utils.IsMap(value); subDict != nil {
			newChart, newContent := s.replaceDictValues(subDict, s.addChartKey(chartKey, key), project, kind)
			chart[key] = newChart
			values[key] = newContent
		} else if utils.IsListObject(value) {
			chart[key] = s.getListContent(s.addChartKey(chartKey, key))
			values[key] = value
		} else if list := utils.IsList(value); list != nil {
			newChart, newContent := s.replaceListValues(list, s.addChartKey(chartKey, key), project, kind)
			chart[key] = newChart
			values[key] = newContent
		}

	}

	fValues := make(map[string]interface{}, 0)
	fValues[kind] = values

	return chart, fValues
}

func (s replaceService) replaceDictValues(dict map[string]interface{}, xKey, project, kind string) (map[string]interface{}, map[string]interface{}) {

	chart := make(map[string]interface{}, 0)
	values := make(map[string]interface{}, 0)

	chartKey := xKey

	for key, value := range dict {

		if key == "labels" && strings.HasSuffix(strings.ToLower(chartKey), ".metadata") {
			chart["chart"] = "{{.Chart.Name}}-{{.Chart.Version}}"
			chart["heritage"] = "{{.Release.Service}}"
			chart["release"] = "{{.Release.Name}}"
		} else if key == "name" && strings.HasSuffix(strings.ToLower(chartKey), ".metadata") {
			chart[key] = fmt.Sprintf("{{ template \"fullname\" . }}-%s", chart[key])
		} else if key == "data" && strings.HasSuffix(strings.ToLower(chartKey), "configmap") {
			newKey := strings.ReplaceAll(key, ".", "_")
			chart[newKey] = s.getListContent(s.addChartKey(chartKey, newKey))
			values[newKey] = value
			continue
		}

		if obj := utils.IsObject(value); obj != nil {
			chart[key] = s.getChartKey(s.addChartKey(chartKey, key))
			values[key] = value

			if key == "app" && strings.HasSuffix(strings.ToLower(chartKey), "labels") {
				chart[key] = fmt.Sprintf("{{.Release.Name}}-%s", chart[key])
			}

		} else if subDict := utils.IsMap(value); subDict != nil {
			newChart, newContent := s.replaceDictValues(subDict, s.addChartKey(chartKey, key), project, kind)
			chart[key] = newChart
			values[key] = newContent
		} else if utils.IsListObject(value) {
			chart[key] = s.getListContent(s.addChartKey(chartKey, key))
			values[key] = value
		} else if list := utils.IsList(value); list != nil {
			newChart, newContent := s.replaceListValues(list, s.addChartKey(chartKey, key), project, kind)
			chart[key] = newChart
			values[key] = newContent
		}

	}

	return chart, values
}

func (s replaceService) replaceListValues(list []interface{}, xKey, project, kind string) ([]interface{}, []interface{}) {

	chart := make([]interface{}, 0)
	values := make([]interface{}, 0)

	chartKey := xKey

	for idx, value := range list {
		chartKey = s.addChartIdx(xKey, idx)
		if obj := utils.IsObject(value); obj != nil {
			chart = append(chart, chartKey)
			values = append(values, value)
		} else if subDict := utils.IsMap(value); subDict != nil {
			newChart, newContent := s.replaceDictValues(subDict, chartKey, project, kind)
			chart = append(chart, newChart)
			values = append(values, newContent)
		} else if subList := utils.IsList(value); subList != nil {
			newChart, newContent := s.replaceListValues(subList, chartKey, project, kind)
			chart = append(chart, newChart)
			values = append(values, newContent)
		}
	}

	return chart, values
}

func (s replaceService) addChartKey(chartKey string, keys ...string) string {

	for _, key := range keys {
		chartKey = fmt.Sprintf("%s.%s", chartKey, key)
	}
	return chartKey
}

func (s replaceService) addChartIdx(chartKey string, idx int) string {
	return fmt.Sprintf("%s%d", chartKey, idx)
}

func (s replaceService) getChartKey(chartKey string) string {
	return fmt.Sprintf("{{%s}}", chartKey)
}

func (s replaceService) getListContent(chartKey string) string {
	return fmt.Sprintf(`{{- range %s}}
{{ . | quote }}
{{- end}}`, chartKey)
}
