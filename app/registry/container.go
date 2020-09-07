package registry

import (
	"sync"

	app "github.com/HeroBcat/kubexport/app/application"
	serv "github.com/HeroBcat/kubexport/app/infrastructure/service"
)

var (
	kubeUseCase  *app.KubeUseCase
	localUseCase *app.LocalUseCase

	kubeOnce  sync.Once
	localOnce sync.Once
)

func BuildKubeUseCase() app.KubeUseCase {

	kubeOnce.Do(func() {
		if kubeUseCase == nil {
			kubectlService := serv.NewKubectlService()
			parseService := serv.NewParseService()
			replaceService := serv.NewReplaceService(parseService)
			useCase := app.NewKubeUseCase(kubectlService, parseService, replaceService)
			kubeUseCase = &useCase
		}
	})

	return *kubeUseCase
}

func BuildLocalUseCase() app.LocalUseCase {

	localOnce.Do(func() {
		if localUseCase == nil {
			parseService := serv.NewParseService()
			useCase := app.NewLocalUseCase(parseService)
			localUseCase = &useCase
		}
	})

	return *localUseCase
}
