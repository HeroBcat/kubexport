package application

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
)

type LocalUseCase interface {
	ExportPath(path, targetDir string, isHelmChart bool)
}

type localUseCase struct {
	cleanup serv.CleanUpService
	parse   serv.ParseService
}

func NewLocalUseCase(cleanup serv.CleanUpService, parse serv.ParseService) LocalUseCase {
	return localUseCase{
		cleanup: cleanup,
		parse:   parse,
	}
}

func (uc localUseCase) ExportPath(localPath, targetDir string, isHelmChart bool) {

	tmp := time.Now().Format("20060102150405")

	uc.CopyFiles(localPath, targetDir, tmp, false, isHelmChart)

	copyFile := filepath.Join(targetDir, tmp)
	uc.SplitFiles(copyFile, targetDir, false)

	os.RemoveAll(copyFile)

	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == targetDir {
				return nil
			}
			uc.exportKustomization(path)
		}
		return nil
	})

}

func (uc localUseCase) CopyFiles(localPath, targetDir string, tempName string, isSubDir, isHelmChart bool) {

	stats, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, stat := range stats {
		if stat.IsDir() {
			dir := filepath.Join(localPath, stat.Name())
			uc.CopyFiles(dir, targetDir, tempName, true, isHelmChart)
			continue
		}

		if strings.HasSuffix(stat.Name(), ".yaml") {

			parent := ""
			if list := strings.Split(localPath, "/"); len(list) > 0 && isSubDir {
				parent = list[len(list)-1]
			}

			dir := filepath.Join(targetDir, tempName, parent)
			os.MkdirAll(dir, os.ModePerm)

			if stat.Name() == "kustomization.yaml" {
				source := filepath.Join(localPath, stat.Name())
				target := filepath.Join(dir, stat.Name())
				uc.copy(source, target)
				continue
			}

			source := filepath.Join(localPath, stat.Name())
			target := filepath.Join(dir, "all.yaml")
			uc.copy(source, target)

		}

	}
}

func (uc localUseCase) SplitFiles(copyDir, targetDir string, isSubDir bool) {

	stats, err := ioutil.ReadDir(copyDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, stat := range stats {
		if stat.IsDir() {
			dir := filepath.Join(copyDir, stat.Name())
			uc.SplitFiles(dir, targetDir, true)
			continue
		}

		if stat.Name() == "all.yaml" {
			parent := ""
			if list := strings.Split(copyDir, "/"); len(list) > 0 && isSubDir {
				parent = list[len(list)-1]
			}

			dir := filepath.Join(targetDir, parent)
			os.MkdirAll(dir, os.ModePerm)

			source := filepath.Join(copyDir, stat.Name())
			yBytes, err := ioutil.ReadFile(source)
			if err != nil {
				continue
			}

			contents := strings.Split(string(yBytes), "\n---")
			for _, content := range contents {

				if content == "" {
					continue
				}
				if content == "\n" {
					continue
				}

				dict := make(map[string]interface{}, 0)
				yaml.Unmarshal([]byte(content), &dict)

				dict = uc.cleanup.CleanUpStatus(dict)
				dict = uc.cleanup.CleanUpMetadata(dict)
				dict = uc.cleanup.CleanUpDeployment(dict)

				kind := uc.parse.GetKubeKind(dict)
				if kind == "" {
					continue
				}

				filename := filepath.Join(dir, kind+".yaml")
				if !isSubDir {
					filename = uc.parse.GetKubeName(dict) + "." + strings.ToLower(kind) + ".yaml"
					filename = filepath.Join(dir, filename)
				}

				uc.exportToYaml(dict, filename)
			}

		}

	}

}

func (uc localUseCase) copy(src, dst string) (int64, error) {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (uc localUseCase) exportKustomization(targetDir string) {

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	kinds := make([]string, 0)
	namespace := ""
	isSameNameSpace := true

	for _, stat := range files {

		if !stat.IsDir() {
			file := filepath.Join(targetDir, stat.Name())
			yBytes, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}

			dict := make(map[string]interface{}, 0)
			yaml.Unmarshal(yBytes, &dict)

			kind := uc.parse.GetKubeKind(dict)
			if kind == "" {
				continue
			}
			kinds = append(kinds, strings.ToLower(kind)+".yaml")

			ns := uc.parse.GetKubeNameSpace(dict)
			if namespace == "" {
				namespace = ns
			}
			if namespace != ns {
				isSameNameSpace = false
			}

		}

	}

	dict := make(map[string]interface{}, 0)
	dict["resources"] = kinds
	if isSameNameSpace {
		dict["namespace"] = namespace
	}

	filename := filepath.Join(targetDir, "kustomization.yaml")
	uc.exportToYaml(dict, filename)
}

func (uc localUseCase) exportToYaml(dict map[string]interface{}, filename string) {
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
