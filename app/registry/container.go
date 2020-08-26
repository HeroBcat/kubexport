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
			cleanUpService := serv.NewCleanUpService()
			useCase := app.NewKubeUseCase(kubectlService, cleanUpService)
			kubeUseCase = &useCase
		}
	})

	return *kubeUseCase
}

func BuildLocalUseCase() app.LocalUseCase {

	localOnce.Do(func() {
		if localUseCase == nil {
			cleanUpService := serv.NewCleanUpService()
			useCase := app.NewLocalUseCase(cleanUpService)
			localUseCase = &useCase
		}
	})

	return *localUseCase
}
