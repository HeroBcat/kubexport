package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

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

func (s replaceService) ReplaceValues(jsonContent, valuesJson string, kind, project string, configs map[string]interface{}) (string, string) {

	var err error
	chartJson := jsonContent
	kind = strings.ToLower(kind)

	if objMap, ok := configs[kind]; ok {

		if config := utils.IsMap(objMap); config != nil {

			for key, objKey := range config {
				valueKey := ""
				if valueKey, ok = objKey.(string); !ok {
					continue
				}

				result := gjson.Get(jsonContent, key)

				if result.IsObject() {

					valuesJson, err = sjson.SetRaw(valuesJson, s.getKey(project, kind, valueKey), result.Raw)
					if err != nil {
						log.Fatal(err)
					}

					chartJson, err = sjson.Set(chartJson, valueKey, s.getChartKey(s.getKey(project, kind, valueKey)))
					if err != nil {
						log.Fatal(err)
					}

				} else if result.IsArray() {

					count := len(result.Array())
					for i := 0; i < count; i++ {

						newKey := valueKey
						if count > 1 {
							newKey = fmt.Sprintf("%s%d", valueKey, i)
						}

						valuesJson, err = sjson.SetRaw(valuesJson, s.getKey(project, kind, newKey), result.Array()[i].Raw)
						if err != nil {
							log.Fatal(err)
						}

						chartJson, err = sjson.Set(chartJson, strings.Replace(key, "#", strconv.Itoa(i), 1), s.getChartKey(s.getKey(project, kind, newKey)))
						if err != nil {
							log.Fatal(err)
						}

					}

				} else {

					valuesJson, err = sjson.Set(valuesJson, s.getKey(project, kind, valueKey), result.Value())
					if err != nil {
						log.Fatal(err)
					}

					chartJson, err = sjson.Set(chartJson, s.getKey(project, kind, valueKey), s.getChartKey(s.getKey(project, kind, valueKey)))
					if err != nil {
						log.Fatal(err)
					}

				}

			}

		}

	}

	return chartJson, valuesJson
}

func (s replaceService) getKey(project, kind string, key string) string {

	if strings.HasPrefix(key, "__") {
		return fmt.Sprintf("global.%s", key[2:])
	} else if strings.HasPrefix(key, "_") {
		return fmt.Sprintf("%s.%s", project, key[1:])
	}
	return fmt.Sprintf("%s.%s.%s", project, kind, key)
}

func (s replaceService) getChartKey(chartKey string) string {
	return fmt.Sprintf("{{.Values.%s}}", chartKey)
}
