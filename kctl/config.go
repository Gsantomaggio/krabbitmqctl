package kctl

type Config struct {
	KubeConfig  string
	ContextName string
	NameSpace   string
	ServiceName string
	CtlCommand  []string
	PodName string
}
