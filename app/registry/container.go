package registry

import (
	"sync"

	app "github.com/HeroBcat/kubexport/app/application"
	cli "github.com/HeroBcat/kubexport/app/infrastructure/client"
	serv "github.com/HeroBcat/kubexport/app/infrastructure/service"
)

var (
	kubeUseCase *app.KubeUseCase

	kubeOnce sync.Once
)

func BuildKubeUseCase() app.KubeUseCase {

	kubeOnce.Do(func() {
		if kubeUseCase == nil {
			kubeClient := cli.NewKubeClient()
			kubeService := serv.NewKubeService(kubeClient)
			useCase := app.NewKubeUseCase(kubeService)
			kubeUseCase = &useCase
		}
	})

	return *kubeUseCase
}
