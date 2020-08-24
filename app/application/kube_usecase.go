package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	serv "github.com/HeroBcat/kubexport/app/domain/service"
)

var defaultFilePerm = os.FileMode(0664)

type KubeUseCase interface {
	ExportYaml(path, targetDir string)
	ExportObjects(objects map[string][]string, targetDir string)

	ExportYamlToHelm(path, targetDir string)
	ExportObjectsToHelm(objects map[string][]string, targetDir string)
}

type kubeUseCase struct {
	service serv.KubeService
}

func NewKubeUseCase(service serv.KubeService) KubeUseCase {
	return kubeUseCase{
		service,
	}
}

func (uc kubeUseCase) ExportYaml(path, targetDir string) {
	kubeClient, err := DefaultKubeClient()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connect k8s success")
	}

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
				yaml := uc.service.ReadKubernetesObject(kubeClient, kind, resource)
				fmt.Println("export " + kind + ": " + resource)
				dir := filepath.Join(targetDir, content.Name)
				createDirIfNotExist(dir)
				filename := resource + "." + strings.ToLower(kind) + ".yaml"
				filename = filepath.Join(dir, filename)
				file, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				file.WriteString(yaml)
				file.Close()
			}

		}

	}
}

func (uc kubeUseCase) ExportObjects(objects map[string][]string, targetDir string) {

	kubeClient, err := DefaultKubeClient()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connect k8s success")
	}

	for kind, resources := range objects {
		for _, resource := range resources {
			yaml := uc.service.ReadKubernetesObject(kubeClient, kind, resource)
			fmt.Println("export " + kind + ": " + resource)
			createDirIfNotExist(targetDir)
			filename := resource + "." + strings.ToLower(kind) + ".yaml"
			filename = filepath.Join(targetDir, filename)
			file, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString(yaml)
			file.Close()
		}

	}

}

func (uc kubeUseCase) ExportYamlToHelm(path, targetDir string) {

}

func (uc kubeUseCase) ExportObjectsToHelm(objects map[string][]string, targetDir string) {

}

type ExportsYaml struct {
	Name      string              `yaml:"name"`
	Resources []map[string]string `yaml:"resources"`
}

func DefaultKubeClient() (k8s.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get kubernetes config: %s", err)
	}
	return k8s.NewForConfig(config)
}

func createDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return err
}
