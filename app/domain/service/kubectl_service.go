package service

type KubectlService interface {
	KubectlGet(kind, name, namespace string) string
}
