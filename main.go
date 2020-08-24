package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/HeroBcat/kubexport/app/registry"
	"github.com/HeroBcat/kubexport/config/constant"
)

func main() {

	var (
		// sourceDir      string
		targetDir      string
		isHelmChart    bool = false
		yamlPath       string
		inputWithKinds = make([][]string, len(constant.KubeKinds))
	)

	rootCmd := &cobra.Command{
		Use:   "kubexport",
		Short: "A tool to export yamls from local or k8s environment",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(1)
			}

			name := args[0]
			isKindExist := false
			isYamlExist := false

			kinds := make(map[string][]string, 0)

			for idx, input := range inputWithKinds {
				if len(input) > 0 {
					kinds[constant.KubeKinds[idx]] = input
				}
			}

			_, err := os.Stat(yamlPath)
			if os.IsNotExist(err) && len(kinds) == 0 {
				fmt.Println("No object given.")
				os.Exit(1)
			}

			isKindExist = len(kinds) > 0
			isYamlExist = err == nil

			if targetDir == "" {
				if isYamlExist {
					targetDir = filepath.Join(filepath.Dir(yamlPath), "kubexport", name)
				} else {
					home, err := homedir.Dir()
					if err != nil {
						fmt.Errorf("ERR: %s", err)
						os.Exit(1)
					}
					targetDir = filepath.Join(home, "kubexport", name)
				}
			}

			fmt.Println(targetDir)

			kubeUseCase := registry.BuildKubeUseCase()

			if isYamlExist {
				if isHelmChart {
					kubeUseCase.ExportYamlToHelm(yamlPath, targetDir)
				} else {
					kubeUseCase.ExportYaml(yamlPath, targetDir)
				}
			}

			if isKindExist {
				if isHelmChart {
					kubeUseCase.ExportObjectsToHelm(kinds, targetDir)
				} else {
					kubeUseCase.ExportObjects(kinds, targetDir)
				}

			}
		},
	}

	// rootCmd.Flags().StringVar(&sourceDir, "local", "", "Specify the directory of the local source yaml files")
	rootCmd.Flags().StringVar(&targetDir, "target", "", "Specify the directory to create files")
	// rootCmd.Flags().BoolVar(&isHelmChart, "helm", false, "Specify conversion to helm files")
	rootCmd.Flags().StringVar(&yamlPath, "yaml", "", "Specify the path of yaml file")

	for idx, kind := range constant.KubeKinds {
		rootCmd.Flags().StringSliceVar(&inputWithKinds[idx], strings.ToLower(kind), nil, "Specify the names of "+kind)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
